package node

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type GRPCNode struct {
	BaseNode
	*gen.GRPC
}

var _ Node = &GRPCNode{}

func NewGRPCNode(node *gen.Node) *GRPCNode {
	return &GRPCNode{
		BaseNode: NewBaseNode(node),
		GRPC:     node.GetGrpc(),
	}
}

func (s *GRPCNode) Execute(executor execute.Executor, input execute.Input) (*execute.Result, error) {
	return executor.Execute(activity.ExecuteGRPCNode, s, input)
}

func (s *GRPCNode) getMethodFromServer(r *resource.GRPCResource, protocol bufcurl.ReflectProtocol) (protoreflect.MethodDescriptor, error) {
	// TODO breadchris I think a grpc resource should have a host that has a protocol
	m := manager.NewReflectionManager("http://"+r.Host, manager.WithProtocol(protocol))
	cleanup, err := m.Init()
	if err != nil {
		return nil, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	serviceName := s.Package + "." + s.Service
	method, err := m.ResolveMethod(serviceName, s.Method)
	if err != nil {
		return nil, errors.Wrapf(err, "error resolving method")
	}
	return method, nil
}

func (s *GRPCNode) Info(res resource.Resource) (*Info, error) {
	r, ok := res.(*resource.GRPCResource)
	if !ok {
		return nil, errors.Errorf("resource is not a grpc resource")
	}

	var (
		method protoreflect.MethodDescriptor
		err    error
	)
	method, err = s.getMethodFromServer(r, bufcurl.ReflectProtocolGRPCV1Alpha)
	if err != nil {
		// TODO breadchris is there a cleaner way to determine if the server supports v1?
		method, err = s.getMethodFromServer(r, bufcurl.ReflectProtocolGRPCV1)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting method from server")
		}
	}

	methodProto, err := manager.GetProtoForMethod(s.Package, s.Service, method)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting proto for method")
	}

	descMethod, err := desc.WrapMethod(method)
	md := grpc.NewMethodDescriptor(descMethod.GetInputType())
	typeInfo := &gen.GRPCTypeInfo{
		Input:      descMethod.GetInputType().AsDescriptorProto(),
		Output:     descMethod.GetOutputType().AsDescriptorProto(),
		DescLookup: md.DescLookup,
		EnumLookup: md.EnumLookup,
		MethodDesc: descMethod.AsMethodDescriptorProto(),
	}

	return &Info{
		MethodProto: methodProto,
		TypeInfo:    typeInfo,
	}, nil
}
