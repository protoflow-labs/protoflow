// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cli

import (
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/protoflow-labs/protoflow/pkg/api"
	"github.com/protoflow-labs/protoflow/pkg/config"
	"github.com/protoflow-labs/protoflow/pkg/db"
	"github.com/protoflow-labs/protoflow/pkg/k8s"
	"github.com/protoflow-labs/protoflow/pkg/log"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/urfave/cli/v2"
)

// Injectors from wire.go:

func Wire(cacheConfig cache.Config) (*cli.App, error) {
	localCache, err := cache.NewLocalCache(cacheConfig)
	if err != nil {
		return nil, err
	}
	provider, err := config.NewProvider(localCache)
	if err != nil {
		return nil, err
	}
	logConfig, err := log.NewConfig(provider)
	if err != nil {
		return nil, err
	}
	clientset, err := k8s.NewClientset()
	if err != nil {
		return nil, err
	}
	dbConfig, err := db.NewConfig(provider)
	if err != nil {
		return nil, err
	}
	gormDB, err := db.NewGormDB(dbConfig, localCache)
	if err != nil {
		return nil, err
	}
	dbStore, err := project.NewDBStore(gormDB)
	if err != nil {
		return nil, err
	}
	workflowConfig, err := workflow.NewConfig(provider)
	if err != nil {
		return nil, err
	}
	manager, err := workflow.NewManager(workflowConfig, provider)
	if err != nil {
		return nil, err
	}
	service, err := project.NewService(clientset, dbStore, manager)
	if err != nil {
		return nil, err
	}
	httpServer := api.NewHTTPServer(service)
	app := New(logConfig, httpServer, service, provider)
	return app, nil
}