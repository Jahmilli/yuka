package models

import (
	"go.uber.org/zap/zapcore"
)

type RegisteredApplication struct {
	Base
	ApplicationId  string `json:"application_id" gorm:"type:uuid"`
	OrganizationId string `json:"organization_id" gorm:"type:uuid"`
	// TODO: Maybe the status' should be json, idk for now...
	DaemonReady      bool `json:"daemon_status"`
	ApplicationReady bool `json:"application_status"`
}

func (c *RegisteredApplication) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Id", c.ID.String())
	enc.AddString("ApplicationId", c.ApplicationId)
	enc.AddString("OrganizationId", c.OrganizationId)
	return nil
}
