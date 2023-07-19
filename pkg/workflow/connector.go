package workflow

import (
	"context"
	"github.com/reactivex/rxgo/v2"
)

type Connector struct {
	observers []rxgo.Observable
}

func NewConnector() *Connector {
	return &Connector{
		observers: []rxgo.Observable{},
	}
}

func (c *Connector) Add(o rxgo.Observable) {
	c.observers = append(c.observers, o)
}

func (c *Connector) Connect(ctx context.Context) rxgo.Observable {
	o := rxgo.Merge(c.observers, rxgo.WithPublishStrategy())

	// TODO breadchris figure out what to do with disposed and cancel
	// disposed, cancel := output.Observable.Connect(ctx)

	for _, obs := range c.observers {
		obs.Connect(ctx)
	}
	o.Connect(ctx)
	return o
}
