package project

import (
	"context"
	"encoding/json"
	"github.com/bufbuild/connect-go"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	http2 "github.com/protoflow-labs/protoflow/pkg/graph/node/http"
	"github.com/protoflow-labs/protoflow/pkg/util/rx"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/reactivex/rxgo/v2"
	"github.com/rs/zerolog/log"
)

func (s *Service) wireWorkflow(
	ctx context.Context,
	w *workflow.Workflow,
	nodeID string,
	workflowInput any,
	// TODO breadchris this should not be needed
	httpStream *http2.HTTPEventStream,
	req *gen.RunWorkflowRequest,
) (*graph.IO, error) {
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
	case *http2.RouteNode:
		inputObs = httpStream.RequestObs
		httpRequest = true
	default:
		inputChan = make(chan rxgo.Item)
		inputObs = rxgo.FromChannel(inputChan, rxgo.WithPublishStrategy())
	}

	m, err := s.manager.WithReq(req).Build()
	if err != nil {
		return nil, err
	}
	io, err := w.WireNodes(ctx, nodeID, inputObs, m)
	if err != nil {
		return nil, err
	}

	// TODO breadchris support streaming input
	if !httpRequest {
		// start workflow by sending the first input
		inputChan <- rx.NewItem(workflowInput)
		close(inputChan)
	}
	return io, nil
}

func (s *Service) RunWorkflow(ctx context.Context, c *connect.Request[gen.RunWorkflowRequest], c2 *connect.ServerStream[gen.NodeExecution]) error {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := FromProto(project)
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

	httpStream := http2.NewHTTPEventStream()

	if c.Msg.StartServer {
		for _, n := range w.NodeLookup {
			switch n.(type) {
			case *http2.RouteNode:
				entrypoints = append(entrypoints, n.ID())
			}
		}
	} else {
		entrypoints = append(entrypoints, c.Msg.NodeId)
	}

	for _, entrypoint := range entrypoints {
		io, err := s.wireWorkflow(ctx, w, entrypoint, workflowInput, httpStream, c.Msg)
		if err != nil {
			return errors.Wrapf(err, "failed to start workflow")
		}
		observables = append(observables, io.Observable)

		// TODO breadchris is there a better way to do this?
		// when the request context is cancelled, we want to cleanup the workflow
		//go func(io *graph.IO) {
		//	<-ctx.Done()
		//	io.Cleanup()
		//}(io)
	}
	log.Debug().Msg("done wiring workflows")

	obs := rxgo.Merge(observables)

	var (
		obsErr error
	)
	<-obs.ForEach(func(item any) {
		log.Debug().Interface("item", item).Msg("workflow item")
		out, err := json.Marshal(item)
		if err != nil {
			obsErr = errors.Wrapf(err, "failed to marshal result data")
			return
		}

		// TODO breadchris node executions should be passed to the observable with the node wID, input, and output
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
	if obsErr != nil {
		log.Error().Err(obsErr).Msg("workflow error")
	}
	return obsErr
}

func (s *Service) GetRunningWorkflows(ctx context.Context, c *connect.Request[gen.GetRunningWorkflowsRequest]) (*connect.Response[gen.GetRunningWorkflowResponse], error) {
	return connect.NewResponse(&gen.GetRunningWorkflowResponse{Traces: s.workflowManager.Traces()}), nil
}

func (s *Service) StopWorkflow(ctx context.Context, c *connect.Request[gen.StopWorkflowRequest]) (*connect.Response[gen.StopWorkflowResponse], error) {
	err := s.workflowManager.Stop(c.Msg.WorkflowId)
	return connect.NewResponse(&gen.StopWorkflowResponse{}), err
}

func (s *Service) GetWorkflowRuns(ctx context.Context, c *connect.Request[gen.GetWorkflowRunsRequest]) (*connect.Response[gen.GetWorkflowRunsResponse], error) {
	// TODO breadchris this should also consider running workflows
	runs, err := s.store.GetWorkflowRunsForProject(c.Msg.ProjectId)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&gen.GetWorkflowRunsResponse{Runs: runs}), nil
}
