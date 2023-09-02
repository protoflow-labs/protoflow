package http

import (
	"context"
	"fmt"
	"github.com/protoflow-labs/protoflow/gen/http"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/reactivex/rxgo/v2"
)

type Response struct {
	*base.Node
	*http.Response
}

var _ graph.Node = &Response{}

func NewResponse(b *base.Node, node *http.Response) *Response {
	return &Response{
		Node:     b,
		Response: node,
	}
}

func NewResponseProto() *http.HTTP {
	return &http.HTTP{
		Type: &http.HTTP_Response{
			Response: &http.Response{},
		},
	}
}

func (n *Response) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	output := make(chan rxgo.Item)
	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, err
	}
	routerResource, ok := p.(*Router)
	if !ok {
		return graph.IO{}, fmt.Errorf("error getting http router resource: %v", n.Response)
	}

	input.Observable.ForEach(func(item any) {
		r, ok := item.(*http.Response)
		if !ok {
			output <- rx.NewError(fmt.Errorf("error getting http request from stream"))
			return
		}
		routerResource.HTTPStream.Responses <- rxgo.Of(r)
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
