package node

import (
	"github.com/protoflow-labs/protoflow/gen"
)

type FunctionNode struct {
	BaseNode
	Function *gen.Function
}

var _ Node = &FunctionNode{}

func NewFunctionNode(node *gen.Node) *FunctionNode {
	return &FunctionNode{
		BaseNode: NewBaseNode(node),
		Function: node.GetFunction(),
	}
}
