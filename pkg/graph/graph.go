package graph

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/grpc/manager"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IO struct {
	Observable rxgo.Observable
	Cleanup    func()
}

type DependencyProvider map[string]Node

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

type IOFunc func(ctx context.Context, input IO) (IO, error)

// TODO breadchris better name?
type Listener interface {
	GetNode() Node
	Transform(ctx context.Context, input *IO) (*IO, error)
}

type Node interface {
	ID() string
	Name() string
	NormalizedName() string
	// TODO breadchris type should probably just return a message descriptor
	Type() (*Info, error)

	Provide() ([]*gen.Node, error)

	// Provider returns the node that this node depends on. (eg. a grpc method node will return the service node)
	Provider() (Node, error)
	SetProvider(n Node) error

	// Dependents returns the nodes that depend on this node.
	Dependents() []Node
	AddDependent(n Node)

	// Wire up the node to an input stream of data and return an output stream of data
	Wire(ctx context.Context, input IO) (IO, error)

	// Subscribers returns the nodes that subscribe to this node.
	Subscribers() []Listener
	AddSubscriber(n Listener)

	// Publishers returns the nodes that this node subscribes to.
	Publishers() []Listener
	AddPublishers(n Listener)

	//DeploymentInfo() (*DeploymentInfo, error)
}

type Edge interface {
	ID() string
	From() string
	To() string
	Connect(from, to Node) error
}
