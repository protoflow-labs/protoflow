package grpc

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/grpc"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
)

func New(b *base.Node, node *grpc.GRPC) graph.Node {
	switch t := node.Type.(type) {
	case *grpc.GRPC_Server:
		return NewServer(b, t.Server)
	case *grpc.GRPC_Method:
		return NewMethod(b, t.Method)
	default:
		return nil
	}
}

func NewProto(name string, g *grpc.GRPC) *gen.Node {
	return &gen.Node{
		Id:   uuid.NewString(),
		Name: name,
		Type: &gen.Node_Grpc{
			Grpc: g,
		},
	}
}
