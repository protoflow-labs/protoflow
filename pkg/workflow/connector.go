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
	cleanup   []func()
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

	for _, obs := range c.observers {
		_, dispose := obs.Observable.Connect(ctx)
		//obs.Observable.DoOnCompleted(func() {
		//	dispose()
		//})
		c.cleanup = append(c.cleanup, dispose)
	}
	_, _ = o.Connect(ctx)
	return o
}

func (c *Connector) Dispose() {
	for _, c := range c.cleanup {
		c()
	}
}
