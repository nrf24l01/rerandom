package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nrf24l01/rerandom/backend/core"
	handlers "github.com/nrf24l01/rerandom/backend/handler"
	"github.com/nrf24l01/rerandom/backend/redis"
)

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Загружаем конфигурацию
	config, err := core.BuildConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к Redis
	redisClient, err := redis.CreateRedisFromCFG(config)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Создаём хендлер
	h := &handlers.Handler{
		Config: config,
		Redis:  redisClient,
	}

	// Настраиваем Gin роутер
	if config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Добавляем CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.AllowOrigins)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Регистрируем роуты
	setupRoutes(router, h)

	// Запускаем сервер
	port := ":8080"
	if config.APPHost != "" {
		port = ":" + config.APPHost
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRoutes настраивает все роуты приложения
func setupRoutes(router *gin.Engine, h *handlers.Handler) {
	// API роуты
	api := router.Group("/")
	{
		api.GET("/integers/", h.GenerateIntegers)
		api.POST("/set", h.SetValue)
	}

	// Swagger UI (опционально)
	router.Static("/swagger", "./swagger.yaml")
}