package code

import (
	"github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
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
