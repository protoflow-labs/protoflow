package workflow

import "go.temporal.io/sdk/client"

func NewClient(config Config) (client.Client, error) {
	//return client.Dial(client.Options{
	//	HostPort:  config.Temporal.Host,
	//	Namespace: config.Temporal.Namespace,
	//})
	return nil, nil
}
