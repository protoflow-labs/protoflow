package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/bufcurl"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

type GRPCResource struct {
	*BaseResource
	*gen.GRPCService
}

func (r *GRPCResource) Init() (func(), error) {
	// TODO breadchris this is a hack to get the grpc server running, this is not ideal
	if !strings.HasPrefix(r.Host, "http://") {
		r.Host = "http://" + r.Host
	}
	//if err := ensureRunning(r.Host); err != nil {
	//	// TODO breadchris ignore errors for now
	//	// return nil, errors.Wrapf(err, "unable to get the %s grpc server running", r.Name())
	//	return nil, nil
	//}
	return nil, nil
}

func (r *GRPCResource) getMethodFromServer(n *node.GRPCNode, protocol bufcurl.ReflectProtocol) (protoreflect.MethodDescriptor, error) {
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

func (r *GRPCResource) Info(n node.Node) (*node.Info, error) {
	// TODO breadchris what if we want to get the proto from a proto file?

	var (
		method protoreflect.MethodDescriptor
		err    error
	)
	gn, ok := n.(*node.GRPCNode)
	if !ok {
		return nil, errors.Errorf("node is not a GRPC node")
	}

	method, err = r.getMethodFromServer(gn, bufcurl.ReflectProtocolGRPCV1Alpha)
	if err != nil {
		// TODO breadchris is there a cleaner way to determine if the server supports v1?
		method, err = r.getMethodFromServer(gn, bufcurl.ReflectProtocolGRPCV1)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting method from server")
		}
	}

	md, err := grpc.NewMethodDescriptor(method)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating method descriptor")
	}
	return &node.Info{
		Method: md,
	}, nil
}
