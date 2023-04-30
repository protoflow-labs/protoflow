package db

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/google/wire"
	"github.com/protoflow-labs/protoflow/pkg/cache"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewGormDB,
)

func NewGormDB(config Config, cache cache.Cache) (*gorm.DB, error) {
	var openedDb gorm.Dialector
	if config.Driver == "postgres" {
		openedDb = postgres.Open(config.DSN)
	} else if config.Driver == "sqlite" {
		dbPath, err := cache.GetFile(config.DSN)
		if err != nil {
			return nil, err
		}
		openedDb = sqlite.Open(dbPath + "?cache=shared&mode=rwc&_journal_mode=WAL")
	} else {
		return nil, fmt.Errorf("unknown db driver: %s", config.Driver)
	}

	db, err := gorm.Open(openedDb, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// TODO breadchris automigrate

	return db.Debug(), nil
}
