package node

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

type SecretNode struct {
	BaseNode
	Secret *gen.Secret
}

var _ graph.Node = &SecretNode{}

func NewSecretNode(node *gen.Node) *SecretNode {
	return &SecretNode{
		BaseNode: NewBaseNode(node),
		Secret:   node.GetSecret(),
	}
}

func (n *SecretNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	return graph.Output{}, errors.New("not implemented")
}
