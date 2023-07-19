package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/reactivex/rxgo/v2"
	"sync"
)

var (
	httpStreamOnce sync.Once
	httpStream     *HTTPEventStream
)

type HTTPEventStream struct {
	Items      chan<- rxgo.Item
	Observable rxgo.Observable
}

// TODO breadchris proper dependency injection will need to be figured out to make this work
func NewHTTPEventStream() *HTTPEventStream {
	httpStreamOnce.Do(func() {
		httpChan := make(chan rxgo.Item)
		httpObs := rxgo.FromEventSource(httpChan)
		httpStream = &HTTPEventStream{
			Items:      httpChan,
			Observable: httpObs,
		}
	})
	return httpStream
}

type HTTPRouterResource struct {
	*BaseResource
	*gen.HTTPRouter
	HTTPStream *HTTPEventStream
}

func (r *HTTPRouterResource) Init() (func(), error) {
	r.HTTPStream = NewHTTPEventStream()
	return nil, nil
}

func (r *HTTPRouterResource) Info(n node.Node) (*node.Info, error) {
	_, ok := n.(*node.RouteNode)
	if !ok {
		return nil, errors.New("node is not a route node")
	}
	return nil, nil
}
