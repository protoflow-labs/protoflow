package http

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/http"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/reactivex/rxgo/v2"
	"net/url"
	"path"
)

type RouteNode struct {
	*base.Node
	*http.Route
}

var _ graph.Node = &RouteNode{}

func NewRouteNode(b *base.Node, node *http.Route) *RouteNode {
	return &RouteNode{
		Node:  b,
		Route: node,
	}
}

func (n *RouteNode) Path(r *Router) string {
	return path.Join(r.Root, n.Route.Path)
}

func (n *RouteNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	p, err := n.Provider()
	if err != nil {
		return graph.Output{}, err
	}
	routerResource, ok := p.(*Router)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting http router resource: %s", n.Route.Path)
	}

	output := make(chan rxgo.Item)
	input.Observable.ForEach(func(item any) {
		r, ok := item.(*http.Request)
		if !ok {
			output <- rx.NewError(fmt.Errorf("error getting http request from stream"))
			return
		}
		u, err := url.Parse(r.Url)
		if err != nil {
			output <- rx.NewError(errors.Wrapf(err, "error parsing request url"))
			return
		}
		if u.Path != n.Path(routerResource) || r.Method != n.Route.Method {
			return
		}
		output <- rx.NewItem(r)
	}, func(err error) {
		output <- rx.NewError(err)
	}, func() {
		close(output)
	})

	return graph.Output{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}
