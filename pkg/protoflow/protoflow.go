package protoflow

import (
	"github.com/protoflow-labs/protoflow/pkg/api"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/store"
)

type Protoflow struct {
	HTTPServer *api.HTTPServer
	Service    *project.Service
	Store      store.Project
}

func New(HTTPServer *api.HTTPServer, Service *project.Service, Store store.Project) *Protoflow {
	return &Protoflow{
		HTTPServer: HTTPServer,
		Service:    Service,
		Store:      Store,
	}
}
