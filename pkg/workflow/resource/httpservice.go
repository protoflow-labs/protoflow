package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
)

type HTTPRouterResource struct {
	*BaseResource
	*gen.HTTPRouter
}

func (r *HTTPRouterResource) Init() (func(), error) {
	return nil, nil
}

func (r *HTTPRouterResource) Info(n node.Node) (*node.Info, error) {
	_, ok := n.(*node.RouteNode)
	if !ok {
		return nil, errors.New("node is not a route node")
	}
	return nil, nil
}
