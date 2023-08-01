package edge

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/robertkrimen/otto"
)

type PublishesToListener struct {
	n graph.Node
	e *PublishesTo
}

var _ graph.Listener = (*PublishesToListener)(nil)

func NewPublishesToListener(n graph.Node, e *PublishesTo) *PublishesToListener {
	return &PublishesToListener{
		n: n,
		e: e,
	}
}

func (p *PublishesToListener) GetNode() graph.Node {
	return p.n
}

func (p *PublishesToListener) Transform(ctx context.Context, input graph.IO) (graph.IO, error) {
	if p.e.CodeAdapter != "" {
		obs := input.Observable.Map(func(c context.Context, i any) (any, error) {
			vm := otto.New()
			err := vm.Set("input", i)
			v, err := vm.Run(p.e.CodeAdapter)
			if err != nil {
				return nil, errors.Wrapf(err, "error running code adapter")
			}
			o, err := v.Export()
			if err != nil {
				return nil, errors.Wrapf(err, "error exporting input from vm")
			}
			return o, nil
		})

		return graph.IO{
			Observable: obs,
		}, nil
	}
	return graph.IO{
		Observable: input.Observable,
	}, nil
}
