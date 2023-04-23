package api

import (
	"fmt"
	"net"

	"github.com/protoflow-labs/protoflow/pkg/workflow"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewGRPCServer(manager *workflow.TemporalManager) *GRPCServer {
	s := grpc.NewServer()

	reflection.Register(s)
	return &GRPCServer{
		server: s,
	}
}

func (s *GRPCServer) Serve(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return s.server.Serve(listener)
}
