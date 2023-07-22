package workflow

import (
	"context"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	client client.Client
}

func NewWorker(client client.Client) *Worker {
	return &Worker{
		client: client,
	}
}

func (s *Worker) Run() error {
	w := worker.New(s.client, TaskQueue, worker.Options{
		BackgroundActivityContext: context.Background(),
	})

	w.RegisterWorkflow(TemporalRun)
	// TODO breadchris register node activities
	// w.RegisterActivity(&execute.Activity{})

	return w.Run(worker.InterruptCh())
}
