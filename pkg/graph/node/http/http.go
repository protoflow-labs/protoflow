package http

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/http"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/reactivex/rxgo/v2"
	"sync"
)

var (
	httpStreamOnce sync.Once
	httpStream     *HTTPEventStream
)

type HTTPEventStream struct {
	Requests   chan rxgo.Item
	Responses  chan *http.Response
	RequestObs rxgo.Observable
}

// TODO breadchris proper dependency injection will need to be figured out to make this work
func NewHTTPEventStream() *HTTPEventStream {
	httpStreamOnce.Do(func() {
		// TODO breadchris there must be an easier way to do this
		// I was thinking of bypassing the need for this altogether and
		// dispatching a workflow job to a workflow service, maybe the executor?
		requestChan := make(chan rxgo.Item)
		responseChan := make(chan *http.Response)
		requestObs := rxgo.FromChannel(requestChan)
		httpStream = &HTTPEventStream{
			Requests:   requestChan,
			Responses:  responseChan,
			RequestObs: requestObs,
		}
	})
	return httpStream
}

type Router struct {
	*base.Node
	*http.Router
	HTTPStream *HTTPEventStream
}

var _ graph.Node = &Router{}

func NewRouterNode(b *base.Node, node *http.Router) *Router {
	return &Router{
		Node:   b,
		Router: node,
	}
}

func (r *Router) Init() (func(), error) {
	// TODO breadchris proper dependency injection will need to be figured out to make this work
	r.HTTPStream = NewHTTPEventStream()
	return nil, nil
}

func (r *Router) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	return graph.IO{}, errors.New("cannot wire router node")
}

//
//type RESTNode struct {
//	Node
//	*gen.REST
//}
//
//var _ graph.Node = &RESTNode{}
//
//func NewRestNode(node *gen.Node) *RESTNode {
//	return &RESTNode{
//		Node: NodeFromProto(node),
//		REST:     node.GetRest(),
//	}
//}
//
//func (n *RESTNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
//	log.Debug().
//		Interface("headers", n.Headers).
//		Str("method", n.Method).
//		Str("path", n.Path).
//		Msgf("executing rest")
//	// TODO breadchris turn this into streamable because why not
//	item, err := input.Observable.First().Get()
//	if err != nil {
//		return graph.IO{}, errors.Wrapf(err, "error getting first item from observable")
//	}
//	res, err := util.InvokeMethodOnUrl(n.Method, n.Path, n.Headers, item.V)
//	if err != nil {
//		return graph.IO{Observable: rxgo.Empty()}, nil
//	}
//	return graph.IO{
//		Observable: rxgo.Just(res)(),
//	}, nil
//}
