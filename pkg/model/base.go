package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Times struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type UUID struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
}

// BeforeCreate will set a UUID.
func (b *UUID) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}
