package workflow

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/gen/genconnect"
)

var ProviderSet = wire.NewSet(
	StoreProviderSet,
	NewConfig,
	NewService,
	NewManager,
	wire.Bind(new(genconnect.ManagerServiceHandler), new(*Service)),
)

var _ genconnect.ManagerServiceHandler = (*Service)(nil)

type Service struct {
	genconnect.UnimplementedManagerServiceHandler

	workflowStore Store
	manager       Manager
}

func NewService(
	store Store,
	manager Manager,
) *Service {
	return &Service{
		workflowStore: store,
		manager:       manager,
	}
}

func (m *Service) CreateWorkflow(
	ctx context.Context,
	req *connect.Request[gen.Workflow],
) (*connect.Response[gen.ID], error) {
	id, err := m.workflowStore.SaveWorkflow(req.Msg)
	if err != nil {
		return nil, err
	}
	return &connect.Response[gen.ID]{
		Msg: &gen.ID{
			Id: id,
		},
	}, nil
}

func (m *Service) StartWorkflow(
	ctx context.Context,
	req *connect.Request[gen.WorkflowEntrypoint],
) (*connect.Response[gen.Run], error) {
	protoflow, err := m.workflowStore.GetWorkflow(req.Msg.WorkflowId)
	if err != nil {
		return nil, err
	}

	w, err := WorkflowFromProtoflow(protoflow)
	if err != nil {
		return nil, err
	}

	runID, err := m.manager.ExecuteWorkflow(ctx, w, req.Msg.NodeId)

	return &connect.Response[gen.Run]{
		Msg: &gen.Run{
			Id: runID,
		},
	}, nil
}
