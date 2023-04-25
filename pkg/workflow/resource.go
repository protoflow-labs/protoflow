package workflow

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ResourceFromProto(r *gen.Resource) (Resource, error) {
	switch t := r.Type.(type) {
	case *gen.Resource_GrpcService:
		g := r.GetGrpcService()
		return &GRPCResource{
			GRPCService: g,
		}, nil
	default:
		return nil, fmt.Errorf("no resource found with type: %s", t)
	}
}

type Resource interface {
	New() (any, error)
}

type GRPCResource struct {
	*gen.GRPCService
}

func (r *GRPCResource) New() (interface{}, error) {
	conn, err := grpc.Dial(r.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to grpc server at %s", r.Host)
	}
	return conn, nil
}
