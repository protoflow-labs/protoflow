package project

import (
	"github.com/google/uuid"
	"github.com/protoflow-labs/protoflow/gen"
	pcode "github.com/protoflow-labs/protoflow/gen/code"
	"github.com/protoflow-labs/protoflow/pkg/graph/edge"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/code"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/data"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/storage"
	"net/url"
)

type ProjectConfig struct {
	ID string
}

func getDefaultProject(name string, bucketDir string) *gen.Project {
	pid := uuid.NewString()
	u := url.URL{Path: bucketDir, Scheme: "file"}
	p := code.NewProto("protoflow", code.NewServerProto(pcode.Runtime_GO))
	n := code.NewProto("NewNode", code.NewFunctionProto())
	c := data.NewProto("config", data.NewConfigProto(ProjectConfig{ID: pid}))
	f := storage.NewProto("files", storage.NewFolderProto(u.String()))
	return &gen.Project{
		Id:   pid,
		Name: name,
		Graph: &gen.Graph{
			Nodes: []*gen.Node{p, n, c, f},
			Edges: []*gen.Edge{edge.NewProvidesProto(p.Id, n.Id)},
		},
	}
}
