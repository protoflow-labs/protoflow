package generate

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/rs/zerolog/log"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
	"os"
	"strings"
)

func unlinkFieldBuilder(f *builder.FileBuilder, t *builder.FieldBuilder, nodeInfo *graph.Info) (*builder.FileBuilder, *builder.FieldBuilder) {
	switch t.GetType().GetType() {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		// TODO breadchris a little hacky, but we need to move the message from the file
		parts := strings.Split(t.GetType().GetTypeName(), ".")
		eb := t.GetFile().GetMessage(strings.Join(parts[1:], "."))
		if eb == nil {
			eb = nodeInfo.Method.FileBuilder.GetMessage(strings.Join(parts[1:], "."))
			if eb == nil {
				log.Warn().Msgf("message %s not found", strings.Join(parts[1:], "."))
				return f, nil
			}
		}
		if ex := f.GetMessage(eb.GetName()); ex != nil {
			f = f.RemoveMessage(ex.GetName())
		}
		f, eb = recursivelyUnlinkBuilder(f, eb, nodeInfo)
		f = f.AddMessage(eb)
		t = t.SetType(builder.FieldTypeMessage(eb))
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		parts := strings.Split(t.GetType().GetTypeName(), ".")
		eb := t.GetFile().GetEnum(strings.Join(parts[1:], "."))
		if eb == nil {
			eb = nodeInfo.Method.FileBuilder.GetEnum(strings.Join(parts[1:], "."))
			if eb == nil {
				log.Warn().Msgf("enum %s not found", strings.Join(parts[1:], "."))
				return f, nil
			}
		}
		if ex := f.GetEnum(eb.GetName()); ex != nil {
			f = f.RemoveEnum(ex.GetName())
		}
		builder.Unlink(eb)
		f = f.AddEnum(eb)
		t = t.SetType(builder.FieldTypeEnum(eb))
	}
	builder.Unlink(t)
	return f, t
}

// TODO breadchris this is a hack to get the protos to generate correctly. removing all references prevents additional imports in the protofile.
func recursivelyUnlinkBuilder(f *builder.FileBuilder, b *builder.MessageBuilder, nodeInfo *graph.Info) (*builder.FileBuilder, *builder.MessageBuilder) {
	// recursively unlink the message
	newB := builder.NewMessage(b.GetName())
	builder.Unlink(newB)
	for _, fb := range b.GetChildren() {
		switch t := fb.(type) {
		case *builder.FieldBuilder:
			f, t = unlinkFieldBuilder(f, t, nodeInfo)
			if t != nil {
				newB.AddField(t)
			}
		case *builder.OneOfBuilder:
			newOneOf := builder.NewOneOf(t.GetName())
			for _, v := range t.GetChildren() {
				switch vt := v.(type) {
				case *builder.FieldBuilder:
					f, vt = unlinkFieldBuilder(f, vt, nodeInfo)
					if vt != nil {
						newOneOf.AddChoice(vt)
					}
				default:
					log.Warn().Msgf("unknown type %T", vt)
				}
			}
			newB.AddOneOf(newOneOf)
		case *builder.EnumBuilder:
			builder.Unlink(t)
			newB.AddNestedEnum(t)
		}
	}
	return f, newB
}

func writeProtoFile(filepath string, fd *desc.FileDescriptor) error {
	file, err := manager.PrintFile(fd)
	if err != nil {
		return errors.Wrapf(err, "error printing file")
	}

	b, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrapf(err, "error opening file")
	}

	_, err = b.WriteString(file)
	b.Close()
	if err != nil {
		return errors.Wrapf(err, "error writing file")
	}
	return nil
}

type protoUpdate struct {
	FilePath string
	FileDesc *desc.FileDescriptor
}

func generateProtoType(fileDescs []*desc.FileDescriptor, n graph.Node, nodeInfo *graph.Info) ([]protoUpdate, error) {
	var updates []protoUpdate
	for _, fd := range fileDescs {
		filename := fd.GetName()
		// TODO breadchris this name should be derived from the language server
		svc := fd.FindService("protoflow.nodejsService")
		if svc == nil {
			log.Warn().Msgf("service not found in proto %s", fd.GetName())
			continue
		}
		sb, err := builder.FromService(svc)
		if err != nil {
			return nil, errors.Wrapf(err, "error building method")
		}

		f, err := builder.FromFile(fd)
		if err != nil {
			return nil, errors.Wrapf(err, "error building proto file")
		}

		// TODO breadchris hack to remove references to the parent
		// does not work for types that have enums and oneofs?
		it, err := desc.WrapMessage(nodeInfo.Method.MethodDesc.Input())
		if err != nil {
			return nil, errors.Wrapf(err, "error wrapping message")
		}
		a, err := builder.FromMessage(it)
		if err != nil {
			return nil, errors.Wrapf(err, "error building message")
		}
		f, newA := recursivelyUnlinkBuilder(f, a, nodeInfo)
		existingMsg := f.GetMessage(newA.GetName())
		if existingMsg != nil {
			builder.Unlink(existingMsg)
		}
		f = f.AddMessage(newA)
		inputType := builder.RpcTypeMessage(newA, false)

		ot, err := desc.WrapMessage(nodeInfo.Method.MethodDesc.Output())
		if err != nil {
			return nil, errors.Wrapf(err, "error wrapping message")
		}
		b, err := builder.FromMessage(ot)
		if err != nil {
			return nil, errors.Wrapf(err, "error building message")
		}
		f, newB := recursivelyUnlinkBuilder(f, b, nodeInfo)
		existingMsg = f.GetMessage(newB.GetName())
		if existingMsg != nil {
			builder.Unlink(existingMsg)
		}
		f = f.AddMessage(newB)
		outputType := builder.RpcTypeMessage(newB, false)

		var mb *builder.MethodBuilder

		// see if this method already exists
		m := svc.FindMethodByName(n.NormalizedName())
		if m == nil {
			// create a new method
			mb = builder.NewMethod(n.NormalizedName(), inputType, outputType)
		} else {
			// refactor the discovered method to match the new node type
			mb, err = builder.FromMethod(m)
			if err != nil {
				return nil, errors.Wrapf(err, "error building method")
			}
			mb = mb.SetRequestType(inputType).SetResponseType(outputType)
			sb = sb.RemoveMethod(m.GetName())
		}

		// replace the existing method and service with the new one
		sb = sb.AddMethod(mb)
		f = f.RemoveService(svc.GetName()).AddService(sb)

		fd, err := f.Build()
		if err != nil {
			return nil, errors.Wrapf(err, "error building proto file")
		}
		updates = append(updates, protoUpdate{
			FilePath: filename,
			FileDesc: fd,
		})
	}
	return updates, nil
}
