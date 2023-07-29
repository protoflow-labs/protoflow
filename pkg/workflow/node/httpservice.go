package node

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"html/template"
	"net/url"
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

func (n *RouteNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	routerResource, ok := input.Resource.(*resource.HTTPRouterResource)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting http router resource: %s", n.Route.Path)
	}

	output := make(chan rxgo.Item)
	input.Observable.ForEach(func(item any) {
		r, ok := item.(*gen.HttpRequest)
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

type TemplateNode struct {
	BaseNode
	Template *gen.Template
}

var _ graph.Node = &TemplateNode{}

func NewTemplateNode(node *gen.Node) *TemplateNode {
	return &TemplateNode{
		BaseNode: NewBaseNode(node),
		Template: node.GetTemplate(),
	}
}

func (n *TemplateNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	_, ok := input.Resource.(*resource.TemplateServiceResource)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting template service resource: %s", n.NormalizedName())
	}

	output := make(chan rxgo.Item)

	input.Observable.ForEach(func(item any) {
		tmpl, err := template.New(n.NormalizedName()).Parse(n.Template.Template)
		if err != nil {
			output <- rx.NewError(err)
			return
		}
		b := &bytes.Buffer{}
		err = tmpl.Execute(b, item)
		if err != nil {
			output <- rx.NewError(err)
			return
		}
		resp := &gen.HttpResponse{
			Headers: []*gen.Header{},
			Body:    b.Bytes(),
		}
		output <- rx.NewItem(resp)
	}, func(err error) {
		output <- rx.NewError(err)
	}, func() {
		close(output)
	})

	return graph.Output{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}

//
//type RESTNode struct {
//	BaseNode
//	*gen.REST
//}
//
//var _ graph.Node = &RESTNode{}
//
//func NewRestNode(node *gen.Node) *RESTNode {
//	return &RESTNode{
//		BaseNode: NewBaseNode(node),
//		REST:     node.GetRest(),
//	}
//}
//
//func (n *RESTNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
//	log.Debug().
//		Interface("headers", n.Headers).
//		Str("method", n.Method).
//		Str("path", n.Path).
//		Msgf("executing rest")
//	// TODO breadchris turn this into streamable because why not
//	item, err := input.Observable.First().Get()
//	if err != nil {
//		return graph.Output{}, errors.Wrapf(err, "error getting first item from observable")
//	}
//	res, err := util.InvokeMethodOnUrl(n.Method, n.Path, n.Headers, item.V)
//	if err != nil {
//		return graph.Output{Observable: rxgo.Empty()}, nil
//	}
//	return graph.Output{
//		Observable: rxgo.Just(res)(),
//	}, nil
//}
