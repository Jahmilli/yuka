package models

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type Invitation struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;" json:"id" example:"aa22666c-0f57-45cb-a449-16efecc04f2e"`
	// FK id of the organization that the invitation is for
	// OrganizationId uuid.UUID `json:"organization_id" gorm:"type:uuid"`
	// FK id of the user that created the invitation
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:uuid"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;default:now()"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamptz;"`
}

// BeforeCreate populates the ID (if not set)
func (base *Invitation) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

func (c *Invitation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Id", c.ID.String())
	enc.AddString("CreatedBy", c.CreatedBy.String())
	enc.AddString("Token", c.Token)
	enc.AddString("CreatedAt", c.CreatedAt.String())
	enc.AddString("ExpiresAt", c.ExpiresAt.String())
	return nil
}
