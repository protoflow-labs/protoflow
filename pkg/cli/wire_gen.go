// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cli

import (
	"github.com/protoflow-labs/protoflow/pkg/api"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/config"
	"github.com/protoflow-labs/protoflow/pkg/db"
	"github.com/protoflow-labs/protoflow/pkg/generate"
	"github.com/protoflow-labs/protoflow/pkg/log"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/store"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/urfave/cli/v2"
)

// Injectors from wire.go:

func Wire(cacheConfig bucket.Config) (*cli.App, error) {
	localCache, err := bucket.NewUserCache(cacheConfig)
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
	apiConfig, err := api.NewConfig(provider)
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
	projectStore, err := store.NewDBStore(gormDB)
	if err != nil {
		return nil, err
	}
	workflowConfig, err := workflow.NewConfig(provider)
	if err != nil {
		return nil, err
	}
	manager, err := workflow.NewManager(workflowConfig, provider, projectStore)
	if err != nil {
		return nil, err
	}
	service, err := project.NewService(projectStore, manager, localCache)
	if err != nil {
		return nil, err
	}
	generateConfig, err := generate.NewConfig(provider)
	if err != nil {
		return nil, err
	}
	generateService, err := generate.NewService(generateConfig, projectStore)
	if err != nil {
		return nil, err
	}
	httpServer, err := api.NewHTTPServer(apiConfig, service, generateService)
	if err != nil {
		return nil, err
	}
	app := New(logConfig, httpServer, service, provider)
	return app, nil
}
