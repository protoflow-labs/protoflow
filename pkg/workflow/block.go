package workflow

import (
	"github.com/pkg/errors"
	protoflow "github.com/protoflow-labs/protoflow/gen"
)

type Block interface {
	Execute(executor Executor, input Input) (*Result, error)
}

type GRPCBlock struct {
	*protoflow.GRPC
}

type RESTBlock struct {
	*protoflow.REST
}

type EntityBlock struct {
	protoflow.Entity
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

func NewBlock(block *protoflow.Block) (Block, error) {
	switch block.Type.(type) {
	case *protoflow.Block_Grpc:
		g := block.GetGrpc()
		return &GRPCBlock{
			GRPC: g,
		}, nil
	case *protoflow.Block_Rest:
		r := block.GetRest()
		return &RESTBlock{
			REST: r,
		}, nil
	default:
		return nil, errors.New("no block found")
	}
}
