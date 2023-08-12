package workflow

import (
	"context"
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
	for _, o := range c.observers {
		obsSlice = append(obsSlice, o.Observable)
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
