package grpc

import (
	"context"
	"github.com/protoflow-labs/protoflow/gen"
	pgrpc "github.com/protoflow-labs/protoflow/gen/grpc"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"strings"
)

type Server struct {
	*base.Node
	*pgrpc.Server
}

type ServerProvider interface {
	GetServer() *Server
}

func NewServer(b *base.Node, n *pgrpc.Server) *Server {
	// TODO breadchris this shouldn't be here
	if !strings.HasPrefix(n.Host, "http://") {
		n.Host = "http://" + n.Host
	}
	return &Server{
		Node:   b,
		Server: n,
	}
}

func NewServerProto(host string) *pgrpc.GRPC {
	return &pgrpc.GRPC{
		Type: &pgrpc.GRPC_Server{
			Server: &pgrpc.Server{
				Host: host,
			},
		},
	}
}

func (n *Server) GetServer() *Server {
	return n
}

func (n *Server) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	return input, nil
}

func (n *Server) Provide() ([]*gen.Node, error) {
	return grpc.EnumerateResourceBlocks(n.Server, false)
}
