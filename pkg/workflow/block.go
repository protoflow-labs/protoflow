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

func NewNode(node *gen.Node, block *gen.Block) (Node, error) {
	switch block.Type.(type) {
	case *gen.Block_Grpc:
		b := block.GetGrpc()
		n := node.GetGrpc()
		if n != nil {
			if n.Service != "" {
				b.Service = n.Service
			}
			if n.Method != "" {
				b.Method = n.Method
			}
		}
		return &GRPCNode{
			GRPC: b,
		}, nil
	default:
		return nil, errors.New("no node found")
	}
}
