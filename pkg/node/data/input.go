package data

import (
	"context"
	"github.com/protoflow-labs/protoflow/gen/data"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

type InputNode struct {
	*base.Node
	*data.Input
}

var _ graph.Node = &InputNode{}

func NewInputNode(base *base.Node, input *data.Input) *InputNode {
	return &InputNode{
		Node:  base,
		Input: input,
	}
}

func (n *InputNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	return graph.Output{
		Observable: input.Observable,
	}, nil
}
