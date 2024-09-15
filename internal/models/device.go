package models

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

type Device struct {
	Base
	PeerId     uuid.UUID `json:"peer_id" gorm:"type:uuid"`
	MacAddress string    `json:"mac_address"`
	Hostname   string    `json:"hostname"`
	LastSeen   time.Time `json:"last_seen" gorm:"type:timestamptz;null"`
	LocalIp    string    `json:"local_ip"`
}

func (c *Device) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Id", c.ID.String())
	enc.AddString("PeerId", c.PeerId.String())
	enc.AddString("MacAddress", c.MacAddress)
	enc.AddString("Hostname", c.Hostname)
	enc.AddTime("LastSeen", c.LastSeen)
	enc.AddString("LocalIp", c.LocalIp)
	return nil
}
