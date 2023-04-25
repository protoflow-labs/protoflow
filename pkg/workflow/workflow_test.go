package workflow

import (
	"context"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"testing"
)

func TestRun(t *testing.T) {
	r := &gen.Resource{
		Id: "1",
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
	r.Blocks = blocks
	p := &gen.Project{
		Graph: &gen.Graph{
			Nodes: []*gen.Node{
				{
					Id:      "1",
					BlockId: getProjectsBlockId,
				},
			},
		},
		Resources: []*gen.Resource{r},
	}

	w, err := FromProject(p)
	if err != nil {
		t.Fatal(err)
	}

	ctx := MemoryContext{context.Background()}
	executor := NewMemoryExecutor(&ctx)
	logger := &MemoryLogger{}
	res, err := w.Run(logger, executor, "1")
	if err != nil {
		t.Fatal(err)
	}

	println(res.Data)
}
