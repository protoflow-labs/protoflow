package workflow

import (
	"github.com/pkg/errors"
	protoflow "github.com/protoflow-labs/protoflow/gen"
)

type Block interface {
	Execute(executor Executor, input Input) (*Result, error)
}

type GRPCBlock struct {
	protoflow.GRPC
}

type RESTBlock struct {
	protoflow.REST
}

var activity = &Activity{}

func (s *GRPCBlock) Init() error {
	return nil
}

func (s *GRPCBlock) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteGRPCBlock, s, input)
}

func (s *RESTBlock) Execute(executor Executor, input Input) (*Result, error) {
	return executor.Execute(activity.ExecuteRestBlock, s, input)
}

func NewBlock(node *protoflow.Node) (Block, error) {
	switch node.Block.Type.(type) {
	case *protoflow.Block_Grpc:
		g := node.Block.GetGrpc()
		return &GRPCBlock{
			GRPC: *g,
		}, nil
	case *protoflow.Block_Rest:
		r := node.Block.GetRest()
		return &RESTBlock{
			REST: *r,
		}, nil
	default:
		return nil, errors.New("no block found")
	}
	return nil, nil
}
