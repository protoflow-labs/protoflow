package workflow

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"go.temporal.io/api/enums/v1"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

type Manager interface {
	StartWorkflow(ctx context.Context, workflowID, nodeID string) (string, error)
}

type TemporalManager struct {
	temporalClient client.Client
	workflowStore  Store
}

var ProviderSet = wire.NewSet(
	NewConfig,
	NewClient,
	DBProviderSet,
	NewManager,
)

func NewManager(
	temporalClient client.Client,
	store Store,
) *TemporalManager {
	return &TemporalManager{
		temporalClient: temporalClient,
		workflowStore:  store,
	}
}

func (m *TemporalManager) StartWorkflow(
	ctx context.Context,
	workflowID,
	nodeID string,
) (string, error) {
	protoflow, err := m.workflowStore.GetWorkflow(workflowID)
	if err != nil {
		return "", err
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        protoflow.Id,
		TaskQueue: "dsl",
		// CronSchedule: workflow.CronSchedule,
	}

	w, err := NewWorkflowFromProtoflow(protoflow)
	if err != nil {
		return "", err
	}

	we, err := m.temporalClient.ExecuteWorkflow(ctx, workflowOptions, w.Run, nodeID)
	if err != nil {
		log.Println("unable to execute workflow", err)
		return "", err
	}
	log.Println("started workflow", "ID", we.GetID(), "RunID", we.GetRunID())

	m.workflowStore.SaveWorkflowRun(we.GetID(), we.GetRunID())

	return we.GetRunID(), nil
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
			log.Println("Error while canceling workflow", err)
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
