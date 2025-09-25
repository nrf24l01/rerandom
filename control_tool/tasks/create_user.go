package tasks

import (
	"errors"

	"github.com/nrf24l01/rerandom/backend/models"
)

func (h *Handler) CreateUser(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}

	var existingUser models.User
	if err := h.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return errors.New("user already exists")
	}

	user := &models.User{
		Username: username,
	}
	user.SetPassword(password)

	if err := h.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}
