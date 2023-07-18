package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
)

type SecretStoreResource struct {
	*BaseResource
	*gen.SecretStore
}

func (r *SecretStoreResource) Init() (func(), error) {
	return nil, nil
}

func (r *SecretStoreResource) Info(n node.Node) (*node.Info, error) {
	_, ok := n.(*node.SecretNode)
	if !ok {
		return nil, errors.New("node is not a prompt node")
	}
	return nil, nil
}
