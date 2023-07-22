package generate

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/protoflow-labs/protoflow/templates"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"
)

type LanguageManager interface {
	GenerateGRPCService(r *resource.LanguageServiceResource) error
}

type NodeJSManager struct {
	codeRoot bucket.Bucket
}

var _ LanguageManager = &NodeJSManager{}

func NewNodeJSManager(c bucket.Bucket) (*NodeJSManager, error) {
	codeRoot, err := c.WithDir("nodejs")
	if err != nil {
		return nil, errors.Wrapf(err, "error creating bucket for %s", "nodejs")
	}
	return &NodeJSManager{
		codeRoot: codeRoot,
	}, nil
}

type Method struct {
	Name string
}

type FunctionTemplate struct {
	Node *node.FunctionNode
}

type ServiceTemplate struct {
	Runtime string
	Methods []Method
}

func (s *NodeJSManager) GenerateGRPCService(r *resource.LanguageServiceResource) error {
	var err error

	tmpl, err := s.generateServiceTemplate(r)
	if err != nil {
		return errors.Wrapf(err, "error scaffolding functions")
	}

	// TODO breadchris this needs to run once for the project
	//err = s.generateServiceProto(tmpl)
	//if err != nil {
	//	return errors.Wrapf(err, "error generating service protos")
	//}

	err = s.generateServices(tmpl)
	if err != nil {
		return errors.Wrapf(err, "error generating services")
	}
	return nil
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

func (s *NodeJSManager) UpdateNodeType(n graph.Node, nodeInfo *graph.Info) error {
	funcDirPath, err := s.codeRoot.GetFolder("protos")
	if err != nil {
		return errors.Wrapf(err, "error getting protos dir")
	}
	fileDescs, err := grpc.ParseProtoDir(funcDirPath)
	if err != nil {
		return errors.Wrapf(err, "error parsing protos")
	}
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
			return errors.Wrapf(err, "error building method")
		}

		f, err := builder.FromFile(fd)
		if err != nil {
			return errors.Wrapf(err, "error building proto file")
		}

		// TODO breadchris hack to remove references to the parent
		// does not work for types that have enums and oneofs?
		it, err := desc.WrapMessage(nodeInfo.Method.MethodDesc.Input())
		if err != nil {
			return errors.Wrapf(err, "error wrapping message")
		}
		a, err := builder.FromMessage(it)
		if err != nil {
			return errors.Wrapf(err, "error building message")
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
			return errors.Wrapf(err, "error wrapping message")
		}
		b, err := builder.FromMessage(ot)
		if err != nil {
			return errors.Wrapf(err, "error building message")
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
				return errors.Wrapf(err, "error building method")
			}
			mb = mb.SetRequestType(inputType).SetResponseType(outputType)
			sb = sb.RemoveMethod(m.GetName())
		}

		// replace the existing method and service with the new one
		sb = sb.AddMethod(mb)
		f = f.RemoveService(svc.GetName()).AddService(sb)

		fd, err := f.Build()
		if err != nil {
			return errors.Wrapf(err, "error building proto file")
		}
		err = writeProtoFile(path.Join(funcDirPath, filename), fd)
		if err != nil {
			return errors.Wrapf(err, "error writing proto file")
		}
	}
	return nil
}

func (s *NodeJSManager) GenerateFunctionImpl(r *resource.LanguageServiceResource, n graph.Node) error {
	switch n.(type) {
	case *node.FunctionNode:
		// create function directory
		funcDir := path.Join("functions", n.NormalizedName())
		funcDirPath, err := s.codeRoot.GetFolder(funcDir)
		if err != nil {
			return errors.Wrapf(err, "error creating function directory %s", funcDir)
		}

		method := Method{
			Name: n.NormalizedName(),
		}

		err = templates.TemplateFile("node/function.index.tmpl.js", path.Join(funcDirPath, "index.js"), method)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (s *NodeJSManager) generateServiceTemplate(r *resource.LanguageServiceResource) (*ServiceTemplate, error) {
	tmpl := &ServiceTemplate{
		Runtime: strings.ToLower(r.Runtime.String()),
		Methods: []Method{},
	}
	for _, resNode := range r.Nodes() {
		switch n := resNode.(type) {
		case *node.FunctionNode:

			method := Method{
				Name: n.NormalizedName(),
			}
			tmpl.Methods = append(tmpl.Methods, method)
		}
	}
	return tmpl, nil
}

// TODO breadchris should generating service protos be done for all services at once? or in each service directory?
// TODO breadchris this should only be run once when a project is initialized
func (s *NodeJSManager) generateServiceProto(tmpl *ServiceTemplate) error {
	protosPath, err := s.codeRoot.GetFolder("protos")
	if err != nil {
		return errors.Wrapf(err, "error getting protos folder %s", path.Join(protosPath, "protos"))
	}

	protoPath := path.Join(protosPath, fmt.Sprintf("%s.proto", tmpl.Runtime))
	return templates.TemplateFile("service.tmpl.proto", protoPath, tmpl)
}

func (s *NodeJSManager) generateServices(tmpl *ServiceTemplate) error {
	projectPath, err := s.codeRoot.GetFolder("/")
	if err != nil {
		return errors.Wrapf(err, "error getting project folder")
	}

	return templates.TemplateDir("node/project", projectPath, tmpl)
}
