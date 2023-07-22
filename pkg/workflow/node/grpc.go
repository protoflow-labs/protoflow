package node

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type GRPCNode struct {
	BaseNode
	*gen.GRPC
}

var _ graph.Node = &GRPCNode{}

func NewGRPCNode(node *gen.Node) *GRPCNode {
	return &GRPCNode{
		BaseNode: NewBaseNode(node),
		GRPC:     node.GetGrpc(),
	}
}

func (n *GRPCNode) getMethodFromServer(r *resource.GRPCResource, protocol bufcurl.ReflectProtocol) (protoreflect.MethodDescriptor, error) {
	// TODO breadchris I think a grpc resource should have a host that has a protocol
	m := manager.NewReflectionManager("http://"+r.Host, manager.WithProtocol(protocol))
	cleanup, err := m.Init()
	if err != nil {
		return nil, errors.Wrapf(err, "error initializing reflection manager")
	}
	defer cleanup()

	serviceName := n.Package + "." + n.Service
	method, err := m.ResolveMethod(serviceName, n.Method)
	if err != nil {
		return nil, errors.Wrapf(err, "error resolving method")
	}
	return method, nil
}

func (n *GRPCNode) Info(r graph.Resource) (*graph.Info, error) {
	// TODO breadchris what if we want to get the proto from a proto file?

	var (
		method protoreflect.MethodDescriptor
		err    error
	)
	gr, ok := r.(*resource.GRPCResource)
	if !ok {
		return nil, errors.New("grpc resource is not supported")
	}
	method, err = n.getMethodFromServer(gr, bufcurl.ReflectProtocolGRPCV1Alpha)
	if err != nil {
		// TODO breadchris is there a cleaner way to determine if the server supports v1?
		method, err = n.getMethodFromServer(gr, bufcurl.ReflectProtocolGRPCV1)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting method from server")
		}
	}

	md, err := grpc.NewMethodDescriptor(method)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating method descriptor")
	}
	return &graph.Info{
		Method: md,
	}, nil
}
