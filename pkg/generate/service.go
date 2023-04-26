package generate

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/project"
)

type Service struct {
	store project.Store
}

var _ genconnect.GenerateServiceHandler = (*Service)(nil)

var ProviderSet = wire.NewSet(
	NewService,
	wire.Bind(new(genconnect.GenerateServiceHandler), new(*Service)),
)

func NewService(store project.Store) (*Service, error) {
	return &Service{
		store: store,
	}, nil
}

func (s *Service) Generate(ctx context.Context, req *connect.Request[gen.GenerateRequest]) (*connect.Response[gen.GenerateResponse], error) {
	project, err := s.store.GetProject(req.Msg.ProjectId)
	if err != nil {
		return nil, err
	}

	generator, err := NewFromProject(project)
	if err != nil {
		return nil, err
	}

	if err := generator.Generate(); err != nil {
		return nil, err
	}

	return &connect.Response[gen.GenerateResponse]{
		Msg: &gen.GenerateResponse{
			ProjectId: project.Id,
		},
	}, nil
}
