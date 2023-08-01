package code

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
)

func New(b *base.Node, node *code.Code) graph.Node {
	switch t := node.Type.(type) {
	case *code.Code_Function:
		return NewFunctionNode(b, t.Function)
	case *code.Code_Server:
		return NewServer(b, t.Server)
	default:
		return nil
	}
}

func NewProto(name string, c *code.Code) *gen.Node {
	return &gen.Node{
		Id:   uuid.NewString(),
		Name: name,
		Type: &gen.Node_Code{
			Code: c,
		},
	}
}
