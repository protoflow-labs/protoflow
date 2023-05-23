package model

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/protoflow-labs/protoflow/gen"

	"google.golang.org/protobuf/encoding/protojson"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type WorkflowRunJSON struct {
	// TODO breadchris couldn't figure out how to make this generic, there is a problem with protojson.Unmarshal/Marshal
	Data *gen.WorkflowRun
}

// Value return json value, implement driver.Valuer interface
func (j *WorkflowRunJSON) Value() (driver.Value, error) {
	return j.MarshalJSON()
}

// Scan scan value into WorkflowRunJSON[T], implements sql.Scanner interface
func (j *WorkflowRunJSON) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return j.UnmarshalJSON(bytes)
}

func (j *WorkflowRunJSON) MarshalJSON() ([]byte, error) {
	marshaler := &protojson.MarshalOptions{}
	b, err := marshaler.Marshal(j.Data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (j *WorkflowRunJSON) UnmarshalJSON(data []byte) error {
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}

	if j.Data == nil {
		j.Data = &gen.WorkflowRun{}
	}

	if err := unmarshaler.Unmarshal(data, j.Data); err != nil {
		return err
	}
	return nil
}

// GormDataType gorm common data type
func (*WorkflowRunJSON) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (*WorkflowRunJSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (j *WorkflowRunJSON) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := j.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}
