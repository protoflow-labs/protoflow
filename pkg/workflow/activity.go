package workflow

import (
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/workflow"
)

type Activity struct{}

func (a *Activity) ExecuteCode(ctx workflow.Context, code *Code, input Input) (Result, error) {
	log.Debug().Msgf("executing code: %s", code.Code)
	return Result{}, nil
}

func (a *Activity) ExecuteInput(ctx workflow.Context, data *Data, input Input) (Result, error) {
	log.Debug().Msgf("executing input: %v", data.Params)
	return Result{}, nil
}
