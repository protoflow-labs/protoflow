package node

import (
	"github.com/protoflow-labs/protoflow/gen"
)

type GRPCNode struct {
	BaseNode
	*gen.GRPC
}

var _ Node = &GRPCNode{}

func NewGRPCNode(node *gen.Node) *GRPCNode {
	return &GRPCNode{
		BaseNode: NewBaseNode(node),
		GRPC:     node.GetGrpc(),
	}
}
