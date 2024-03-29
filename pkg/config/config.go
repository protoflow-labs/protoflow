package config

import (
	"github.com/protoflow-labs/protoflow/pkg/generate"
	"github.com/protoflow-labs/protoflow/pkg/llm"
	"github.com/protoflow-labs/protoflow/pkg/server"
	"os"
	"path"

	"github.com/protoflow-labs/protoflow/pkg/temporal"
	"github.com/protoflow-labs/protoflow/pkg/workflow"

	"github.com/protoflow-labs/protoflow/pkg/bucket"
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
	Server   server.Config   `yaml:"server"`
	Generate generate.Config `yaml:"generate"`
	LLM      llm.Config      `yaml:"llm"`
}

func NewDefaultConfig() BaseConfig {
	return BaseConfig{
		Log:      log.NewDefaultConfig(),
		Workflow: workflow.NewDefaultConfig(),
		DB:       db.NewDefaultConfig(),
		Temporal: temporal.NewDefaultConfig(),
		Server:   server.NewDefaultConfig(),
		Generate: generate.NewDefaultConfig(),
		LLM:      llm.NewDefaultConfig(),
	}
}

func NewProvider(cache bucket.Bucket) (config.Provider, error) {
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
