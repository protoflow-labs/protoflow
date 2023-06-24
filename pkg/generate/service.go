package generate

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/pkg/errors"
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

func (s *Service) generatorAndProject(projectID string) (*Generate, *project.Project, error) {
	projProto, err := s.store.GetProject(projectID)
	if err != nil {
		return nil, nil, err
	}

	p, err := project.FromProto(projProto)
	if err != nil {
		return nil, nil, err
	}

	generator, err := NewGenerate(s.config)
	if err != nil {
		return nil, nil, err
	}
	return generator, p, nil
}

func (s *Service) GenerateImplementation(ctx context.Context, c *connect.Request[gen.GenerateImplementationRequest]) (*connect.Response[gen.GenerateImplementationResponse], error) {
	generator, p, err := s.generatorAndProject(c.Msg.ProjectId)
	if err != nil {
		return nil, err
	}

	n, ok := p.Workflow.NodeLookup[c.Msg.NodeId]
	if !ok {
		return nil, errors.Wrapf(err, "node %s not found", c.Msg.NodeId)
	}

	if err := generator.GenerateImplementation(p, n); err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GenerateImplementationResponse{}), nil
}

func (s *Service) InferNodeType(ctx context.Context, c *connect.Request[gen.InferNodeTypeRequest]) (*connect.Response[gen.InfertNodeTypeResponse], error) {
	generator, p, err := s.generatorAndProject(c.Msg.ProjectId)
	if err != nil {
		return nil, err
	}

	n, ok := p.Workflow.NodeLookup[c.Msg.NodeId]
	if !ok {
		return nil, errors.Wrapf(err, "node %s not found", c.Msg.NodeId)
	}

	if err := generator.InferNodeType(p, n); err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.InfertNodeTypeResponse{}), nil
}

func (s *Service) Generate(ctx context.Context, req *connect.Request[gen.GenerateRequest]) (*connect.Response[gen.GenerateResponse], error) {
	generator, p, err := s.generatorAndProject(req.Msg.ProjectId)
	if err != nil {
		return nil, err
	}

	if err := generator.Generate(p); err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GenerateResponse{ProjectId: p.Base.Id}), nil
}
