package tasks

import (
	"github.com/nrf24l01/rerandom/control_tool/core"
	"gorm.io/gorm"
)

type Handler struct {
	DB      *gorm.DB
	Config  *core.Config
}