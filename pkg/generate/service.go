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
	store     project.Store
	generator Generator
}

var _ genconnect.GenerateServiceHandler = (*Service)(nil)

var ProviderSet = wire.NewSet(
	NewService,
	NewGenerate,
	wire.Bind(new(genconnect.GenerateServiceHandler), new(*Service)),
	wire.Bind(new(Generator), new(*Generate)),
)

func NewService(store project.Store, generator Generator) (*Service, error) {
	return &Service{
		store:     store,
		generator: generator,
	}, nil
}

func (s *Service) Generate(ctx context.Context, req *connect.Request[gen.GenerateRequest]) (*connect.Response[gen.GenerateResponse], error) {
	p, err := s.store.GetProject(req.Msg.ProjectId)
	if err != nil {
		return nil, err
	}

	if err := s.generator.Generate(p); err != nil {
		return nil, err
	}

	return &connect.Response[gen.GenerateResponse]{
		Msg: &gen.GenerateResponse{
			ProjectId: p.Id,
		},
	}, nil
}
