package code

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/grpc"
	igrpc "github.com/protoflow-labs/protoflow/pkg/grpc"
)

type Server struct {
	*base.Node
	*code.Server
	GRPC *grpc.Server
}

var _ graph.Node = &Server{}
var _ grpc.ServerProvider = &Server{}

func NewServer(b *base.Node, node *code.Server) *Server {
	return &Server{
		Node:   b,
		Server: node,
		// TODO breadchris maybe there is a graph relationship here between the server and the grpc resource
		GRPC: grpc.NewServer(b, node.Grpc),
	}
}

func NewServerProto(runtime code.Runtime) *code.Code {
	return &code.Code{
		Type: &code.Code_Server{
			Server: &code.Server{
				Runtime: runtime,
				// TODO breadchris should there be a default URL?
				Grpc: grpc.NewServerProto("localhost:8000").GetServer(),
			},
		},
	}
}

func (r *Server) GetServer() *grpc.Server {
	return r.GRPC
}

func (r *Server) Type() (*graph.Info, error) {
	return nil, errors.New("implement me")
}

func (r *Server) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	return input, nil
}

func (r *Server) Provide() ([]*gen.Node, error) {
	return igrpc.EnumerateResourceBlocks(r.Grpc, false)
}
