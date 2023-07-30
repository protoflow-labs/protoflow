package project

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
)

func (s *Service) EnumerateProviders(ctx context.Context, c *connect.Request[gen.GetProvidersRequest]) (*connect.Response[gen.GetProvidersResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	providers, err := enumerateProvidersFromNodes(project.Graph.GetNodes())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GetProvidersResponse{
		Providers: providers,
	}), nil
}
