package generate

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/cache"
	"github.com/protoflow-labs/protoflow/pkg/project"
	"path"
)

type Service struct {
	config Config
	store  project.Store
}

var _ genconnect.GenerateServiceHandler = (*Service)(nil)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewService,
	wire.Bind(new(genconnect.GenerateServiceHandler), new(*Service)),
)

func NewService(config Config, store project.Store) (*Service, error) {
	return &Service{
		config: config,
		store:  store,
	}, nil
}

func (s *Service) Generate(ctx context.Context, req *connect.Request[gen.GenerateRequest]) (*connect.Response[gen.GenerateResponse], error) {
	p, err := s.store.GetProject(req.Msg.ProjectId)
	if err != nil {
		return nil, err
	}

	var projectDir string
	// if there is a project path defined, use this for where the codeBucket goes
	if s.config.ProjectPath != "" {
		projectDir = s.config.ProjectPath
	} else {
		projectDir = path.Join("projects", p.Name)
	}
	c, err := cache.FromDir(projectDir)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating codeBucket from %s", projectDir)
	}

	generator := NewGenerate(c)
	if err := generator.Generate(p); err != nil {
		return nil, err
	}

	return &connect.Response[gen.GenerateResponse]{
		Msg: &gen.GenerateResponse{
			ProjectId: p.Id,
		},
	}, nil
}
