package handlers

import (
	"fmt"
	"yuka/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateUserInput struct {
	AuthID                string `json:"auth_id" binding:"required"`
	CurrentOrganizationId string `json:"current_organization_id" binding:"required,uuid4"`
	Username              string `json:"username" binding:"required"`
	DeviceToken           string `json:"device_token" binding:"required"`
}

type UpdateUserInput struct {
	Username              string `json:"username"`
	DeviceToken           string `json:"device_token"`
	CurrentOrganizationId string `json:"current_organization_id"`
}

type UserKey string

const (
	FindUserKeyUserID      UserKey = "id"
	FindUserKeyAuthID      UserKey = "auth_id"
	FindUserKeyUsername    UserKey = "username"
	FindUserKeyDeviceToken UserKey = "device_token"
)

func (k UserKey) String() string {
	return string(k)
}

type UserHandler struct {
	Db     *gorm.DB
	Logger *zap.Logger
}

func NewUserHandler(logger *zap.Logger, db *gorm.DB) UserHandler {
	return UserHandler{
		Db:     db,
		Logger: logger,
	}
}

func (c *UserHandler) FindUsers(key UserKey, val string) ([]models.User, error) {
	var users []models.User
	if err := c.Db.Where(fmt.Sprintf("%s = ?", key), val).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// CreateUser creates a User if it doesnt exist otherwise returns the existing user from the database
func (c *UserHandler) CreateUser(input CreateUserInput) (*models.User, error) {
	var existingUser models.User
	if err := c.Db.Where("auth_id = ?", input.AuthID).First(&existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		// If we have no user then we need to create one and so we can proceed
	} else {
		// We already have a user so just return it
		return &existingUser, nil
	}

	user := models.User{
		AuthID:      input.AuthID,
		Username:    input.Username,
		DeviceToken: input.DeviceToken,
	}

	c.Logger.Info("Created user", zap.Object("user", &user))
	return &user, nil
}

// FindUser returns the user along with their associated organizations
func (c *UserHandler) FindUser(key UserKey, val string) (*models.User, error) {
	var user models.User

	if err := c.Db.Where(fmt.Sprintf("%s = ?", key), val).Preload("Organizations").First(&user).Error; err != nil {
		return nil, err
	}

	c.Logger.Info("Found user", zap.Object("user", &user))
	return &user, nil
}

func (c *UserHandler) UpdateUser(id string, input UpdateUserInput) (*models.User, error) {
	var user models.User
	slogger := c.Logger.Sugar()
	slogger.Debugf("Updating user with id %s and input %v", id, input)
	if err := c.Db.Where("id = ?", id).Preload("Organizations").First(&user).Error; err != nil {
		slogger.Errorf("Error finding user %s", id)
		return nil, err
	}

	if err := c.Db.Model(&user).Updates(input).Error; err != nil {
		return nil, err
	}

	c.Logger.Info("Updated user", zap.Object("user", &user))
	return &user, nil
}

func (c *UserHandler) DeleteUser(id string) error {
	var user models.User

	if err := c.Db.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	if err := c.Db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
