package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/rerandom/handlers"
)

func RegisterRoutes(e *echo.Echo, h *handlers.Handler) {
	RegisterIntegerRoutes(e, h)
}
