package models

import "go.uber.org/zap/zapcore"

type User struct {
	Base
	AuthID                string `json:"auth_id"`
	CurrentOrganizationId string `json:"current_organization_id" gorm:"type:uuid;default:null"`
	// Organizations         []Organization `json:"organizations" gorm:"many2many:user_organizations;"`
	Username string `json:"username"`
	// TODO: This should be in a separate table but for now we'll just store it here
	DeviceToken string `json:"device_token"`
}

func (c *User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Id", c.ID.String())
	enc.AddString("UserId", c.ID.String())
	enc.AddString("CurrentOrganizationId", c.CurrentOrganizationId)
	enc.AddString("Username", c.Username)
	enc.AddString("DeviceToken", c.DeviceToken)
	return nil
}
