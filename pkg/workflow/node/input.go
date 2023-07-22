package node

import (
	"context"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

type InputNode struct {
	BaseNode
	*gen.Input
}

var _ graph.Node = &InputNode{}

func NewInputNode(node *gen.Node) *InputNode {
	return &InputNode{
		BaseNode: NewBaseNode(node),
		Input:    node.GetInput(),
	}
}

func (n *InputNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	return graph.Output{
		Observable: input.Observable,
	}, nil
}
