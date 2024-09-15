package models

import "go.uber.org/zap/zapcore"

type Organization struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`
	// Users       []*User `json:"-" gorm:"many2many:user_organizations;"`
}

func (c *Organization) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Id", c.ID.String())
	enc.AddString("Name", c.Name)
	enc.AddString("Description", c.Description)
	enc.AddTime("CreatedAt", c.CreatedAt)
	enc.AddTime("UpdatedAt", c.UpdatedAt)
	enc.AddTime("DeletedAt", c.DeletedAt)
	return nil
}
