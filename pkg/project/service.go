package project

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	genconnect "github.com/protoflow-labs/protoflow/gen/genconnect"
	"github.com/protoflow-labs/protoflow/pkg/workflow"
	"github.com/protoflow-labs/protoflow/templates"
	"google.golang.org/protobuf/types/known/anypb"
	"html/template"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	blockProtoTemplate, err := template.New("block").ParseFS(templates.Templates, "templates/*.template.proto")
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

	w, err := workflow.FromGraph(project.Graph)
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
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("project.ProjectService.GetProject is not implemented"))
}

func (s *Service) GetProjects(ctx context.Context, req *connect.Request[gen.GetProjectsRequest]) (*connect.Response[gen.GetProjectsResponse], error) {
	projects := make([]*gen.Project, 0)

	namespaces, err := s.clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, namespace := range namespaces.Items {
		projects = append(projects, &gen.Project{
			Id:   namespace.Name,
			Name: namespace.Name,
		})
	}

	return connect.NewResponse(&gen.GetProjectsResponse{Projects: projects}), nil
}

func (s *Service) CreateProject(context.Context, *connect.Request[gen.CreateProjectRequest]) (*connect.Response[gen.CreateProjectResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("project.ProjectService.CreateProject is not implemented"))
}

func (s *Service) DeleteProject(context.Context, *connect.Request[gen.DeleteProjectRequest]) (*connect.Response[gen.DeleteProjectResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("project.ProjectService.DeleteProject is not implemented"))
}

func (s *Service) GetBlocks(ctx context.Context, req *connect.Request[gen.GetBlocksRequest]) (*connect.Response[gen.GetBlocksResponse], error) {
	blocks := make([]*gen.Block, 0)

	files, err := os.ReadDir(".persistence/blocks")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, file := range files {
		dat, err := os.ReadFile(".persistence/blocks/" + file.Name())
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if len(dat) == 0 {
			continue
		}

		b := &gen.Block{}

		if err := json.Unmarshal(dat, b); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		blocks = append(blocks, b)
	}

	return connect.NewResponse(&gen.GetBlocksResponse{Blocks: blocks}), nil
}

func (s *Service) AddBlock(ctx context.Context, req *connect.Request[gen.AddBlockRequest]) (*connect.Response[gen.AddBlockResponse], error) {
	blockJson, _ := json.Marshal(req.Msg.Block)

	if err := os.WriteFile(".persistence/blocks/"+req.Msg.Block.Id+".dat", blockJson, 0644); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.generateProto(req.Msg.Block); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&gen.AddBlockResponse{Block: req.Msg.Block}), nil
}

func (s *Service) RemoveBlock(ctx context.Context, req *connect.Request[gen.RemoveBlockRequest]) (*connect.Response[gen.RemoveBlockResponse], error) {
	dat, err := os.ReadFile(".persistence/blocks/" + req.Msg.BlockId + ".dat")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b := &gen.Block{}

	if err := json.Unmarshal(dat, b); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	os.Remove(".persistence/blocks/" + req.Msg.BlockId + ".dat")
	os.Remove(".persistence/proto/" + b.Name + ".proto")

	return connect.NewResponse(&gen.RemoveBlockResponse{Block: b}), nil
}

func (s *Service) UpdateBlock(ctx context.Context, req *connect.Request[gen.UpdateBlockRequest]) (*connect.Response[gen.UpdateBlockResponse], error) {
	blockJson, _ := json.Marshal(req.Msg.Block)

	if err := os.WriteFile(".persistence/blocks/"+req.Msg.Block.Id+".dat", blockJson, 0644); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.generateProto(req.Msg.Block); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&gen.UpdateBlockResponse{Block: req.Msg.Block}), nil
}

func (s *Service) GetEdges(ctx context.Context, req *connect.Request[gen.GetEdgesRequest]) (*connect.Response[gen.GetEdgesResponse], error) {
	edges := make([]*gen.Edge, 0)

	files, err := os.ReadDir(".persistence/edges")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, file := range files {
		dat, err := os.ReadFile(".persistence/edges/" + file.Name())
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if len(dat) == 0 {
			continue
		}

		e := &gen.Edge{}

		if err := json.Unmarshal(dat, e); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		edges = append(edges, e)
	}

	return connect.NewResponse(&gen.GetEdgesResponse{Edges: edges}), nil
}

func (s *Service) AddEdge(ctx context.Context, req *connect.Request[gen.AddEdgeRequest]) (*connect.Response[gen.AddEdgeResponse], error) {
	edgeJson, _ := json.Marshal(req.Msg.Edge)

	if err := os.WriteFile(".persistence/edges/"+req.Msg.Edge.Id+".dat", edgeJson, 0644); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&gen.AddEdgeResponse{
		Edge: req.Msg.Edge,
	}), nil
}

func (s *Service) RemoveEdge(ctx context.Context, req *connect.Request[gen.RemoveEdgeRequest]) (*connect.Response[gen.RemoveEdgeResponse], error) {
	dat, err := os.ReadFile(".persistence/edges/" + req.Msg.EdgeId + ".dat")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	e := &gen.Edge{}

	if err := json.Unmarshal(dat, e); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	os.Remove(".persistence/edges/" + req.Msg.EdgeId + ".dat")

	return connect.NewResponse(&gen.RemoveEdgeResponse{Edge: e}), nil
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
