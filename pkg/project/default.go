package project

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	pcode "github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/graph/edge"
	code2 "github.com/protoflow-labs/protoflow/pkg/graph/node/code"
)

func getDefaultProject(name string, bucketDir string) *gen.Project {
	p := code2.NewProto("protoflow", code2.NewServerProto(pcode.Runtime_GO))
	n := code2.NewProto("NewNode", code2.NewFunctionProto())
	return &gen.Project{
		Id:   uuid.NewString(),
		Name: name,
		Graph: &gen.Graph{
			Nodes: []*gen.Node{p, n},
			Edges: []*gen.Edge{edge.NewProvidesProto(p.Id, n.Id)},
		},
	}
}
