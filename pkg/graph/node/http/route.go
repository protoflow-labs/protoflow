package http

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/http"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
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

func NewRouteProto() *http.HTTP {
	return &http.HTTP{
		Type: &http.HTTP_Route{
			Route: &http.Route{},
		},
	}
}

func (n *RouteNode) Path(r *Router) string {
	return path.Join(r.Root, n.Route.Path)
}

func (n *RouteNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, err
	}
	routerResource, ok := p.(*Router)
	if !ok {
		return graph.IO{}, fmt.Errorf("error getting http router resource: %s", n.Route.Path)
	}

	// TODO breadchris log what path is being listened to

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
			log.Debug().Msgf("skipping request %s != %s || %s != %s", u.Path, n.Path(routerResource), r.Method, n.Route.Method)
			return
		}
		output <- rx.NewItem(r)
	}, func(err error) {
		output <- rx.NewError(err)
	}, func() {
		close(output)
	})

	return graph.IO{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}
