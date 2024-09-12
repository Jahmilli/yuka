package models

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

type DeviceDNSQuery struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;" json:"-" example:"aa22666c-0f57-45cb-a449-16efecc04f2e"`
	PeerID uuid.UUID `json:"-" gorm:"type:uuid"`
	// To minimise serialization cost slightly, chaning domain and query time to single characters
	Domain    string    `json:"d" gorm:"type:string"`
	QueryTime time.Time `json:"t" gorm:"type:timestamptz"`
}

func (c *DeviceDNSQuery) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Id", c.ID.String())
	enc.AddString("PeerId", c.PeerID.String())
	enc.AddString("Domain", c.Domain)
	enc.AddTime("QueryTime", c.QueryTime)
	return nil
}
