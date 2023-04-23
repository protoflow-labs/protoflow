package project

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	genconnect.UnimplementedProjectServiceHandler

	clientset *kubernetes.Clientset
}

var ProviderSet = wire.NewSet(
	NewService,
	wire.Bind(new(genconnect.ProjectServiceHandler), new(*Service)),
)

func NewService(clientset *kubernetes.Clientset) *Service {
	return &Service{
		clientset: clientset,
	}
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

	files, err := os.ReadDir("blocks")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, file := range files {
		dat, err := os.ReadFile("blocks/" + file.Name())
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
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

	if err := os.WriteFile("blocks/"+req.Msg.Block.Id+".dat", blockJson, 0644); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&gen.AddBlockResponse{
		Block: req.Msg.Block,
	}), nil
}

func (s *Service) RemoveBlock(ctx context.Context, req *connect.Request[gen.RemoveBlockRequest]) (*connect.Response[gen.RemoveBlockResponse], error) {
	dat, err := os.ReadFile("blocks/" + req.Msg.BlockId + ".dat")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	b := &gen.Block{}

	if err := json.Unmarshal(dat, b); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	os.Remove("blocks/" + req.Msg.BlockId + ".dat")

	return connect.NewResponse(&gen.RemoveBlockResponse{Block: b}), nil
}

func (s *Service) UpdateBlock(ctx context.Context, req *connect.Request[gen.UpdateBlockRequest]) (*connect.Response[gen.UpdateBlockResponse], error) {
	blockJson, _ := json.Marshal(req.Msg.Block)

	if err := os.WriteFile("blocks/"+req.Msg.Block.Id+".dat", blockJson, 0644); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&gen.UpdateBlockResponse{Block: req.Msg.Block}), nil
}

func (s *Service) GetEdges(ctx context.Context, req *connect.Request[gen.GetEdgesRequest]) (*connect.Response[gen.GetEdgesResponse], error) {
	edges := make([]*gen.Edge, 0)

	files, err := os.ReadDir("edges")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, file := range files {
		dat, err := os.ReadFile("edges/" + file.Name())
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
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

	if err := os.WriteFile("edges/"+req.Msg.Edge.Id+".dat", edgeJson, 0644); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&gen.AddEdgeResponse{
		Edge: req.Msg.Edge,
	}), nil
}

func (s *Service) RemoveEdge(ctx context.Context, req *connect.Request[gen.RemoveEdgeRequest]) (*connect.Response[gen.RemoveEdgeResponse], error) {
	dat, err := os.ReadFile("edges/" + req.Msg.EdgeId + ".dat")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	e := &gen.Edge{}

	if err := json.Unmarshal(dat, e); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	os.Remove("edges/" + req.Msg.EdgeId + ".dat")

	return connect.NewResponse(&gen.RemoveEdgeResponse{Edge: e}), nil
}
