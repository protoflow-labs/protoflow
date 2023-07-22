package resource

import (
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/reactivex/rxgo/v2"
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
