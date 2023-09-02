package store

import (
	"github.com/google/uuid"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/db"
	"github.com/protoflow-labs/protoflow/pkg/model"
	"gorm.io/gorm"
)

type Project interface {
	CreateProject(w *gen.Project) (string, error)
	SaveProject(w *gen.Project) (string, error)
	GetProject(projectID string) (*gen.Project, error)
	ListProjects() ([]*gen.Project, error)
	CreateWorkflowRun(run *gen.WorkflowTrace) (string, error)
	SaveNodeExecution(workflowRunID string, nodeExecution *gen.NodeExecution) (string, error)
	GetWorkflowRunsForProject(project string) ([]*gen.WorkflowTrace, error)
}

var _ Project = (*ProjectStore)(nil)

type ProjectStore struct {
	db *gorm.DB
}

type WorkflowRunModel struct {
	WorkflowID string
	RunID      string
}

var ProviderSet = wire.NewSet(
	db.ProviderSet,
	NewDBStore,
	wire.Bind(new(Project), new(*ProjectStore)),
)

func NewDBStore(db *gorm.DB) (*ProjectStore, error) {
	err := db.AutoMigrate(
		&model.Project{},
		&model.WorkflowRun{},
		&model.NodeExecution{},
	)
	if err != nil {
		return nil, err
	}
	return &ProjectStore{
		db: db,
	}, nil
}

func (s *ProjectStore) CreateProject(w *gen.Project) (string, error) {
	project := model.Project{
		UUID: model.UUID{
			ID: uuid.MustParse(w.Id),
		},
		ProjectJSON: &model.ProjectJSON{
			Data: w,
		},
	}

	res := s.db.Create(&project)
	if res.Error != nil {
		return "", res.Error
	}

	return project.ID.String(), nil
}

func (s *ProjectStore) SaveProject(w *gen.Project) (string, error) {
	project := model.Project{
		UUID: model.UUID{
			ID: uuid.MustParse(w.Id),
		},
		ProjectJSON: &model.ProjectJSON{
			Data: w,
		},
	}

	res := s.db.Save(&project)
	if res.Error != nil {
		return "", res.Error
	}

	return project.ID.String(), nil
}

func (s *ProjectStore) GetProject(projectID string) (*gen.Project, error) {
	w := model.Project{}
	res := s.db.First(&w, "id = ?", projectID)
	if res.Error != nil {
		return nil, res.Error
	}

	if w.Data == nil {
		return nil, nil
	}

	return w.Data, nil
}

func (s *ProjectStore) ListProjects() ([]*gen.Project, error) {
	var projects []*model.Project

	res := s.db.Find(&projects)
	if res.Error != nil {
		return nil, res.Error
	}

	var result []*gen.Project
	for _, p := range projects {
		result = append(result, p.Data)
	}

	return result, nil
}

func (s *ProjectStore) CreateWorkflowRun(run *gen.WorkflowTrace) (string, error) {
	w := model.WorkflowRun{
		UUID: model.UUID{
			ID: uuid.MustParse(run.Id),
		},
		ProjectID: run.Request.ProjectId,
		WorkflowRunJSON: &model.WorkflowRunJSON{
			Data: run,
		},
	}

	res := s.db.Create(&w)
	if res.Error != nil {
		return "", res.Error
	}

	return w.ID.String(), nil
}

func (s *ProjectStore) SaveNodeExecution(workflowRunID string, nodeExecution *gen.NodeExecution) (string, error) {
	w := model.NodeExecution{
		UUID: model.UUID{
			ID: uuid.New(),
		},
		WorkflowRunID: workflowRunID,
		NodeExecutionJSON: &model.NodeExecutionJSON{
			Data: nodeExecution,
		},
	}

	res := s.db.Create(&w)
	if res.Error != nil {
		return "", res.Error
	}

	return w.ID.String(), nil
}

func (s *ProjectStore) GetWorkflowRunsForProject(projectID string) ([]*gen.WorkflowTrace, error) {
	var runs []*model.WorkflowRun
	res := s.db.Find(&runs, "project_id = ?", projectID)
	if res.Error != nil {
		return nil, res.Error
	}

	var result []*gen.WorkflowTrace
	for _, r := range runs {
		var execs []*model.NodeExecution
		res := s.db.Find(&execs, "workflow_run_id = ?", r.ID)
		if res.Error != nil {
			return nil, errors.Wrapf(res.Error, "failed to get node executions for workflow run %s", r.ID)
		}
		for _, e := range execs {
			r.Data.NodeExecs = append(r.Data.NodeExecs, e.Data)
		}
		result = append(result, r.Data)
	}

	return result, nil
}
