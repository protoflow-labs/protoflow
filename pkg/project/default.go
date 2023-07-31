package project

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	pcode "github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/node"
	"github.com/protoflow-labs/protoflow/pkg/node/code"
)

func getDefaultProject(name string, bucketDir string) *gen.Project {
	p := code.NewProto("protoflow", code.NewServerProto(pcode.Runtime_GO))

	n := code.NewProto("NewNode", code.NewFunctionProto())

	return &gen.Project{
		Id:   uuid.NewString(),
		Name: name,
		Graph: &gen.Graph{
			Nodes: []*gen.Node{p, n},
			Edges: []*gen.Edge{node.NewEdge(p.Id, n.Id)},
		},
	}
}
