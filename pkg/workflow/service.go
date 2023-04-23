package workflow

import (
	"context"
<<<<<<<< HEAD:pkg/workflow/service.go
	"github.com/bufbuild/connect-go"
========
	"time"

	connect_go "github.com/bufbuild/connect-go"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"

>>>>>>>> b308591b69c96f5b4aa93f141f29a0b8a8f2c82d:pkg/workflow/manager.go
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow-editor/protoflow/gen"
	"github.com/protoflow-labs/protoflow-editor/protoflow/gen/genconnect"
	"github.com/rs/zerolog/log"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"time"
)

const taskQueue = "protoflow"

var ProviderSet = wire.NewSet(
	NewConfig,
	StoreProviderSet,
	NewClient,
	NewService,
	wire.Bind(new(genconnect.ManagerHandler), new(*Service)),
)

var _ genconnect.ManagerHandler = (*Service)(nil)

type Service struct {
	genconnect.UnimplementedManagerHandler

	temporalClient client.Client
	workflowStore  Store
	config         Config
}

func NewService(
	temporalClient client.Client,
	store Store,
	config Config,
) *Service {
	return &Service{
		temporalClient: temporalClient,
		workflowStore:  store,
		config:         config,
	}
}

func (m *Service) CreateWorkflow(
	ctx context.Context,
	req *connect.Request[gen.Workflow],
) (*connect.Response[gen.ID], error) {
	id, err := m.workflowStore.SaveWorkflow(req.Msg)
	if err != nil {
		return nil, err
	}
	return &connect.Response[gen.ID]{
		Msg: &gen.ID{
			Id: id,
		},
	}, nil
}

func (m *Service) StartWorkflow(
	ctx context.Context,
	req *connect.Request[gen.WorkflowEntrypoint],
) (*connect.Response[gen.Run], error) {
	protoflow, err := m.workflowStore.GetWorkflow(req.Msg.WorkflowId)
	if err != nil {
		return nil, err
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        req.Msg.WorkflowId,
		TaskQueue: taskQueue,
		// CronSchedule: workflow.CronSchedule,
	}

	w, err := NewWorkflowFromProtoflow(protoflow)
	if err != nil {
		return nil, err
	}

	we, err := m.temporalClient.ExecuteWorkflow(ctx, workflowOptions, Run, w, req.Msg.NodeId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start workflow %s", protoflow.Id)
	}
	log.Debug().
		Str("id", we.GetID()).
		Str("run id", we.GetRunID()).
		Msg("started workflow")

	return &connect.Response[gen.Run]{
		Msg: &gen.Run{
			Id: we.GetRunID(),
		},
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

func (m *Service) CancelWorkflows(deploymentID string, workflowRuns []WorkflowRunModel) error {
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
func (m *Service) StopAllOpenWorkflows(c *fiber.Ctx) error {
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
