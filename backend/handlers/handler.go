package handlers

import (
	"gorm.io/gorm"

	"github.com/nrf24l01/rerandom/core"
)

type Handler struct {
	DB     *gorm.DB
	Config *core.Config
}
