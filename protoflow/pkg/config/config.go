package config

import (
	"github.com/breadchris/protoflow/pkg/log"
	"github.com/breadchris/protoflow/pkg/workflow"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"go.uber.org/config"
	"os"
	"path"
)

const (
	localConfigFile = ".protoflow.yaml"
	homeConfigFile  = "config.yaml"
)

type BaseConfig struct {
	Log      log.Config      `yaml:"log"`
	Workflow workflow.Config `yaml:"workflow"`
}

func NewDefaultConfig() BaseConfig {
	return BaseConfig{
		Log:      log.NewDefaultConfig(),
		Workflow: workflow.NewDefaultConfig(),
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
