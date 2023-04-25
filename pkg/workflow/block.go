package workflow

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
)

type Node interface {
	Execute(executor Executor, input Input) (*Result, error)
}

type GRPCNode struct {
	*gen.GRPC
}

type RESTNode struct {
	*gen.REST
}

type EntityBlock struct {
	gen.Entity
}

var activity = &Activity{}

func (s *GRPCNode) Init() error {
	return nil
}

func (s *GRPCNode) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteGRPCNode, s, input)
}

func (s *RESTNode) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteRestNode, s, input)
}

func NewNode(node *gen.Node) (Node, error) {
	switch node.Config.(type) {
	case *gen.Node_Grpc:
		g := node.GetGrpc()
		return &GRPCNode{
			GRPC: g,
		}, nil
	case *gen.Node_Rest:
		r := node.GetRest()
		return &RESTNode{
			REST: r,
		}, nil
	default:
		return nil, errors.New("no node found")
	}
}
