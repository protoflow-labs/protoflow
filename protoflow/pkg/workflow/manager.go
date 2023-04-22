package workflow

import (
	"context"
	"fmt"
	"github.com/breadchris/protoflow/gen/workflow"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.temporal.io/api/enums/v1"
	"time"

	"go.temporal.io/sdk/client"
)

type TemporalManager struct {
	temporalClient client.Client
	workflowStore  Store
}

var ManagerProviderSet = wire.NewSet(
	NewManager,
	wire.Bind(new(workflow.Manager), new(*TemporalManager)),
)

var _ workflow.Manager = (*TemporalManager)(nil)

func NewManager(
	temporalClient client.Client,
	store Store,
) *TemporalManager {
	return &TemporalManager{
		temporalClient: temporalClient,
		workflowStore:  store,
	}
}

func (m *TemporalManager) CreateWorkflow(ctx context.Context, protoflow *workflow.Workflow) (*workflow.ID, error) {
	id, err := m.workflowStore.SaveWorkflow(protoflow)
	if err != nil {
		return nil, err
	}
	return &workflow.ID{
		Id: id,
	}, nil
}

func (m *TemporalManager) StartWorkflow(ctx context.Context, entry *workflow.WorkflowEntrypoint) (*workflow.Run, error) {
	protoflow, err := m.workflowStore.GetWorkflow(entry.WorkflowId)
	if err != nil {
		return nil, err
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        protoflow.Id,
		TaskQueue: taskQueue,
		// CronSchedule: workflow.CronSchedule,
	}

	w, err := NewWorkflowFromProtoflow(protoflow)
	if err != nil {
		return nil, err
	}

	we, err := m.temporalClient.ExecuteWorkflow(ctx, workflowOptions, w.Run, entry.NodeId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start workflow %s", protoflow.Id)
	}
	log.Debug().
		Str("id", we.GetID()).
		Str("run id", we.GetRunID()).
		Msg("started workflow")

	m.workflowStore.SaveWorkflowRun(we.GetID(), we.GetRunID())

	return &workflow.Run{
		Id: we.GetRunID(),
	}, nil
}

// ProcessNewWorkflows starts any workflow that needs to be started now and saves them.
func getWorkflowResult(ctx context.Context, temporalClient client.Client, workflowID, runID string) {
	var workflowCompeted bool

	for !workflowCompeted {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(125 * time.Millisecond)
			resp, err := temporalClient.DescribeWorkflowExecution(ctx, workflowID, "")
			if err != nil {
				//workflowResult.Error = err
				return
			}
			executionInfo := resp.GetWorkflowExecutionInfo()
			if enums.WORKFLOW_EXECUTION_STATUS_COMPLETED == executionInfo.GetStatus() {
				workflowCompeted = true
				break
			}
		}
	}

	//workflowRun := temporalClient.GetWorkflow(ctx, workflowID, runID)
	//workflowResult.Error = workflowRun.Get(ctx, &workflowResult.Data)
	return
}

func (m *TemporalManager) CancelWorkflows(deploymentID string, workflowRuns []WorkflowRunModel) error {
	for _, workflowRun := range workflowRuns {
		err := m.temporalClient.CancelWorkflow(context.Background(), workflowRun.WorkflowID, workflowRun.RunID)
		if err != nil {
			continue
		}
	}
	/*
		request := &workflowservice.ListOpenWorkflowExecutionsRequest{
			Namespace: "refinery",
			Filters: &workflowservice.ListOpenWorkflowExecutionsRequest_ExecutionFilter{
				ExecutionFilter: filter.ExecutionFilter{
					WorkflowId: "",
				},
			},
		}
	*/

	// TODO this should only delete the workflow runs specified
	err := m.workflowStore.DeleteDeploymentWorkflows(deploymentID)
	if err != nil {
		return fmt.Errorf("unable to delete workflows for deployment: %s", err)
	}
	return nil
}

func (m *TemporalManager) CancelWorkflowsForDeployment(deploymentID string) error {
	workflowRuns, err := m.workflowStore.GetWorkflowRunsForDeployment(deploymentID)
	if err != nil {
		return fmt.Errorf("unable to get workflow runs for deployment: %s", err)
	}
	return m.CancelWorkflows(deploymentID, workflowRuns)
}
