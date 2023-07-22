package node

import (
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
