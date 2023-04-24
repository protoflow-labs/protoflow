package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	return nil
}
