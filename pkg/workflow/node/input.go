package node

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
)

type InputNode struct {
	BaseNode
	*gen.Input
}

var _ Node = &InputNode{}

func NewInputNode(node *gen.Node) *InputNode {
	return &InputNode{
		BaseNode: NewBaseNode(node),
		Input:    node.GetInput(),
	}
}

func (s *InputNode) Execute(executor execute.Executor, input execute.Input) (*execute.Result, error) {
	return &execute.Result{
		Data: input.Params,
	}, nil
}
