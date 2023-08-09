package grpc

import (
	"context"
	pgrpc "github.com/protoflow-labs/protoflow/gen/grpc"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/base"
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

func (n *Server) Init() (func(), error) {
	// TODO breadchris this is a hack to get the grpc server running, this is not ideal
	if !strings.HasPrefix(n.Host, "http://") {
		n.Host = "http://" + n.Host
	}
	//if err := ensureRunning(r.Host); err != nil {
	//	// TODO breadchris ignore errors for now
	//	// return nil, errors.Wrapf(err, "unable to get the %s grpc server running", r.Name())
	//	return nil, nil
	//}
	return nil, nil
}

func (n *Server) Wire(ctx context.Context, input graph.IO) (graph.IO, error) {
	//TODO implement me
	panic("implement me")
}
