package resource

import (
	"github.com/protoflow-labs/protoflow/gen"
)

type SecretStoreResource struct {
	*BaseResource
	*gen.SecretStore
}

func (r *SecretStoreResource) Init() (func(), error) {
	return nil, nil
}
