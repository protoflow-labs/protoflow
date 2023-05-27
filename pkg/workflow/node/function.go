package node

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"strings"
)

type FunctionNode struct {
	BaseNode
	Function *gen.Function
}

var _ Node = &FunctionNode{}

func NewFunctionNode(node *gen.Node) *FunctionNode {
	return &FunctionNode{
		BaseNode: NewBaseNode(node),
		Function: node.GetFunction(),
	}
}

func (f *FunctionNode) Execute(executor execute.Executor, input execute.Input) (*execute.Result, error) {
	return executor.Execute(activity.ExecuteFunctionNode, f, input)
}

func (f *FunctionNode) Info(res resource.Resource) (*Info, error) {
	r, ok := res.(*resource.LanguageServiceResource)
	if !ok {
		return nil, errors.Errorf("resource is not a language resource")
	}
	grpcNode := f.ToGRPC(r)

	// TODO breadchris we should know where the function node is located and should read/write from the proto
	return grpcNode.Info(r.GRPCResource)
}

func (f *FunctionNode) ToGRPC(res *resource.LanguageServiceResource) *GRPCNode {
	serviceName := strings.ToLower(res.Runtime.String()) + "Service"
	return &GRPCNode{
		GRPC: &gen.GRPC{
			Package: "protoflow",
			Service: serviceName,
			Method:  f.Name,
		},
	}
}
