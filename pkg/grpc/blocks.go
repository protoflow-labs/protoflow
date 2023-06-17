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
	"google.golang.org/protobuf/reflect/protoreflect"
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

type MethodDescriptor struct {
	MethodDesc protoreflect.MethodDescriptor
	Input      protoreflect.MessageDescriptor
	Output     protoreflect.MessageDescriptor
	DescLookup map[string]protoreflect.MessageDescriptor
	EnumLookup map[string]protoreflect.EnumDescriptor
}

type MethodDescriptorProto struct {
	DescLookup map[string]*descriptorpb.DescriptorProto
	EnumLookup map[string]*descriptorpb.EnumDescriptorProto
}

func NewMethodDescriptor(md protoreflect.MethodDescriptor) *MethodDescriptor {
	m := &MethodDescriptor{
		MethodDesc: md,
		DescLookup: map[string]protoreflect.MessageDescriptor{},
		EnumLookup: map[string]protoreflect.EnumDescriptor{},
		Input:      md.Input(),
		Output:     md.Output(),
	}
	m.buildTypeLookup(md.Input())
	m.buildTypeLookup(md.Output())
	return m
}

func (m *MethodDescriptor) buildTypeLookup(msgDesc protoreflect.MessageDescriptor) {
	msgs := []protoreflect.MessageDescriptor{msgDesc}
	for len(msgs) > 0 {
		msg := msgs[0]
		msgs = msgs[1:]
		m.DescLookup[string(msg.FullName())] = msg
		fields := msg.Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.Get(i)
			lookupName := string(f.FullName())

			oneOf := f.ContainingOneof()
			if oneOf != nil {
				oneOfFields := oneOf.Fields()
				for j := 0; j < oneOfFields.Len(); j++ {
					c := oneOfFields.Get(j)
					if _, ok := m.DescLookup[lookupName]; ok {
						continue
					}
					msgs = append(msgs, c.Message())
				}
			} else {
				switch f.Kind() {
				case protoreflect.MessageKind:
					if _, ok := m.DescLookup[lookupName]; ok {
						continue
					}
					m.DescLookup[lookupName] = f.Message()
					msgs = append(msgs, f.Message())
				case protoreflect.EnumKind:
					m.EnumLookup[lookupName] = f.Enum()
				}
			}
		}
	}
}

// Proto returns a proto representation of the MethodDescriptor
func (m *MethodDescriptor) Proto() (*gen.GRPCTypeInfo, error) {
	d := &MethodDescriptorProto{
		DescLookup: map[string]*descriptorpb.DescriptorProto{},
		EnumLookup: map[string]*descriptorpb.EnumDescriptorProto{},
	}
	for k, v := range m.DescLookup {
		m, err := desc.WrapMessage(v)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to wrap message %s", k)
			continue
		}
		d.DescLookup[k] = m.AsDescriptorProto()
	}
	for k, v := range m.EnumLookup {
		e, err := desc.WrapEnum(v)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to wrap enum %s", k)
			continue
		}
		d.EnumLookup[k] = e.AsEnumDescriptorProto()
	}

	// TODO breadchris does protoreflect have a way to get the proto?
	descMethod, err := desc.WrapMethod(m.MethodDesc)
	if err != nil {
		return nil, errors.Wrapf(err, "error wrapping method")
	}

	return &gen.GRPCTypeInfo{
		Input:      descMethod.GetInputType().AsDescriptorProto(),
		Output:     descMethod.GetOutputType().AsDescriptorProto(),
		DescLookup: d.DescLookup,
		EnumLookup: d.EnumLookup,
		MethodDesc: descMethod.AsMethodDescriptorProto(),
	}, nil
}

func (m *MethodDescriptor) Print() (string, error) {
	// TODO breadchris implement
	return "", nil
}
