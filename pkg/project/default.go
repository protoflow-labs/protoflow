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
			// TODO breadchris programmatically add resources such as language services
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
				Name: "doc store",
				Type: &gen.Resource_DocStore{
					DocStore: &gen.DocStore{
						Url: "mem://",
					},
				},
			},
			{
				Id:   uuid.NewString(),
				Name: "file store",
				Type: &gen.Resource_FileStore{
					FileStore: &gen.FileStore{
						Url: "file://" + bucketDir,
					},
				},
			},
			{
				Id:   uuid.NewString(),
				Name: "openai",
				Type: &gen.Resource_ReasoningEngine{
					ReasoningEngine: &gen.ReasoningEngine{},
				},
			},
			{
				Id:   uuid.NewString(),
				Name: "config provider",
				Type: &gen.Resource_ConfigProvider{
					ConfigProvider: &gen.ConfigProvider{},
				},
			},
		},
	}
}
