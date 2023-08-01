package data

import (
	"context"
	"github.com/protoflow-labs/protoflow/gen/data"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
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

func (n *InputNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	return graph.IO{
		Observable: input.Observable,
	}, nil
}
