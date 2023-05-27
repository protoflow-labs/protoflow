package project

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
)

func getDefaultProject(name string, bucketDir string) gen.Project {
	return gen.Project{
		Id:   uuid.NewString(),
		Name: name,
		Resources: []*gen.Resource{
			{
				Id:   uuid.NewString(),
				Name: "protoflow",
				Type: &gen.Resource_GrpcService{
					GrpcService: &gen.GRPCService{
						Host: "localhost:8080",
					},
				},
			},
			{
				Id:   uuid.NewString(),
				Name: "js",
				Type: &gen.Resource_LanguageService{
					LanguageService: &gen.LanguageService{
						Runtime: gen.Runtime_NODEJS,
						Grpc: &gen.GRPCService{
							Host: "localhost:8086",
						},
					},
				},
			},
			{
				Id:   uuid.NewString(),
				Name: "docs",
				Type: &gen.Resource_Docstore{
					Docstore: &gen.Docstore{
						Url: "mem://",
					},
				},
			},
			{
				Id:   uuid.NewString(),
				Name: "bucket",
				Type: &gen.Resource_Blobstore{
					Blobstore: &gen.Blobstore{
						Url: "file://" + bucketDir,
					},
				},
			},
		},
	}
}
