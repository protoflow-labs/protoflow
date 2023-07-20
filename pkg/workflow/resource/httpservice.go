package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/reactivex/rxgo/v2"
	"path"
	"sync"
)

var (
	httpStreamOnce sync.Once
	httpStream     *HTTPEventStream
)

type HTTPEventStream struct {
	Requests   chan rxgo.Item
	Responses  chan *gen.HttpResponse
	RequestObs rxgo.Observable
}

// TODO breadchris proper dependency injection will need to be figured out to make this work
func NewHTTPEventStream() *HTTPEventStream {
	httpStreamOnce.Do(func() {
		// TODO breadchris there must be an easier way to do this
		// I was thinking of bypassing the need for this altogether and
		// dispatching a workflow job to a workflow service, maybe the executor?
		requestChan := make(chan rxgo.Item)
		responseChan := make(chan *gen.HttpResponse)
		requestObs := rxgo.FromChannel(requestChan)
		httpStream = &HTTPEventStream{
			Requests:   requestChan,
			Responses:  responseChan,
			RequestObs: requestObs,
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
	// TODO breadchris proper dependency injection will need to be figured out to make this work
	r.HTTPStream = NewHTTPEventStream()
	return nil, nil
}

func (r *HTTPRouterResource) Path(n *node.RouteNode) string {
	return path.Join(r.HTTPRouter.Root, n.Route.Path)
}

func (r *HTTPRouterResource) Info(n node.Node) (*node.Info, error) {
	_, ok := n.(*node.RouteNode)
	if !ok {
		return nil, errors.New("node is not a route node")
	}
	return nil, nil
}
