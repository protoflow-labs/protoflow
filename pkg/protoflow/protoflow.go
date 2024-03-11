package protoflow

import (
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/server"
	"github.com/protoflow-labs/protoflow/pkg/store"
)

type Protoflow struct {
	HTTPServer *server.HTTPServer
	Service    *project.Service
	Store      store.Project
}

func New(HTTPServer *server.HTTPServer, Service *project.Service, Store store.Project) *Protoflow {
	return &Protoflow{
		HTTPServer: HTTPServer,
		Service:    Service,
		Store:      Store,
	}
}
