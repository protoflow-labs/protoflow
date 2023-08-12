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

	io *graph.IO
}

var _ graph.Node = &InputNode{}

type InputOption func(*InputNode) *InputNode

func WithObservable(io *graph.IO) InputOption {
	return func(n *InputNode) *InputNode {
		n.io = io
		return n
	}
}

func NewInputNode(base *base.Node, input *data.Input, ops ...InputOption) *InputNode {
	i := &InputNode{
		Node:  base,
		Input: input,
	}
	for _, op := range ops {
		i = op(i)
	}
	return i
}

func NewInputProto() *data.Data {
	return &data.Data{
		Type: &data.Data_Input{
			Input: &data.Input{},
		},
	}
}

func (n *InputNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	if n.io != nil {
		return *n.io, nil
	}
	return graph.IO{
		Observable: input.Observable,
	}, nil
}
