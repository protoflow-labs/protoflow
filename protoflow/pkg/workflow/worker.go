package workflow

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Activity struct {
}

func (a *Activity) ExecuteCode(ctx context.Context, code Code) (*Result, error) {
	log.Debug().Msgf("executing code: %s", code.code)
	return nil, nil
}

func (a *Activity) ExecuteInput(ctx context.Context, input Input) (*Result, error) {
	log.Debug().Msgf("executing input: %v", input.params)
	return nil, nil
}

type Worker struct {
	config Config
	client client.Client
}

func NewWorker(client client.Client, config Config) *Worker {
	return &Worker{
		config: config,
		client: client,
	}
}

func (s *Worker) Run() error {
	w := worker.New(s.client, taskQueue, worker.Options{
		BackgroundActivityContext: context.Background(),
	})

	w.RegisterActivity(&Activity{})

	return w.Run(worker.InterruptCh())
}
