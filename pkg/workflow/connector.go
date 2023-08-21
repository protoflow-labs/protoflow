package workflow

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/reactivex/rxgo/v2"
)

type Connector struct {
	observers map[string]*graph.IO
}

func NewConnector() *Connector {
	return &Connector{
		observers: map[string]*graph.IO{},
	}
}

func (c *Connector) Get(nodeID string) (*graph.IO, bool) {
	o, ok := c.observers[nodeID]
	return o, ok
}

func (c *Connector) Add(nodeID string, o *graph.IO) {
	c.observers[nodeID] = o
}

func (c *Connector) Connect(ctx context.Context) rxgo.Observable {
	var obsSlice []rxgo.Observable
	for nodeID, o := range c.observers {
		obsSlice = append(obsSlice, o.Observable.Map(func(ctx context.Context, i any) (any, error) {
			output, err := json.Marshal(i)
			if err != nil {
				return nil, errors.Wrapf(err, "error marshalling output")
			}
			return &gen.NodeExecution{
				NodeId: nodeID,
				Output: string(output),
			}, nil
		}))
	}
	o := rxgo.Merge(obsSlice, rxgo.WithPublishStrategy())

	// TODO breadchris figure out what to do with disposed and cancel
	// disposed, cancel := output.RequestObs.Connect(ctx)

	for _, obs := range c.observers {
		obs.Observable.Connect(ctx)
	}
	o.Connect(ctx)
	return o
}
