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

func EnumerateResourceBlocks(resource *gen.Resource) ([]*gen.Node, error) {
	var (
		g             *gen.GRPCService
		nodes         []*gen.Node
		err           error
		isLangService bool
	)

	switch resource.Type.(type) {
	case *gen.Resource_LanguageService:
		l := resource.GetLanguageService()
		g = l.Grpc
		isLangService = true
	case *gen.Resource_GrpcService:
		g = resource.GetGrpcService()
	default:
		log.Debug().Interface("type", resource.Type).Msg("resource cannot be enumerated")
	}

	if g != nil {
		nodes, err = nodesFromGRPC(resource.Id, g, isLangService)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to enumerate grpc service %s", g.Host)
			nodes = []*gen.Node{}
		}
	}
	return nodes, nil
}

type MethodDescriptor struct {
	DescLookup map[string]*descriptorpb.DescriptorProto
	EnumLookup map[string]*descriptorpb.EnumDescriptorProto
}

func NewMethodDescriptor(msg *desc.MessageDescriptor) *MethodDescriptor {
	m := &MethodDescriptor{
		DescLookup: make(map[string]*descriptorpb.DescriptorProto),
		EnumLookup: make(map[string]*descriptorpb.EnumDescriptorProto),
	}
	m.buildTypeLookup(msg)
	return m
}

func (m *MethodDescriptor) buildTypeLookup(msgDesc *desc.MessageDescriptor) {
	msgs := []*desc.MessageDescriptor{msgDesc}
	for len(msgs) > 0 {
		msg := msgs[0]
		msgs = msgs[1:]
		m.DescLookup[msg.GetFullyQualifiedName()] = msg.AsDescriptorProto()
		for _, f := range msg.GetFields() {
			lookupName := f.GetFullyQualifiedName()

			oneOf := f.GetOneOf()
			if oneOf != nil {
				choices := oneOf.GetChoices()
				for _, c := range choices {
					if _, ok := m.DescLookup[lookupName]; ok {
						continue
					}
					msgs = append(msgs, c.GetMessageType())
				}
			} else {
				switch f.GetType() {
				case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
					if _, ok := m.DescLookup[lookupName]; ok {
						continue
					}
					m.DescLookup[lookupName] = f.GetMessageType().AsDescriptorProto()
					msgs = append(msgs, f.GetMessageType())
				case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
					m.EnumLookup[lookupName] = f.GetEnumType().AsEnumDescriptorProto()
				}
			}
		}
	}
}

func nodesFromGRPC(resourceID string, service *gen.GRPCService, isLangService bool) ([]*gen.Node, error) {
	if service.Host == "" {
		return nil, errors.New("host is required")
	}

	conn, err := grpc.Dial(service.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to python server at %s", service.Host)
	}

	// TODO breadchris there is some repeat code, the grpc package has some code from Buf that does reflection already
	methodDesc, err := allMethodsViaReflection(context.Background(), conn)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get all methods via reflection")
	}

	log.Debug().Str("service", service.Host).Msgf("found %d methods", len(methodDesc))

	var blocks []*gen.Node
	for _, m := range methodDesc {
		serviceName := m.GetService().GetName()
		methodName := m.GetName()

		grpcInfo := &gen.GRPC{
			Package: m.GetFile().GetPackage(),
			Service: serviceName,
			Method:  methodName,
		}

		block := &gen.Node{
			Id:         uuid.New().String(),
			Name:       methodName,
			ResourceId: resourceID,
			Config: &gen.Node_Grpc{
				Grpc: grpcInfo,
			},
		}
		if isLangService {
			block.Config = &gen.Node_Function{
				Function: &gen.Function{},
			}
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}
