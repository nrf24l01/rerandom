package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/rerandom/schemas"
)

func (h *Handler) GetRandomInteger(c echo.Context) error {
	// Создаем экземпляр структуры для запроса
	req := &schemas.IntegerRequest{}
	
	// Биндим GET параметры к структуре
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request parameters",
		})
	}
	
	// Валидируем структуру
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	tmp := schemas.TemplateAnswer{
		RandomNumber: "42",
		Min:          req.Min,
		Max:          req.Max,
		Timestamp:    "2023-01-01T00:00:00Z",
	}

	return c.Render(http.StatusOK, "random_return.html", tmp)
}