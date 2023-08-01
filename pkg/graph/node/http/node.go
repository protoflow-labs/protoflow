package http

import (
	"github.com/protoflow-labs/protoflow/gen/http"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
)

func New(b *base.Node, node *http.HTTP) graph.Node {
	switch t := node.Type.(type) {
	case *http.HTTP_Route:
		return NewRouteNode(b, t.Route)
	case *http.HTTP_Router:
		return NewRouterNode(b, t.Router)
	case *http.HTTP_Template:
		return NewTemplateNode(b, t.Template)
	default:
		return nil
	}
}
