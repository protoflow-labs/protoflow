package config

import (
	"github.com/protoflow-labs/protoflow/pkg/temporal"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"os"
	"path"

	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/protoflow-labs/protoflow/pkg/db"
	"github.com/protoflow-labs/protoflow/pkg/log"
	"go.uber.org/config"
)

const (
	localConfigFile = ".protoflow.yaml"
	homeConfigFile  = "config.yaml"
)

type BaseConfig struct {
	Log      log.Config      `yaml:"log"`
	Workflow workflow.Config `yaml:"workflow"`
	DB       db.Config       `yaml:"db"`
	Temporal temporal.Config `yaml:"temporal"`
}

func NewDefaultConfig() BaseConfig {
	return BaseConfig{
		Log:      log.NewDefaultConfig(),
		Workflow: workflow.NewDefaultConfig(),
		DB:       db.NewDefaultConfig(),
		Temporal: temporal.NewDefaultConfig(),
	}
}

func NewProvider(cache cache.Cache) (config.Provider, error) {
	opts := []config.YAMLOption{
		config.Permissive(),
		config.Expand(os.LookupEnv),
		config.Static(NewDefaultConfig()),
	}

	if f, err := os.Stat(localConfigFile); err == nil {
		opts = append(opts, config.File(path.Join(f.Name())))
	}

	homeFile, err := cache.GetFile(homeConfigFile)
	if err == nil {
		if _, err := os.Stat(homeFile); err == nil {
			opts = append(opts, config.File(homeFile))
		}
	}
	return config.NewYAML(opts...)
}
