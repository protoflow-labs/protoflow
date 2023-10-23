package project

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/bucket"
	"github.com/protoflow-labs/protoflow/pkg/grpc"
	"github.com/protoflow-labs/protoflow/pkg/store"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/rs/zerolog/log"
)

type Service struct {
	store          store.Project
	cache          bucket.Bucket
	defaultProject *gen.Project
	manager        *workflow.ManagerBuilder
	// TODO breadchris rename this to something that is more relevant
	workflowManager *workflow.WorkflowManager
}

var ProviderSet = wire.NewSet(
	store.ProviderSet,
	NewService,
	workflow.ProviderSet,
	wire.Bind(new(genconnect.ProjectServiceHandler), new(*Service)),
)

var _ genconnect.ProjectServiceHandler = (*Service)(nil)

func NewDefaultProject(cache bucket.Bucket) (*gen.Project, error) {
	// TODO breadchris this folder should be configurable
	bucketDir, err := cache.GetFolder("filestore")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get bucket dir")
	}
	return getDefaultProject("local", bucketDir), nil
}

func NewService(
	store store.Project,
	cache bucket.Bucket,
	defaultProject *gen.Project,
	manager *workflow.ManagerBuilder,
	workflowManager *workflow.WorkflowManager,
) (*Service, error) {
	return &Service{
		store:           store,
		cache:           cache,
		defaultProject:  defaultProject,
		manager:         manager,
		workflowManager: workflowManager,
	}, nil
}

func (s *Service) EnumerateProviders(ctx context.Context, c *connect.Request[gen.GetProvidersRequest]) (*connect.Response[gen.GetProvidersResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := FromProto(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	return connect.NewResponse(&gen.GetProvidersResponse{
		Providers: w.EnumerateProviders(),
	}), nil
}

func getProjectTypes() (*gen.ProjectTypes, error) {
	// TODO breadchris when types are bound to a project, this should be specific to a project
	// return the rules for different layers
	n := &gen.Node{}
	nd, err := grpc.SerializeType(n)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to wrap message")
	}
	e := &gen.Edge{}
	ed, err := grpc.SerializeType(e)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to wrap message")
	}

	// TODO breadchris cleanup this code, see blocks.go:76
	tr := grpc.NewTypeResolver()
	tr = tr.ResolveLookup(n)
	tr = tr.ResolveLookup(e)

	sr := tr.Serialize()

	return &gen.ProjectTypes{
		NodeType:   nd.AsDescriptorProto(),
		EdgeType:   ed.AsDescriptorProto(),
		DescLookup: sr.DescLookup,
		EnumLookup: sr.EnumLookup,
	}, nil
}

func maybeDefaultPath(path string) string {
	if path == "" {
		return "protoflow.bin"
	}
	return path
}

func (s *Service) ExportProject(ctx context.Context, c *connect.Request[gen.ExportProjectRequest]) (*connect.Response[gen.ExportProjectResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	return connect.NewResponse(
		&gen.ExportProjectResponse{},
	), SaveToFile(project, maybeDefaultPath(c.Msg.Path))
}

func (s *Service) LoadProject(ctx context.Context, c *connect.Request[gen.LoadProjectRequest]) (*connect.Response[gen.LoadProjectResponse], error) {
	project, err := LoadFromFile(maybeDefaultPath(c.Msg.Path))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load project from %s", c.Msg.Path)
	}
	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save project %s", project.Id)
	}
	return connect.NewResponse(
		&gen.LoadProjectResponse{
			Project: project,
		},
	), nil
}

func (s *Service) GetProjectTypes(ctx context.Context, c *connect.Request[gen.GetProjectTypesRequest]) (*connect.Response[gen.ProjectTypes], error) {
	projectTypes, err := getProjectTypes()
	return &connect.Response[gen.ProjectTypes]{
		Msg: projectTypes,
	}, err
}

func (s *Service) GetNodeInfo(ctx context.Context, c *connect.Request[gen.GetNodeInfoRequest]) (*connect.Response[gen.GetNodeInfoResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := FromProto(project)
	if err != nil {
		return nil, err
	}
	n, err := w.GetNode(c.Msg.NodeId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node %s", c.Msg.NodeId)
	}
	nodeInfo, err := w.GetNodeInfo(n)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get node info for node %s", c.Msg.NodeId)
	}
	if nodeInfo == nil {
		log.Warn().Str("node", c.Msg.NodeId).Msg("node has no info")
		return connect.NewResponse(&gen.GetNodeInfoResponse{}), nil
	}
	if nodeInfo.Method == nil {
		log.Warn().Str("node", c.Msg.NodeId).Msg("node has no method")
		return connect.NewResponse(&gen.GetNodeInfoResponse{}), nil
	}
	typeInfo, err := nodeInfo.Method.Proto()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get proto for node %s", c.Msg.NodeId)
	}
	proto, err := nodeInfo.BuildProto()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build proto for node %s", c.Msg.NodeId)
	}
	return connect.NewResponse(&gen.GetNodeInfoResponse{
		MethodProto: proto,
		TypeInfo:    typeInfo,
	}), nil
}
