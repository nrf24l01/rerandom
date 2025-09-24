package handlers

import (
	"gorm.io/gorm"

	"github.com/nrf24l01/rerandom/core"
	"github.com/nrf24l01/rerandom/redis"
)

type Handler struct {
	DB     *gorm.DB
	Config *core.Config
	Redis  *redis.RedisClient
}
