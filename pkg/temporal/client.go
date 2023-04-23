package temporal

import (
	"github.com/google/wire"
	"go.temporal.io/sdk/client"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewClient,
)

func NewClient(config Config) (client.Client, error) {
	return client.Dial(client.Options{
		HostPort:  config.Host,
		Namespace: config.Namespace,
	})
}
