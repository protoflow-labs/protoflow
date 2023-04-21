package workflow

import (
	"context"
	"fmt"
	"github.com/breadchris/protoflow/gen/workflow"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lunabrain-ai/lunabrain/pkg/store/temporal/dsl"
	"github.com/rs/zerolog/log"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

const taskQueue = "protoflow"

type Service struct {
	temporalClient client.Client
	store          Store
	config         Config
}

func NewService(
	temporalClient client.Client,
	store Store,
	config Config,
) *Service {
	return &Service{
		temporalClient: temporalClient,
		store:          store,
		config:         config,
	}
}

func (s *Service) CreateWorkflow(ctx context.Context, w *workflow.Workflow) (*workflow.ID, error) {
	workflowID := uuid.New().String()
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQueue,
	}

	we, err := s.temporalClient.ExecuteWorkflow(ctx, workflowOptions, dsl.SimpleDSLWorkflow, w)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Str("id", we.GetID()).
		Str("run id", we.GetRunID()).
		Msg("started workflow")

	s.store.SaveWorkflowRun(we.GetID(), we.GetRunID())

	return &workflow.ID{
		Id: we.GetRunID(),
	}, nil
}

func (s *Service) StartWorkflow(c *fiber.Ctx) error {
	waitForCompletion := false

	ctx := context.Background()
	if waitForCompletion {
		ctx = c.Context()
	}

	deploymentID := c.Params("deploymentID")
	workflowID := c.Params("workflowID")

	log.Debug().Str("workflow id", workflowID).Msg("starting workflow")

	fmt.Printf("%s", c.Body())

	if waitForCompletion {
		workflowResults := s.manager.WaitForWorkflowRunsResults(ctx, workflowID, startedWorkflows)

		returnData, err := formatReturnData(workflowResults)
		if err != nil {
			return err
		}
		c.Write(returnData)
		return nil
	}
	return nil
}

// CancelWorkflowsForDeployment cancels all workflows for a given deployment.
func (s *Service) CancelWorkflowsForDeployment(c *fiber.Ctx) error {
	deploymentID := c.Params("deploymentID")

	if err := s.manager.CancelWorkflowsForDeployment(deploymentID); err != nil {
		return err
	}
	return nil
}

// StopAllOpenWorkflows gets all currently open workflows and cancels them.
func (s *Service) StopAllOpenWorkflows(c *fiber.Ctx) error {
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
