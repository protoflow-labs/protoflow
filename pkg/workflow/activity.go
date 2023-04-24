package workflow

import (
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/workflow"
)

type Activity struct{}

func (a *Activity) ExecuteGRPCBlock(ctx workflow.Context, block *GRPCBlock, input Input) (Result, error) {
	log.Debug().Msgf("executing block: %s", block.Host)
	return Result{}, nil
}

func (a *Activity) ExecuteRestBlock(ctx workflow.Context, block *RESTBlock, input Input) (Result, error) {
	log.Debug().Msgf("executing input: %v", block.Url)
	return Result{}, nil
}
