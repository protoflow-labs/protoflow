package project

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/graph"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
)

type Project struct {
	Base     *gen.Project
	Workflow *workflow.Workflow
}

func FromProto(project *gen.Project) (*Project, error) {
	w, err := workflow.Default().
		WithProtoProject(graph.ConvertProto(project)).
		Build()
	if err != nil {
		return nil, err
	}
	return &Project{
		Base:     project,
		Workflow: w,
	}, nil
}

func (s *Service) GetProject(context.Context, *connect.Request[gen.GetProjectRequest]) (*connect.Response[gen.GetProjectResponse], error) {
	proj, err := s.store.GetProject("local")
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GetProjectResponse{Project: proj}), nil
}

func (s *Service) GetProjects(ctx context.Context, req *connect.Request[gen.GetProjectsRequest]) (*connect.Response[gen.GetProjectsResponse], error) {
	projects, err := s.store.ListProjects()
	if err != nil {
		return nil, err
	}

	var p []*gen.GetProjectResponse
	projectTypes, err := getProjectTypes()
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		p = append(p, &gen.GetProjectResponse{
			Project: project,
			Types:   projectTypes,
		})
	}

	return connect.NewResponse(&gen.GetProjectsResponse{Projects: p}), nil
}

func (s *Service) CreateProject(ctx context.Context, req *connect.Request[gen.CreateProjectRequest]) (*connect.Response[gen.CreateProjectResponse], error) {
	_, err := s.store.CreateProject(s.defaultProject)

	if err != nil {
		return connect.NewResponse(&gen.CreateProjectResponse{Project: nil}), nil
	}

	return connect.NewResponse(&gen.CreateProjectResponse{Project: s.defaultProject}), nil
}

func (s *Service) DeleteProject(context.Context, *connect.Request[gen.DeleteProjectRequest]) (*connect.Response[gen.DeleteProjectResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("project.ProjectService.DeleteProject is not implemented"))
}

func (s *Service) SaveProject(ctx context.Context, req *connect.Request[gen.SaveProjectRequest]) (*connect.Response[gen.SaveProjectResponse], error) {
	project, err := s.store.GetProject(req.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", req.Msg.ProjectId)
	}

	project.Graph = req.Msg.Graph

	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save project %s", project.Id)
	}

	return connect.NewResponse(&gen.SaveProjectResponse{Project: project}), nil
}

func (s *Service) NewNode(ctx context.Context, c *connect.Request[gen.NewNodeRequest]) (*connect.Response[gen.NewNodeResponse], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	n := c.Msg.Node
	n.Id = uuid.NewString()
	project.Graph.Nodes = append(project.Graph.Nodes, n)

	_, err = s.store.SaveProject(project)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to save project %s", project.Id)
	}

	return connect.NewResponse(&gen.NewNodeResponse{Node: n}), nil
}
