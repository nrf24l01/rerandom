package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/go-web-utils/echokit"
	"github.com/nrf24l01/rerandom/handlers"
	"github.com/nrf24l01/rerandom/schemas"
)

func RegisterPredictRoutes(e *echo.Echo, h *handlers.Handler) {
	group := e.Group("/predict")
	
	// JWT middleware для всех маршрутов предиктов
	group.Use(echokit.JWTMiddleware([]byte(h.Config.JWTAccessSecret)))

	// Добавить предикт
	group.POST("/add", h.AddPredict, echokit.ValidationMiddleware(func() interface{} {
		return &schemas.PredictAddRequest{}
	}))
	group.GET("/list", h.PredictList)
}
