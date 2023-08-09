package code

import (
	"context"
	"fmt"
	"github.com/protoflow-labs/protoflow/gen/code"
	pgrpc "github.com/protoflow-labs/protoflow/gen/grpc"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/grpc"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
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
	return func(ctx context.Context, input graph.IO) (graph.IO, error) {
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
		return graph.IO{
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

func NewFunctionProto() *code.Code {
	return &code.Code{
		Type: &code.Code_Function{
			Function: &code.Function{},
		},
	}
}

func (n *FunctionNode) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	log.Debug().
		Str("name", n.NormalizedName()).
		Msg("setting up function")

	if n.inMemory {
		return n.f(ctx, input)
	}

	p, err := n.Provider()
	if err != nil {
		return graph.IO{}, err
	}

	g, ok := p.(*Server)
	if !ok {
		return graph.IO{}, fmt.Errorf("error getting language service resource: %s", n.Name)
	}

	grpcNode := n.ToGRPC(g)
	return grpc.WireMethod(ctx, g.GRPC, grpcNode, input.Observable)
}

func (n *FunctionNode) Type() (*graph.Info, error) {
	p, err := n.Provider()
	if err != nil {
		return nil, err
	}
	ls, ok := p.(*Server)
	if !ok {
		return nil, fmt.Errorf("error getting language server provider: %s", n.Name)
	}
	grpcNode := n.ToGRPC(ls)
	// TODO breadchris we should know where the function node is located and should read/write from the proto
	md, err := grpc.GetMethodDescriptor(ls.GRPC, grpcNode)
	if err != nil {
		return nil, err
	}
	return &graph.Info{
		Method: md,
	}, nil
}

func (n *FunctionNode) ToGRPC(r *Server) *grpc.Method {
	// TODO breadchris fix
	//serviceName := strings.ToLower(r.Runtime.String()) + "Service"
	return grpc.NewMethod(n.Node, &pgrpc.Method{
		Package: "project",
		Service: "ProjectService",
		Method:  n.Name,
	})
}
