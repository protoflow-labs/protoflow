//go:build wireinject
// +build wireinject

package temporal

import (
	"github.com/google/wire"
	"go.temporal.io/sdk/client"
	"go.uber.org/config"
)

func Wire(provider config.Provider) (client.Client, error) {
	panic(wire.Build(
		ProviderSet,
	))
}
