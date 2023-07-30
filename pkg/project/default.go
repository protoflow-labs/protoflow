package project

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/node/grpc"
)

func getDefaultProject(name string, bucketDir string) *gen.Project {
	return &gen.Project{
		Id:   uuid.NewString(),
		Name: name,
		Graph: &gen.Graph{
			Nodes: []*gen.Node{
				{
					Id:   uuid.NewString(),
					Name: "protoflow",
					Type: &gen.Node_Grpc{
						Grpc: grpc.NewServerProto("localhost:8080"),
					},
				},
			},
		},
	}
}
