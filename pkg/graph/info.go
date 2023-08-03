package graph

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/samber/lo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// TODO breadchris if there is only one field, set name of message to just the name of the one field.
// this is canonical in grpc
func MessageFromTypes(name string, types []protoreflect.MessageDescriptor) (*desc.MessageDescriptor, error) {
	if len(types) == 1 {
		return desc.WrapMessage(types[0])
	}
	mb := builder.NewMessage(name)
	if len(types) == 0 {
		return mb.Build()
	}
	var addedFields []string
	for _, t := range types {
		wt, err := desc.WrapMessage(t)
		if err != nil {
			return nil, errors.Wrapf(err, "error wrapping message %s", name)
		}
		msgBuilder, err := builder.FromMessage(wt)
		if err != nil {
			return nil, errors.Wrapf(err, "error building message %s", name)
		}
		fm := builder.FieldTypeMessage(msgBuilder)

		fieldName := string(t.Name())
		if lo.Contains(addedFields, fieldName) {
			return nil, errors.Errorf("duplicate field %s", name)
		}

		mb = mb.AddField(builder.NewField(fieldName, fm))
	}
	return mb.Build()
}

func NewInfoFromType(name string, msg protoreflect.ProtoMessage) (*Info, error) {
	d, err := desc.WrapMessage(msg.ProtoReflect().Descriptor())
	if err != nil {
		return nil, err
	}
	req := builder.RpcTypeImportedMessage(d, false)
	res := builder.RpcTypeImportedMessage(d, false)

	s := builder.NewService("Service")
	b := builder.NewMethod(name, req, res)
	s.AddMethod(b)

	m, err := b.Build()
	if err != nil {
		return nil, err
	}

	mthd, err := grpc.NewMethodDescriptor(m.UnwrapMethod())
	if err != nil {
		return nil, err
	}
	return &Info{
		Method: mthd,
	}, nil
}
