// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cli

import (
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/api"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/config"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/db"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/log"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/workflow"
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
	workflowConfig, err := workflow.NewConfig(provider)
	if err != nil {
		return nil, err
	}
	client, err := workflow.NewClient(workflowConfig)
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
	dbStore, err := workflow.NewDBStore(gormDB)
	if err != nil {
		return nil, err
	}
	temporalManager := workflow.NewManager(client, dbStore, workflowConfig)
	httpServer := api.NewHTTPServer(temporalManager)
	grpcServer := api.NewGRPCServer(temporalManager)
	worker := workflow.NewWorker(client, workflowConfig)
	app := New(logConfig, httpServer, grpcServer, worker)
	return app, nil
}
