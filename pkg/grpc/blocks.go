package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/jhump/protoreflect/desc"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/descriptorpb"
)

func EnumerateResourceBlocks(resource *gen.Resource) ([]*gen.Block, error) {
	var (
		blocks []*gen.Block
		err    error
	)
	switch resource.Type.(type) {
	case *gen.Resource_GrpcService:
		g := resource.GetGrpcService()
		blocks, err = blocksFromGRPC(g)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to enumerate grpc service %s", g.Host)
			blocks = []*gen.Block{}
		}
	}

	resource.Blocks = blocks
	return blocks, nil
}

type MethodDescriptor struct {
	descLookup map[string]*descriptorpb.DescriptorProto
	enumLookup map[string]*descriptorpb.EnumDescriptorProto
}

func NewMethodDescriptor(msg *desc.MessageDescriptor) *MethodDescriptor {
	m := &MethodDescriptor{
		descLookup: make(map[string]*descriptorpb.DescriptorProto),
		enumLookup: make(map[string]*descriptorpb.EnumDescriptorProto),
	}
	m.buildTypeLookup(msg)
	return m
}

func (m *MethodDescriptor) buildTypeLookup(msgDesc *desc.MessageDescriptor) {
	msgs := []*desc.MessageDescriptor{msgDesc}
	for len(msgs) > 0 {
		msg := msgs[0]
		msgs = msgs[1:]
		m.descLookup[msg.GetFullyQualifiedName()] = msg.AsDescriptorProto()
		for _, f := range msg.GetFields() {
			lookupName := f.GetFullyQualifiedName()

			oneOf := f.GetOneOf()
			if oneOf != nil {
				choices := oneOf.GetChoices()
				for _, c := range choices {
					if _, ok := m.descLookup[lookupName]; ok {
						continue
					}
					msgs = append(msgs, c.GetMessageType())
				}
			} else {
				switch f.GetType() {
				case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
					if _, ok := m.descLookup[lookupName]; ok {
						continue
					}
					m.descLookup[lookupName] = f.GetMessageType().AsDescriptorProto()
					msgs = append(msgs, f.GetMessageType())
				case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
					m.enumLookup[lookupName] = f.GetEnumType().AsEnumDescriptorProto()
				}
			}
		}
	}
}

func blocksFromGRPC(service *gen.GRPCService) ([]*gen.Block, error) {
	if service.Host == "" {
		return nil, errors.New("host is required")
	}

	conn, err := grpc.Dial(service.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to python server at %s", service.Host)
	}

	methodDesc, err := AllMethodsViaReflection(context.Background(), conn)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get all methods via reflection")
	}

	log.Debug().Str("service", service.Host).Msgf("found %d methods", len(methodDesc))

	var blocks []*gen.Block
	for _, m := range methodDesc {
		serviceName := m.GetService().GetName()
		methodName := m.GetName()

		md := NewMethodDescriptor(m.GetInputType())

		blocks = append(blocks, &gen.Block{
			Id:   uuid.New().String(),
			Name: serviceName + "." + methodName,
			Type: &gen.Block_Grpc{
				Grpc: &gen.GRPC{
					Package:    m.GetFile().GetPackage(),
					Service:    serviceName,
					Method:     methodName,
					Input:      m.GetInputType().AsDescriptorProto(),
					Output:     m.GetOutputType().AsDescriptorProto(),
					DescLookup: md.descLookup,
					EnumLookup: md.enumLookup,
					MethodDesc: m.AsMethodDescriptorProto(),
				},
			},
		})
	}
	return blocks, nil
}
