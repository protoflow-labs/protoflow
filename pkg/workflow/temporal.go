package workflow

import (
	"context"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

//var _ Manager = (*TemporalManager)(nil)

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

func (m *TemporalManager) ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string, input interface{}) (string, error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        w.ID,
		TaskQueue: TaskQueue,
		// CronSchedule: workflow.CronSchedule,
	}

	we, err := m.client.ExecuteWorkflow(ctx, workflowOptions, TemporalRun, w, nodeID, input)
	if err != nil {
		return "", errors.Wrapf(err, "failed to start workflow %s", w.ID)
	}

	log.Debug().
		Str("id", we.GetID()).
		Str("run id", we.GetRunID()).
		Msg("started workflow")
	return we.GetRunID(), nil
}

func (m *TemporalManager) ExecuteWorkflowSync(ctx context.Context, w *Workflow, nodeID string, input any) (rxgo.Observable, error) {
	//TODO implement me
	panic("implement me")
}

// TemporalRun is the entrypoint for a Temporal workflow that will run on a worker
func TemporalRun(ctx workflow.Context, w *Workflow, nodeID string, input string) (rxgo.Observable, error) {
	//if w.NodeLookup == nil || w.Graph == nil {
	//	return nil, fmt.Errorf("workflow is not initialized")
	//}
	//
	//ao := workflow.ActivityOptions{
	//	ScheduleToStartTimeout: time.Minute,
	//	StartToCloseTimeout:    time.Minute,
	//	HeartbeatTimeout:       time.Second * 20,
	//}
	//ctx = workflow.WithActivityOptions(ctx, ao)
	//logger := workflow.GetLogger(ctx)
	//
	//executor := execute.NewTemporalExecutor(ctx)
	//
	//// Adding context to a workflow
	//// ctx = workflow.WithValue(ctx, AccountIDContextKey, dslWorkflow.AccountID)
	//
	//logger.Type("Starting workflow", "workflowID", workflow.GetInfo(ctx).WorkflowExecution.ID, "nodeID", nodeID)
	//
	//return w.WireNodes(context.Background(), logger, executor, nodeID, input)
	return nil, nil
}

func (m *TemporalManager) CleanupResources() error {
	//TODO implement me
	panic("implement me")
}
