package project

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
)

func (s *Service) EnumerateProviders(ctx context.Context, c *connect.Request[gen.GetProvidersRequest]) (*connect.Response[gen.GetProvidersResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := workflow.Default().
		WithProtoProject(graph.ConvertProto(project)).
		Build()

	var providers []*gen.EnumeratedProvider
	for _, node := range project.Graph.Nodes {
		info := &gen.ProviderInfo{
			State: gen.ProviderState_READY,
			Error: "",
		}
		n, ok := w.NodeLookup[node.Id]
		if !ok {
			info.State = gen.ProviderState_ERROR
			info.Error = "node not found"
		}

		providedNodes, err := n.Provide()
		if len(providedNodes) == 0 && err == nil {
			continue
		}

		if err != nil {
			info.State = gen.ProviderState_ERROR
			info.Error = err.Error()
		}
		providers = append(providers, &gen.EnumeratedProvider{
			Provider: node,
			Nodes:    providedNodes,
			Info:     info,
		})
	}

	return connect.NewResponse(&gen.GetProvidersResponse{
		Providers: providers,
	}), nil
}
