package workflow

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/pkg/temporal"
	"go.uber.org/config"
)

const TaskQueue = "protoflow"

type Manager interface {
	ExecuteWorkflow(ctx context.Context, w *Workflow, nodeID string, input string) (string, error)
	ExecuteWorkflowSync(ctx context.Context, w *Workflow, nodeID string, input string) (*Result, error)
	CleanupResources() error
}

var ProviderSet = wire.NewSet(
	NewConfig,
	NewManager,
)

func NewManager(config Config, provider config.Provider) (Manager, error) {
	switch config.ManagerType {
	case MemoryManagerType:
		return NewMemoryManager(), nil
	case TemporalManagerType:
		// TODO breadchris we do this because we don't want a temporal client to try to connect on startup
		// Is there a way to run this more intelligently? maybe with sync.Once?
		client, err := temporal.Wire(provider)
		if err != nil {
			return nil, err
		}
		return NewTemporalManager(client), nil
	default:
		return nil, fmt.Errorf("unknown manager type %s", config.ManagerType)
	}
}

/*
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

	workflowRun := temporalClient.GetProject(ctx, workflowID, runID)
	workflowResult.Error = workflowRun.Get(ctx, &workflowResult.RESTNode)
	return
}

func CancelWorkflows(deploymentID string, workflowRuns []WorkflowRunModel) error {
	for _, workflowRun := range workflowRuns {
		err := m.temporalClient.CancelWorkflow(context.Background(), workflowRun.WorkflowID, workflowRun.RunID)
		if err != nil {
			continue
		}
	}
		request := &workflowservice.ListOpenWorkflowExecutionsRequest{
			Namespace: "refinery",
			Filters: &workflowservice.ListOpenWorkflowExecutionsRequest_ExecutionFilter{
				ExecutionFilter: filter.ExecutionFilter{
					WorkflowId: "",
				},
			},
		}
	return nil
}

// StopAllOpenWorkflows gets all currently open workflows and cancels them.
func StopAllOpenWorkflows(c *fiber.Ctx) error {
	request := &workflowservice.ListOpenWorkflowExecutionsRequest{
		Namespace: m.config.Temporal.Namespace,
	}

	// TODO call this in loop to get all executions
	response, err := m.temporalClient.ListOpenWorkflow(context.Background(), request)
	if err != nil {
		return err
	}
	for _, execution := range response.GetExecutions() {
		workflowExec := execution.GetExecution()
		err := m.temporalClient.CancelWorkflow(context.Background(), workflowExec.WorkflowId, workflowExec.RunId)
		if err != nil {
			return err
		}
	}
	return nil
}
*/
