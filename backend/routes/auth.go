package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/go-web-utils/echokit"
	"github.com/nrf24l01/rerandom/backend/handlers"
	"github.com/nrf24l01/rerandom/backend/schemas"
)

func RegisterAuthRoutes(e *echo.Echo, h *handlers.Handler) {
	g := e.Group("/auth")

	g.POST("/login", h.Login, echokit.ValidationMiddleware(func() interface{} {
		return &schemas.LoginRequest{}
	}))

	g.POST("/refresh", h.Refresh)
}
