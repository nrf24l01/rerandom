package handlers

import (
	"github.com/nrf24l01/rerandom/backend/core"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Config *core.Config
}
