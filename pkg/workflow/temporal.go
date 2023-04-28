package workflow

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"time"
)

var _ Manager = (*TemporalManager)(nil)

var TemporalManagerProviderSet = wire.NewSet(
	NewTemporalManager,
	wire.Bind(new(Manager), new(*TemporalManager)),
)

type TemporalManager struct {
	client client.Client
}

func NewTemporalManager(client client.Client) *TemporalManager {
	return &TemporalManager{
		client: client,
	}
}

func (m *TemporalManager) ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string) (string, error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        w.ID,
		TaskQueue: TaskQueue,
		// CronSchedule: workflow.CronSchedule,
	}

	we, err := m.client.ExecuteWorkflow(ctx, workflowOptions, TemporalRun, w, nodeID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to start workflow %s", w.ID)
	}

	log.Debug().
		Str("id", we.GetID()).
		Str("run id", we.GetRunID()).
		Msg("started workflow")
	return we.GetRunID(), nil
}

func (m *TemporalManager) ExecuteWorkflowSync(ctx context.Context, w *Workflow, nodeID string) (*Result, error) {
	//TODO implement me
	panic("implement me")
}

// TemporalRun is the entrypoint for a Temporal workflow that will run on a worker
func TemporalRun(ctx workflow.Context, w *Workflow, nodeID string) (*Result, error) {
	if w.NodeLookup == nil || w.Graph == nil {
		return nil, fmt.Errorf("workflow is not initialized")
	}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	executor := NewTemporalExecutor(ctx)

	// Adding context to a workflow
	// ctx = workflow.WithValue(ctx, AccountIDContextKey, dslWorkflow.AccountID)

	logger.Info("Starting workflow", "workflowID", workflow.GetInfo(ctx).WorkflowExecution.ID, "nodeID", nodeID)

	return w.Run(logger, executor, nodeID)
}

func (m *TemporalManager) CleanupResources() error {
	//TODO implement me
	panic("implement me")
}
