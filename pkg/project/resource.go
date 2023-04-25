package project

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	grpcanal "github.com/protoflow-labs/protoflow/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func EnumerateResource(resource *gen.Resource) (*gen.Resource, error) {
	resource.Id = uuid.New().String()

	var (
		blocks []*gen.Block
		err    error
	)
	switch resource.Type.(type) {
	case *gen.Resource_GrpcService:
		g := resource.GetGrpcService()
		blocks, err = blocksFromGRPC(g)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating blocks for resource %s", resource.Id)
		}
	}

	resource.Blocks = blocks
	return resource, nil
}

func blocksFromGRPC(service *gen.GRPCService) ([]*gen.Block, error) {
	if service.Host == "" {
		return nil, errors.New("host is required")
	}

	conn, err := grpc.Dial(service.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to python server at %s", service.Host)
	}

	methodDesc, err := grpcanal.AllMethodsViaReflection(context.Background(), conn)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get all methods via reflection")
	}

	var blocks []*gen.Block
	for _, m := range methodDesc {
		serviceName := m.GetService().GetName()
		methodName := m.GetName()
		blocks = append(blocks, &gen.Block{
			Id:   uuid.New().String(),
			Name: serviceName + "." + methodName,
			Type: &gen.Block_Grpc{
				Grpc: &gen.GRPC{
					Service: serviceName,
					Method:  methodName,
				},
			},
		})
	}
	return blocks, nil
}