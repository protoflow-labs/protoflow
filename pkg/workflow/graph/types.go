package graph

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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

type Node interface {
	NormalizedName() string
	ID() string
	ResourceID() string
	Info(r Resource) (*Info, error)
	Represent() (string, error)
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
