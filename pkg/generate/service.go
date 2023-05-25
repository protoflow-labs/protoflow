package generate

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	project "github.com/protoflow-labs/protoflow/pkg/project"
	"github.com/protoflow-labs/protoflow/pkg/store"
)

type Service struct {
	config Config
	store  store.Project
}

var _ genconnect.GenerateServiceHandler = (*Service)(nil)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewService,
	wire.Bind(new(genconnect.GenerateServiceHandler), new(*Service)),
)

func NewService(config Config, store store.Project) (*Service, error) {
	return &Service{
		config: config,
		store:  store,
	}, nil
}

func (s *Service) Generate(ctx context.Context, req *connect.Request[gen.GenerateRequest]) (*connect.Response[gen.GenerateResponse], error) {
	projProto, err := s.store.GetProject(req.Msg.ProjectId)
	if err != nil {
		return nil, err
	}

	p, err := project.FromProto(projProto)
	if err != nil {
		return nil, err
	}

	generator, err := NewGenerate(s.config)
	if err != nil {
		return nil, err
	}
	if err := generator.Generate(p); err != nil {
		return nil, err
	}

	return &connect.Response[gen.GenerateResponse]{
		Msg: &gen.GenerateResponse{
			ProjectId: p.Base.Id,
		},
	}, nil
}
