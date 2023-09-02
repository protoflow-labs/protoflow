package workflow

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/rs/zerolog/log"
	"sync"
)

type RunningWorkflow struct {
	Trace     *gen.WorkflowTrace
	Ctx       context.Context
	Cancel    context.CancelFunc
	connector *Connector
}

type WorkflowManager struct {
	mu        sync.RWMutex
	workflows map[string]*RunningWorkflow
}

var ProviderSet = wire.NewSet(
	NewManagerBuilder,
	NewWorkflowManager,
)

func NewWorkflowManager() *WorkflowManager {
	return &WorkflowManager{
		workflows: map[string]*RunningWorkflow{},
	}
}

func (m *WorkflowManager) Start(ctx context.Context, id string, req *gen.RunWorkflowRequest) (context.Context, *Connector) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO breadchris configurable timeout?
	ctx, cancel := context.WithCancel(ctx)
	c := NewConnector()
	m.workflows[id] = &RunningWorkflow{
		Trace: &gen.WorkflowTrace{
			Id:      id,
			Request: req,
		},
		Ctx:       ctx,
		Cancel:    cancel,
		connector: c,
	}
	return ctx, c
}

func (m *WorkflowManager) Stop(id string) error {
	// TODO breadchris need to close the input stream for the workflow

	m.mu.Lock()
	defer m.mu.Unlock()

	w, exists := m.workflows[id]
	if !exists {
		return errors.New("workflow not found")
	}

	// TODO breadchris is this correct?
	// cleanup observers
	w.connector.Dispose()

	log.Debug().Str("id", id).Msg("stopping workflow")
	w.Cancel()
	delete(m.workflows, id)
	return nil
}

func (m *WorkflowManager) Traces() []*gen.WorkflowTrace {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var traces []*gen.WorkflowTrace
	for _, w := range m.workflows {
		traces = append(traces, w.Trace)
	}
	return traces
}

type ManagerBuilder struct {
	wm  *WorkflowManager
	m   *Manager
	err error
}

type Manager struct {
	id  string
	req *gen.RunWorkflowRequest

	wm *WorkflowManager
}

func NewManagerBuilder(wm *WorkflowManager) *ManagerBuilder {
	return &ManagerBuilder{
		m:  NewManager(wm),
		wm: wm,
	}
}

func (m *ManagerBuilder) WithReq(req *gen.RunWorkflowRequest) *ManagerBuilder {
	nm := *m
	nm.m.req = req
	return &nm
}

func (m *ManagerBuilder) Build() (*Manager, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.m, nil
}

func NewManager(wm *WorkflowManager) *Manager {
	return &Manager{
		wm: wm,
		id: uuid.NewString(),
	}
}

func (m *Manager) Id() string {
	return m.id
}

func (m *Manager) Start(ctx context.Context) (context.Context, *Connector) {
	return m.wm.Start(ctx, m.id, m.req)
}

func (m *Manager) Stop() error {
	return m.wm.Stop(m.id)
}
