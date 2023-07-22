package node

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"path"
)

type RouteNode struct {
	BaseNode
	Route *gen.Route
}

var _ graph.Node = &RouteNode{}

func NewRouteNode(node *gen.Node) *RouteNode {
	return &RouteNode{
		BaseNode: NewBaseNode(node),
		Route:    node.GetRoute(),
	}
}

func (n *RouteNode) Path(r *resource.HTTPRouterResource) string {
	return path.Join(r.HTTPRouter.Root, n.Route.Path)
}
