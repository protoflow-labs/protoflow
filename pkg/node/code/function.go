package code

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/code"
	pgrpc "github.com/protoflow-labs/protoflow/gen/grpc"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/node/grpc"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"strings"
	"time"
)

type FunctionNode struct {
	*base.Node
	*code.Function
	inMemory bool
	f        graph.IOFunc
}

var _ graph.Node = &FunctionNode{}

type FunctionNodeOption func(*FunctionNode) *FunctionNode

func InMemoryObserver(name string) graph.IOFunc {
	return func(ctx context.Context, input graph.Input) (graph.Output, error) {
		output := make(chan rxgo.Item)
		input.Observable.ForEach(func(item any) {
			log.Info().
				Str("name", name).
				Interface("item", item).
				Msg("observing item")
			output <- rxgo.Of(item)
			close(output)
		}, func(err error) {
			log.Info().Str("name", name).Err(err).Msg("err")
		}, func() {
			log.Info().Str("name", name).Msg("complete")
		})
		return graph.Output{
			Observable: rxgo.FromChannel(output),
		}, nil
	}
}

func WithFunction(f graph.IOFunc) FunctionNodeOption {
	return func(n *FunctionNode) *FunctionNode {
		n.inMemory = true
		n.f = f
		return n
	}
}

func NewFunctionNode(b *base.Node, node *code.Function, ops ...FunctionNodeOption) *FunctionNode {
	f := &FunctionNode{
		Node:     b,
		Function: node,
	}
	for _, op := range ops {
		f = op(f)
	}
	return f
}

func (n *FunctionNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	log.Debug().
		Str("name", n.NormalizedName()).
		Msg("setting up function")

	if n.inMemory {
		return n.f(ctx, input)
	}

	p, err := n.Provider()
	if err != nil {
		return graph.Output{}, err
	}

	g, ok := p.(*Server)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting language service resource: %s", n.Name)
	}

	grpcNode := n.ToGRPC(g)
	return grpcNode.Wire(ctx, input)
}

func (n *FunctionNode) Info() (*graph.Info, error) {
	p, err := n.Provider()
	if err != nil {
		return nil, err
	}
	ls, ok := p.(*Server)
	if ok {
		return nil, errors.New("language service resource is not supported")
	}
	grpcNode := n.ToGRPC(ls)
	// TODO breadchris we should know where the function node is located and should read/write from the proto
	return grpcNode.Info()
}

func (n *FunctionNode) ToGRPC(r *Server) *grpc.Method {
	serviceName := strings.ToLower(r.Runtime.String()) + "Service"
	return grpc.NewMethod(n.Node, &pgrpc.Method{
		Package: "protoflow",
		Service: serviceName,
		Method:  n.Name,
	})
}

type Server struct {
	*base.Node
	*code.Server
	GRPC *grpc.Server
}

var _ graph.Node = &Server{}

func NewServer(b *base.Node, node *code.Server) *Server {
	return &Server{
		Node:   b,
		Server: node,
		// TODO breadchris maybe there is a graph relationship here between the server and the grpc resource
		GRPC: grpc.NewServer(b, node.Grpc),
	}
}

func (r *Server) Init() (func(), error) {
	if r.Grpc != nil {
		return r.GRPC.Init()
	}
	return nil, nil
}

func (r *Server) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	//TODO implement me
	panic("implement me")
}

func ensureRunning(host string) error {
	maxRetries := 1
	retryInterval := 2 * time.Second

	u, err := url.Parse(host)
	if err != nil {
		return errors.Wrapf(err, "unable to parse url %s", host)
	}

	log.Debug().Str("host", host).Msg("waiting for host to come online")
	for i := 1; i <= maxRetries; i++ {
		conn, err := net.DialTimeout("tcp", u.Host, time.Second)
		if err == nil {
			conn.Close()
			log.Debug().Str("host", host).Msg("host is not listening")
			return nil
		} else {
			log.Debug().Err(err).Int("attempt", i).Int("max", maxRetries).Msg("error connecting to host")
			time.Sleep(retryInterval)
		}
	}
	return errors.New("host did not come online in time")
}
