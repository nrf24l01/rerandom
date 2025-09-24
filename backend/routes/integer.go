package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/rerandom/handlers"
)

func RegisterIntegerRoutes(e *echo.Echo, h *handlers.Handler) {
	group := e.Group("/integers")

	group.GET("/", h.GetRandomInteger)
}