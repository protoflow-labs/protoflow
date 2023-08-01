package http

import (
	"bytes"
	"context"
	"github.com/protoflow-labs/protoflow/gen/http"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/reactivex/rxgo/v2"
	"html/template"
)

type TemplateNode struct {
	*base.Node
	*http.Template
}

var _ graph.Node = &TemplateNode{}

func NewTemplateNode(b *base.Node, node *http.Template) *TemplateNode {
	return &TemplateNode{
		Node:     b,
		Template: node,
	}
}

func (n *TemplateNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
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
		resp := &http.Response{
			Headers: []*http.Header{},
			Body:    b.Bytes(),
		}
		output <- rx.NewItem(resp)
	}, func(err error) {
		output <- rx.NewError(err)
	}, func() {
		close(output)
	})

	return graph.IO{
		Observable: rxgo.FromChannel(output, rxgo.WithPublishStrategy()),
	}, nil
}
