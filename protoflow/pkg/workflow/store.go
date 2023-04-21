package workflow

import (
	"github.com/breadchris/protoflow/gen/workflow"
	"github.com/breadchris/protoflow/pkg/db"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type Store interface {
	SaveDeploymentWorkflows(deploymentID string, workflows map[string]Workflow) (err error)
	SaveWorkflowRun(workflowID, runID string)
	GetWorkflow(workflowID string) (protoflow *workflow.Workflow, err error)
	GetWorkflowRunsForDeployment(deploymentID string) (workflowRuns []WorkflowRunModel, err error)
	DeleteDeploymentWorkflows(deploymentID string) (err error)
}

type DBStore struct {
	db *gorm.DB
}

type WorkflowRunModel struct {
	WorkflowID string
	RunID      string
}

var DBProviderSet = wire.NewSet(
	db.NewGormDB,
	NewDBStore,
)

func NewDBStore(db *gorm.DB) *DBStore {
	return &DBStore{
		db: db,
	}
}

func (s *DBStore) SaveDeploymentWorkflows(deploymentID string, workflows map[string]Workflow) (err error) {
	return
}

func (s *DBStore) SaveWorkflowRun(workflowID, runID string) {
}

func (s *DBStore) GetWorkflow(workflowID string) (protoflow *workflow.Workflow, err error) {
	return nil, nil
}

func (s *DBStore) GetWorkflowRunsForDeployment(deploymentID string) (workflowRuns []WorkflowRunModel, err error) {
	return nil, nil
}

func (s *DBStore) DeleteDeploymentWorkflows(deploymentID string) (err error) {
	return
}
