package edge

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/robertkrimen/otto"
)

type MapListener struct {
	n graph.Node
	e *Map
}

var _ graph.Listener = (*MapListener)(nil)

func NewMapListener(n graph.Node, e *Map) *MapListener {
	return &MapListener{
		n: n,
		e: e,
	}
}

func (p *MapListener) GetNode() graph.Node {
	return p.n
}

func (p *MapListener) Transform(ctx context.Context, input *graph.IO) (*graph.IO, error) {
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

		return &graph.IO{
			Observable: obs,
		}, nil
	}
	return &graph.IO{
		Observable: input.Observable,
	}, nil
}
