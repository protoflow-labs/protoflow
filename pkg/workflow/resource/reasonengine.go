package resource

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	openaiclient "github.com/protoflow-labs/protoflow/pkg/openai"
	"go.uber.org/config"
)

type ReasoningEngineResource struct {
	*BaseResource
	*gen.ReasoningEngine
	QAClient openaiclient.QAClient
}

func (r *ReasoningEngineResource) Init() (func(), error) {
	// TODO breadchris replace with some type of dependency injection capability
	var (
		configProvider config.Provider
		err            error
	)
	staticConfig := map[string]interface{}{
		"openai": openaiclient.NewDefaultConfig(),
	}
	for _, n := range r.dependencyLookup {
		switch t := n.(type) {
		case *ConfigProviderResource:
			// TODO breadchris how do we handle resources that need to be initialized before others?
			configProvider, err = t.NewConfigProvider(config.Static(staticConfig))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to build config provider")
			}
		}
	}
	if configProvider == nil {
		return nil, errors.New("config provider not found")
	}
	c, err := openaiclient.Wire(configProvider)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize openai client")
	}
	r.QAClient = c
	return nil, nil
}
