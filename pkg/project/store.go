package project

import (
	"github.com/google/uuid"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/protoflow-labs/protoflow/gen"
	"github.com/protoflow-labs/protoflow/pkg/db"
	"github.com/protoflow-labs/protoflow/pkg/model"
	"gorm.io/gorm"
)

type Store interface {
	CreateProject(w *gen.Project) (string, error)
	SaveProject(w *gen.Project) (string, error)
	GetProject(projectID string) (*gen.Project, error)
	ListProjects() ([]*gen.Project, error)
}

var _ Store = (*DBStore)(nil)

type DBStore struct {
	db *gorm.DB
}

type WorkflowRunModel struct {
	WorkflowID string
	RunID      string
}

var StoreProviderSet = wire.NewSet(
	db.ProviderSet,
	NewDBStore,
	wire.Bind(new(Store), new(*DBStore)),
)

func NewDBStore(db *gorm.DB) (*DBStore, error) {
	err := db.AutoMigrate(&model.Project{})
	if err != nil {
		return nil, err
	}
	return &DBStore{
		db: db,
	}, nil
}

func (s *DBStore) CreateProject(w *gen.Project) (string, error) {
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

func (s *DBStore) SaveProject(w *gen.Project) (string, error) {
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

func (s *DBStore) GetProject(projectID string) (*gen.Project, error) {
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

func (s *DBStore) ListProjects() ([]*gen.Project, error) {
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
