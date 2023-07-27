package project

import (
	"context"
	"encoding/json"
	"github.com/bufbuild/connect-go"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/protoflow-labs/protoflow/pkg/workflow/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow/node"
	"github.com/protoflow-labs/protoflow/pkg/workflow/resource"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
)

func (s *Service) startWorkflow(
	ctx context.Context,
	w *workflow.Workflow,
	nodeID string,
	workflowInput any,
	// TODO breadchris this should not be needed
	httpStream *resource.HTTPEventStream,
) (rxgo.Observable, error) {
	log.Debug().
		Str("workflow", w.ID).
		Str("node", nodeID).
		Msg("workflow starting")

	n, ok := w.NodeLookup[nodeID]
	if !ok {
		return nil, errors.Errorf("node %s not found in workflow", nodeID)
	}

	var (
		inputChan   chan rxgo.Item
		inputObs    rxgo.Observable
		httpRequest bool
	)
	switch n.(type) {
	case *node.RouteNode:
		inputObs = httpStream.RequestObs
		httpRequest = true
	default:
		inputChan = make(chan rxgo.Item)
		inputObs = rxgo.FromChannel(inputChan, rxgo.WithPublishStrategy())
	}

	obs, err := s.manager.ExecuteWorkflow(ctx, w, nodeID, inputObs)
	if err != nil {
		return nil, err
	}

	// TODO breadchris support streaming input
	if !httpRequest {
		inputChan <- rx.NewItem(workflowInput)
		close(inputChan)
	}
	return obs, nil
}

func (s *Service) RunWorkflow(ctx context.Context, c *connect.Request[gen.RunWorkflowRequest], c2 *connect.ServerStream[gen.NodeExecution]) error {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := workflow.Default().
		WithProtoProject(graph.ConvertProto(project)).
		Build()
	if err != nil {
		return err
	}

	// TODO breadchris temporary for when the input is not set
	if c.Msg.Input == "" {
		c.Msg.Input = "{}"
	}

	// TODO breadchris this is a _little_ sketchy, we would like to be able to use the correct type, which might just be some data!
	var workflowInput map[string]any
	err = json.Unmarshal([]byte(c.Msg.Input), &workflowInput)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal workflow input")
	}

	var (
		entrypoints []string
		observables []rxgo.Observable
	)

	httpStream := resource.NewHTTPEventStream()

	if c.Msg.StartServer {
		for _, n := range w.NodeLookup {
			switch n.(type) {
			case *node.RouteNode:
				entrypoints = append(entrypoints, n.ID())
			}
		}
	} else {
		entrypoints = append(entrypoints, c.Msg.NodeId)
	}

	for _, entrypoint := range entrypoints {
		obs, err := s.startWorkflow(ctx, w, entrypoint, workflowInput, httpStream)
		if err != nil {
			return errors.Wrapf(err, "failed to start workflow")
		}
		observables = append(observables, obs)
	}

	obs := rxgo.Merge(observables)

	var (
		obsErr error
	)
	<-obs.ForEach(func(item any) {
		switch t := item.(type) {
		case *gen.HttpResponse:
			httpStream.Responses <- t
			log.Debug().Msg("sent http response")
		}

		log.Debug().Interface("item", item).Msg("workflow item")
		out, err := json.Marshal(item)
		if err != nil {
			obsErr = errors.Wrapf(err, "failed to marshal result data")
			return
		}

		// TODO breadchris node executions should be passed to the observable with the node id, input, and output
		err = c2.Send(&gen.NodeExecution{
			Output: string(out),
		})
		if err != nil {
			obsErr = errors.Wrapf(err, "failed to send node execution")
			return
		}
	}, func(err error) {
		obsErr = err
	}, func() {
		log.Debug().
			Str("workflow", w.ID).
			Str("node", c.Msg.NodeId).
			Msg("workflow finished")
	})
	return obsErr
}

func (s *Service) StopWorkflow(ctx context.Context, c *connect.Request[gen.StopWorkflowRequest]) (*connect.Response[gen.StopWorkflowResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetWorkflowRuns(ctx context.Context, c *connect.Request[gen.GetWorkflowRunsRequest]) (*connect.Response[gen.GetWorkflowRunsResponse], error) {
	runs, err := s.store.GetWorkflowRunsForProject(c.Msg.ProjectId)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&gen.GetWorkflowRunsResponse{Runs: runs}), nil
}
