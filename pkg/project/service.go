package project

import (
	"context"
	"encoding/json"
	"html/template"
	"os"

	"github.com/pkg/errors"
	genconnect "github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/protoflow-labs/protoflow/templates"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	store              Store
	manager            workflow.Manager
	clientset          *kubernetes.Clientset
	blockProtoTemplate *template.Template
}

var ProviderSet = wire.NewSet(
	StoreProviderSet,
	workflow.ProviderSet,
	NewService,
	wire.Bind(new(genconnect.ProjectServiceHandler), new(*Service)),
)

var _ genconnect.ProjectServiceHandler = (*Service)(nil)

func NewService(
	clientset *kubernetes.Clientset,
	store Store,
	manager workflow.Manager,
) (*Service, error) {
	// TODO breadchris this should be loading from an embedded file system
	blockProtoTemplate, err := template.New("block").ParseFS(templates.Templates, "*.template.proto")
	if err != nil {
		return nil, err
	}

	return &Service{
		store:              store,
		manager:            manager,
		clientset:          clientset,
		blockProtoTemplate: blockProtoTemplate,
	}, nil
}

func (s *Service) GetResources(ctx context.Context, c *connect.Request[gen.GetResourcesRequest]) (*connect.Response[gen.GetResourcesResponse], error) {
	//TODO implement me
	panic("implement me")
}

func resultToAny(res *workflow.Result) (*anypb.Any, error) {
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	output, err := anypb.New(&gen.Result{
		Data: data,
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (s *Service) RunWorklow(ctx context.Context, c *connect.Request[gen.RunWorkflowRequest]) (*connect.Response[gen.RunOutput], error) {
	project, err := s.store.GetProject(c.Msg.ProjectId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project %s", c.Msg.ProjectId)
	}

	w, err := workflow.FromProject(project)
	if err != nil {
		return nil, err
	}

	res, err := s.manager.ExecuteWorkflowSync(ctx, w, c.Msg.NodeId)
	if err != nil {
		return nil, err
	}

	output, err := resultToAny(res)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.RunOutput{
		Output: output,
	}), nil
}

func (s *Service) RunBlock(ctx context.Context, c *connect.Request[gen.RunBlockRequest]) (*connect.Response[gen.RunOutput], error) {
	//TODO implement me
	panic("implement me")
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

	return connect.NewResponse(&gen.GetProjectsResponse{Projects: projects}), nil
}

func (s *Service) CreateProject(ctx context.Context, req *connect.Request[gen.CreateProjectRequest]) (*connect.Response[gen.CreateProjectResponse], error) {
	project := gen.Project{
		Id:   uuid.NewString(),
		Name: req.Msg.Name,
	}
	_, err := s.store.SaveProject(&project)

	if err != nil {
		return connect.NewResponse(&gen.CreateProjectResponse{Project: nil}), nil
	}

	return connect.NewResponse(&gen.CreateProjectResponse{Project: &project}), nil
}

func (s *Service) DeleteProject(context.Context, *connect.Request[gen.DeleteProjectRequest]) (*connect.Response[gen.DeleteProjectResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("project.ProjectService.DeleteProject is not implemented"))
}

func (s *Service) SaveProject(ctx context.Context, req *connect.Request[gen.SaveProjectRequest]) (*connect.Response[gen.SaveProjectResponse], error) {
	s.store.SaveProject(&gen.Project{
		Id:    req.Msg.ProjectId,
		Name:  "local",
		Graph: req.Msg.Graph,
	})
	return nil, nil
}

func (s *Service) generateProto(block *gen.Block) error {
	file, err := os.Create(".persistence/proto/" + block.Name + ".proto")
	if err != nil {
		return err
	}

	defer file.Close()

	err = s.blockProtoTemplate.Execute(file, block)
	if err != nil {
		return err
	}

	return nil
}
