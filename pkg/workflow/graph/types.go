package graph

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Input struct {
	Observable rxgo.Observable
	Resource   Resource
}

type Output struct {
	Observable rxgo.Observable
}

type DependencyProvider map[string]Resource

type Info struct {
	Method *grpc.MethodDescriptor
}

func (s *Info) BuildProto() (string, error) {
	s.Method.MethodDesc.ParentFile().Package()
	svc := s.Method.MethodDesc.Parent().(protoreflect.ServiceDescriptor)
	pkgName := string(svc.ParentFile().Package())
	svcName := string(svc.Name())
	methodProto, err := manager.GetProtoForMethod(pkgName, svcName, s.Method.MethodDesc)
	if err != nil {
		return "", errors.Wrapf(err, "error getting proto for method %s", s.Method.MethodDesc.Name())
	}
	return methodProto, nil
}

type IOFunc func(ctx context.Context, input Input) (Output, error)

type Node interface {
	NormalizedName() string
	ID() string
	ResourceID() string
	Info(r Resource) (*Info, error)
	// Represent the node as a string
	Represent() (string, error)
	// Wire up the node to an input stream of data and return an output stream of data
	Wire(ctx context.Context, input Input) (Output, error)
}

type Edge struct {
	From Node
	To   Node
}

type Resource interface {
	Name() string
	Init() (func(), error)
	ID() string
	AddNode(n Node)
	Nodes() []Node
	ResolveDependencies(dp DependencyProvider) error
	//DeploymentInfo() (*DeploymentInfo, error)
}
