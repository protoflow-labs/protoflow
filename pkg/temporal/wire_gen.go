// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package temporal

import (
	"go.temporal.io/sdk/client"
	"go.uber.org/config"
)

// Injectors from wire.go:

func Wire(provider config.Provider) (client.Client, error) {
	temporalConfig, err := NewConfig(provider)
	if err != nil {
		return nil, err
	}
	clientClient, err := NewClient(temporalConfig)
	if err != nil {
		return nil, err
	}
	return clientClient, nil
}