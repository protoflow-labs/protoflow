package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/code"
	pgrpc "github.com/protoflow-labs/protoflow/gen/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"net/url"
)

func GetGRPCTypeInfo(host string) ([]*gen.GRPCService, error) {
	m := manager.NewReflectionManager(host, manager.WithProtocol(bufcurl.ReflectProtocolGRPCV1Alpha))
	cleanup, err := m.Init()
	if err != nil {
		return nil, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	services, err := m.ResolveServices()
	if err != nil {
		return nil, errors.Wrapf(err, "error resolving services")
	}

	var blocks []*gen.GRPCService
	seen := map[protoreflect.FullName]struct{}{}
	for _, sd := range services {
		if _, ok := seen[sd.FullName()]; ok {
			continue
		}
		seen[sd.FullName()] = struct{}{}

		var methods []*gen.GRPCMethod
		for i := 0; i < sd.Methods().Len(); i++ {
			m := sd.Methods().Get(i)
			md, err := NewMethodDescriptor(m)
			if err != nil {
				return nil, errors.Wrapf(err, "error creating method descriptor")
			}
			ti, err := md.Proto()
			if err != nil {
				return nil, errors.Wrapf(err, "error getting proto")
			}
			methods = append(methods, &gen.GRPCMethod{
				Name:     string(m.Name()),
				TypeInfo: ti,
			})
		}
		blocks = append(blocks, &gen.GRPCService{
			Name:    string(sd.Name()),
			Package: string(sd.ParentFile().Package()),
			Methods: methods,
			// TODO breadchris what are the cases where passing the file path here is not correct?
			// for example, the path is not absolute, it is relative to the proto directory.
			File: sd.ParentFile().Path(),
		})
	}
	return blocks, nil
}

func EnumerateResourceBlocks(server *pgrpc.Server, isLangService bool) ([]*gen.Node, error) {
	if server.Host == "" {
		return nil, errors.New("host is required")
	}

	u, err := url.Parse(server.Host)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse host %s", server.Host)
	}

	conn, err := grpc.Dial(u.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to python server at %s", server.Host)
	}

	// TODO breadchris there is some repeat code, the grpc package has some code from Buf that does reflection already
	svcDesc, err := allMethodsViaReflection(context.Background(), conn)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get all methods via reflection")
	}

	methodDesc := AllMethodsForServices(svcDesc)

	log.Debug().Str("server", server.Host).Msgf("found %d methods", len(methodDesc))

	var blocks []*gen.Node
	for _, m := range methodDesc {
		serviceName := m.GetService().GetName()
		methodName := m.GetName()

		grpcInfo := &pgrpc.Method{
			Package: m.GetFile().GetPackage(),
			Service: serviceName,
			Method:  methodName,
		}

		block := &gen.Node{
			Id:   uuid.New().String(),
			Name: methodName,
			Type: &gen.Node_Grpc{
				Grpc: &pgrpc.GRPC{
					Type: &pgrpc.GRPC_Method{
						Method: grpcInfo,
					},
				},
			},
		}
		if isLangService {
			block.Type = &gen.Node_Code{
				Code: &code.Code{
					Type: &code.Code_Function{
						Function: &code.Function{},
					},
				},
			}
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

type MethodDescriptor struct {
	MethodDesc protoreflect.MethodDescriptor
	DescLookup map[string]protoreflect.MessageDescriptor
	EnumLookup map[string]protoreflect.EnumDescriptor
	//FileDesc    protoreflect.FileDescriptor
	FileBuilder *builder.FileBuilder
}

type MethodDescriptorProto struct {
	DescLookup map[string]*descriptorpb.DescriptorProto
	EnumLookup map[string]*descriptorpb.EnumDescriptorProto
}

// TODO breadchris make this more generic to allow different type of descriptors such as MessageDescriptor
func NewMethodDescriptor(md protoreflect.MethodDescriptor) (*MethodDescriptor, error) {
	m := &MethodDescriptor{
		MethodDesc:  md,
		DescLookup:  map[string]protoreflect.MessageDescriptor{},
		EnumLookup:  map[string]protoreflect.EnumDescriptor{},
		FileBuilder: builder.NewFile(string(md.Name()) + "File"),
	}
	m.buildTypeLookup(md.Input())
	m.buildTypeLookup(md.Output())
	return m, nil
}

// TODO breadchris placeholder until the above is implemented
func ResolveTypeLookup(
	msgDesc protoreflect.MessageDescriptor,
	descLookup map[string]protoreflect.MessageDescriptor,
	enumLookup map[string]protoreflect.EnumDescriptor,
) (map[string]protoreflect.MessageDescriptor, map[string]protoreflect.EnumDescriptor) {
	msgs := []protoreflect.MessageDescriptor{msgDesc}
	fileBuilder := builder.NewFile("File")
	for len(msgs) > 0 {
		msg := msgs[0]
		msgs = msgs[1:]
		descLookup[string(msg.FullName())] = msg

		wmsg, err := desc.WrapMessage(msg)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to wrap message %s", msg.FullName())
			continue
		}
		mb, err := builder.FromMessage(wmsg)
		if e := fileBuilder.GetMessage(wmsg.GetName()); e == nil {
			fileBuilder = fileBuilder.AddMessage(mb)
		}

		fields := msg.Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.Get(i)
			lookupName := string(f.FullName())

			oneOf := f.ContainingOneof()
			if oneOf != nil {
				oneOfFields := oneOf.Fields()
				for j := 0; j < oneOfFields.Len(); j++ {
					c := oneOfFields.Get(j)
					// TODO breadchris replace with m.FileBuilder.GetMessage
					msgName := string(c.Message().FullName())
					if _, ok := descLookup[msgName]; ok {
						continue
					}
					msgs = append(msgs, c.Message())
				}
			} else {
				switch f.Kind() {
				case protoreflect.MessageKind:
					// TODO breadchris replace with m.FileBuilder.GetMessage
					msgName := string(f.Message().FullName())
					if _, ok := descLookup[msgName]; ok {
						continue
					}
					msgs = append(msgs, f.Message())
				case protoreflect.EnumKind:
					enumLookup[lookupName] = f.Enum()
					wenum, err := desc.WrapEnum(f.Enum())
					if err != nil {
						log.Warn().Err(err).Msgf("unable to wrap message %s", f.Enum())
						continue
					}
					eb, err := builder.FromEnum(wenum)
					if e := fileBuilder.GetEnum(eb.GetName()); e == nil {
						fileBuilder = fileBuilder.AddEnum(eb)
					}
				}
			}
		}
	}
	return descLookup, enumLookup
}

