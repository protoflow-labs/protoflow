package node

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"strings"
)

type FunctionNode struct {
	BaseNode
	Function *gen.Function
	inMemory bool
	f        graph.IOFunc
}

var _ graph.Node = &FunctionNode{}

func NewFunctionProto(name, resourceID string) *gen.Node {
	return &gen.Node{
		Id:         uuid.NewString(),
		Name:       name,
		ResourceId: resourceID,
		Config: &gen.Node_Function{
			Function: &gen.Function{},
		},
	}
}

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

func NewFunctionNode(node *gen.Node, ops ...FunctionNodeOption) *FunctionNode {
	f := &FunctionNode{
		BaseNode: NewBaseNode(node),
		Function: node.GetFunction(),
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

	g, ok := input.Resource.(*resource.LanguageServiceResource)
	if !ok {
		return graph.Output{}, fmt.Errorf("error getting language service resource: %s", n.Name)
	}

	// provide the grpc resource to the grpc gn call. Is this the best place for this? Should this be provided on injection? Probably.
	input.Resource = g.GRPCResource

	grpcNode := n.ToGRPC(g)
	return grpcNode.Wire(ctx, input)
}

func (n *FunctionNode) Info(r graph.Resource) (*graph.Info, error) {
	ls, ok := r.(*resource.LanguageServiceResource)
	if ok {
		return nil, errors.New("language service resource is not supported")
	}
	grpcNode := n.ToGRPC(ls)
	// TODO breadchris we should know where the function node is located and should read/write from the proto
	return grpcNode.Info(ls.GRPCResource)
}

func (n *FunctionNode) ToGRPC(r *resource.LanguageServiceResource) *GRPCNode {
	serviceName := strings.ToLower(r.Runtime.String()) + "Service"
	return &GRPCNode{
		BaseNode: n.BaseNode,
		GRPC: &gen.GRPC{
			Package: "protoflow",
			Service: serviceName,
			Method:  n.Name,
		},
	}
}
