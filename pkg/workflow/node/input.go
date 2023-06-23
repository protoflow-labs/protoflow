package node

import (
	"github.com/protoflow-labs/protoflow/gen"
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
