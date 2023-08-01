package data

import (
	"github.com/protoflow-labs/protoflow/gen/data"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
)

func New(b *base.Node, node *data.Data) graph.Node {
	switch t := node.Type.(type) {
	case *data.Data_Input:
		return NewInputNode(b, t.Input)
	case *data.Data_Config:
		return NewConfigNode(b, t.Config)
	default:
		return nil
	}
}
