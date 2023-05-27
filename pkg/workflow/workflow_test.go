package workflow

import (
	"context"
	"github.com/protoflow-labs/protoflow/pkg/workflow/execute"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"testing"

	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
)

func TestRun(t *testing.T) {
	// TODO breadchris start server to listen for localhost:8080?
	nodeID := "1"
	lr := &gen.Resource{
		Id: "2",
		Type: &gen.Resource_LanguageService{
			LanguageService: &gen.LanguageService{
				Runtime: gen.Runtime_NODEJS,
				Grpc: &gen.GRPCService{
					Host: "localhost:8086",
				},
			},
		},
	}
	r := &gen.Resource{
		Id: nodeID,
		Type: &gen.Resource_GrpcService{
			GrpcService: &gen.GRPCService{
				Host: "localhost:8080",
			},
		},
	}
	blocks, err := grpc.EnumerateResourceBlocks(r)
	if err != nil {
		t.Fatal(err)
	}
	getProjectsBlockId := blocks[1].Id
	p := &gen.Project{
		Graph: &gen.Graph{
			Nodes: []*gen.Node{
				{
					Id:      nodeID,
					BlockId: getProjectsBlockId,
					Config: &gen.Node_Grpc{
						Grpc: &gen.GRPC{
							Service: "ProjectService",
							Method:  "GetProjects",
						},
					},
				},
			},
		},
		Resources: []*gen.Resource{lr, r},
	}

	w, err := FromProject(p)
	if err != nil {
		t.Fatal(err)
	}

	ctx := execute.MemoryContext{context.Background()}
	executor := execute.NewMemoryExecutor(&ctx, nil)
	logger := &MemoryLogger{}
	_, err = w.Run(logger, executor, nodeID, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuildingGraph(t *testing.T) {
	p := &gen.Project{
		Graph: &gen.Graph{
			Edges: []*gen.Edge{
				{
					From: "input-node",
					To:   "crawl-node",
				},
				{
					From: "crawl-node",
					To:   "normalize-html-node",
				},
				{
					From: "normalize-html-node",
					To:   "create-embeddings-node",
				},
			},
			Nodes: []*gen.Node{
				{
					Id:   "input-node",
					Name: "Website",
					Config: &gen.Node_Input{
						Input: &gen.Input{
							Fields: []*gen.FieldDefinition{
								{
									Name: "url",
								},
							},
						},
					},
				},
				{
					Id:   "crawl-node",
					Name: "Crawl Website",
					Config: &gen.Node_Function{
						Function: &gen.Function{
							Runtime: "go",
						},
					},
					BlockId: "crawl-block",
				},
				{
					Id:   "normalize-html-node",
					Name: "Crawl Website",
					Config: &gen.Node_Function{
						Function: &gen.Function{
							Runtime: "go",
						},
					},
					BlockId: "normalize-html-block",
				},
				{
					Id:   "create-embeddings-node",
					Name: "Create Embeddings for HTML",
					Config: &gen.Node_Function{
						Function: &gen.Function{
							Runtime: "go",
						},
					},
					BlockId: "create-embeddings-block",
				},
			},
		},
	}

	w, err := FromProject(p, node.ResourceMap{})
	if err != nil {
		t.Fatal(err)
	}

	ctx := execute.MemoryContext{context.Background()}
	executor := execute.NewMemoryExecutor(&ctx)
	logger := &MemoryLogger{}

	entrypointNode := "input-node"
	input := `{"url": "https://example.com"}`
	_, err = w.Run(logger, executor, entrypointNode, input)
	if err != nil {
		t.Fatal(err)
	}
}
