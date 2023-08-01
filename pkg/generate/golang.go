package generate

import (
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/graph/node/code"
)

type GoManager struct {
	codeRoot bucket.Bucket
}

func (g GoManager) GenerateGRPCService(r *code.Server) error {
	//TODO implement me
	panic("implement me")
}

var _ LanguageManager = &GoManager{}

func NewGoManager(c bucket.Bucket) (*GoManager, error) {
	codeRoot, err := c.WithDir("golang")
	if err != nil {
		return nil, errors.Wrapf(err, "error creating bucket for %s", "nodejs")
	}
	return &GoManager{
		codeRoot: codeRoot,
	}, nil
}
