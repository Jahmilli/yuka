package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id" example:"aa22666c-0f57-45cb-a449-16efecc04f2e"`
	CreatedAt time.Time `json:"-" gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `json:"-" gorm:"type:timestamptz;default:now()"`
	DeletedAt time.Time `gorm:"type:timestamptz;null;default:null" json:"-"`
}

// BeforeCreate populates the ID (if not set)
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}
