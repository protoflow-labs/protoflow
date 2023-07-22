package node

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/rs/zerolog/log"
	"strings"
)

type FunctionNode struct {
	BaseNode
	Function *gen.Function
}

var _ graph.Node = &FunctionNode{}

func NewFunctionNode(node *gen.Node) *FunctionNode {
	return &FunctionNode{
		BaseNode: NewBaseNode(node),
		Function: node.GetFunction(),
	}
}

func (n *FunctionNode) Wire(ctx context.Context, input graph.Input) (graph.Output, error) {
	log.Debug().
		Str("name", n.NormalizedName()).
		Msg("setting up function")
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
