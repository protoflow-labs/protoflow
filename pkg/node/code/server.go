package code

import (
	"context"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/node/base"
	"github.com/protoflow-labs/protoflow/pkg/node/grpc"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"time"
)

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

func NewServerProto(runtime code.Runtime) *code.Code {
	return &code.Code{
		Type: &code.Code_Server{
			Server: &code.Server{
				Runtime: runtime,
				Grpc:    grpc.NewServerProto("localhost:8080").GetServer(),
			},
		},
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