func (m *MethodDescriptor) buildTypeLookup(msgDesc protoreflect.MessageDescriptor) {
	msgs := []protoreflect.MessageDescriptor{msgDesc}
	for len(msgs) > 0 {
		msg := msgs[0]
		msgs = msgs[1:]
		m.DescLookup[string(msg.FullName())] = msg

		wmsg, err := desc.WrapMessage(msg)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to wrap message %s", msg.FullName())
			continue
		}
		mb, err := builder.FromMessage(wmsg)
		if e := m.FileBuilder.GetMessage(wmsg.GetName()); e == nil {
			m.FileBuilder = m.FileBuilder.AddMessage(mb)
		}

		fields := msg.Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.Get(i)
			lookupName := string(f.FullName())

			oneOf := f.ContainingOneof()
			if oneOf != nil {
				oneOfFields := oneOf.Fields()
				for j := 0; j < oneOfFields.Len(); j++ {
					c := oneOfFields.Get(j)
					// TODO breadchris replace with m.FileBuilder.GetMessage
					msgName := string(c.Message().FullName())
					if _, ok := m.DescLookup[msgName]; ok {
						continue
					}
					msgs = append(msgs, c.Message())
				}
			} else {
				switch f.Kind() {
				case protoreflect.MessageKind:
					// TODO breadchris replace with m.FileBuilder.GetMessage
					msgName := string(f.Message().FullName())
					if _, ok := m.DescLookup[msgName]; ok {
						continue
					}
					msgs = append(msgs, f.Message())
				case protoreflect.EnumKind:
					m.EnumLookup[lookupName] = f.Enum()
					wenum, err := desc.WrapEnum(f.Enum())
					if err != nil {
						log.Warn().Err(err).Msgf("unable to wrap message %s", f.Enum())
						continue
					}
					eb, err := builder.FromEnum(wenum)
					if e := m.FileBuilder.GetEnum(eb.GetName()); e == nil {
						m.FileBuilder = m.FileBuilder.AddEnum(eb)
					}
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
		Input:       descMethod.GetInputType().AsDescriptorProto(),
		Output:      descMethod.GetOutputType().AsDescriptorProto(),
		DescLookup:  d.DescLookup,
		EnumLookup:  d.EnumLookup,
		MethodDesc:  descMethod.AsMethodDescriptorProto(),
		PackageName: string(m.MethodDesc.ParentFile().Package()),
	}, nil
}

func (m *MethodDescriptor) Print() (string, error) {
	// TODO breadchris implement
	return "", nil
}
