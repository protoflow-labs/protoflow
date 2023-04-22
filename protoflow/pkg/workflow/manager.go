package workflow

import (
	"context"
	"github.com/breadchris/protoflow/gen/workflow"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"time"

	"go.temporal.io/sdk/client"
)

const taskQueue = "protoflow"

type TemporalManager struct {
	workflow.UnimplementedManagerServer

	temporalClient client.Client
	workflowStore  Store
	config         Config
}

var ProviderSet = wire.NewSet(
	NewConfig,
	StoreProviderSet,
	NewClient,
	NewManager,
	wire.Bind(new(workflow.Manager), new(*TemporalManager)),
)

var _ workflow.Manager = (*TemporalManager)(nil)

func NewManager(
	temporalClient client.Client,
	store Store,
	config Config,
) *TemporalManager {
	return &TemporalManager{
		temporalClient: temporalClient,
		workflowStore:  store,
		config:         config,
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
	return nil
}

// StopAllOpenWorkflows gets all currently open workflows and cancels them.
func (s *TemporalManager) StopAllOpenWorkflows(c *fiber.Ctx) error {
	request := &workflowservice.ListOpenWorkflowExecutionsRequest{
		Namespace: s.config.Temporal.Namespace,
	}

	// TODO call this in loop to get all executions
	response, err := s.temporalClient.ListOpenWorkflow(context.Background(), request)
	if err != nil {
		return err
	}
	for _, execution := range response.GetExecutions() {
		workflowExec := execution.GetExecution()
		err := s.temporalClient.CancelWorkflow(context.Background(), workflowExec.WorkflowId, workflowExec.RunId)
		if err != nil {
			return err
		}
	}
	return nil
}
