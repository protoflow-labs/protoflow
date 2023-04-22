package workflow

import (
	"github.com/google/wire"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/protoflow-labs/protoflow-editor/protoflow/gen"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/db"
	"github.com/protoflow-labs/protoflow-editor/protoflow/pkg/model"
	"gorm.io/gorm"
)

type Store interface {
	SaveWorkflow(w *gen.Workflow) (id string, err error)
	GetWorkflow(workflowID string) (protoflow *gen.Workflow, err error)
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
	err := db.AutoMigrate(&model.Workflow{})
	if err != nil {
		return nil, err
	}
	return &DBStore{
		db: db,
	}, nil
}

func (s *DBStore) SaveWorkflow(w *gen.Workflow) (id string, err error) {
	work := model.Workflow{
		Protoflow: model.ProtoJSON[interface{}]{
			Data: *w,
		},
	}
	res := s.db.Create(&work)
	if res.Error != nil {
		return "", res.Error
	}
	return work.ID.String(), nil
}

func (s *DBStore) GetWorkflow(workflowID string) (protoflow *gen.Workflow, err error) {
	w := model.Workflow{}
	res := s.db.First(&w, "id = ?", workflowID)
	if res.Error != nil {
		return nil, res.Error
	}
	return &w.Protoflow.Data, nil
}
