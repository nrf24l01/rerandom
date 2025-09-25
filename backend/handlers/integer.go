package handlers

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/rerandom/backend/models"
	"github.com/nrf24l01/rerandom/backend/schemas"
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

	var value string
	
	// Один запрос для поиска подходящих записей с подсчетом использований
	type DigitWithDrops struct {
		models.Digit
		DropCount        int64 `json:"drop_count"`
		RemainingDrops   int   `json:"remaining_drops"`
		Priority         int   `json:"priority"`
	}
	
	var results []DigitWithDrops
	
	err := h.DB.Raw(`
		SELECT 
			d.*,
			COALESCE(drop_counts.drop_count, 0) as drop_count,
			(d.max_drops - COALESCE(drop_counts.drop_count, 0)) as remaining_drops,
			CASE 
				WHEN d.min = ? AND d.max = ? THEN 1
				WHEN (d.min IS NULL OR d.min = 0) AND (d.max IS NULL OR d.max = 0) AND d.value >= ? AND d.value <= ? THEN 2
				ELSE 3
			END as priority
		FROM digits d
		LEFT JOIN (
			SELECT digit_id, COUNT(*) as drop_count
			FROM digit_drops
			GROUP BY digit_id
		) drop_counts ON d.id = drop_counts.digit_id
		WHERE 
			(
				(d.min = ? AND d.max = ?) OR
				((d.min IS NULL OR d.min = 0) AND (d.max IS NULL OR d.max = 0) AND d.value >= ? AND d.value <= ?)
			)
			AND (d.max_drops - COALESCE(drop_counts.drop_count, 0)) > 0
		ORDER BY priority ASC, remaining_drops DESC
		LIMIT 1
	`, req.Min, req.Max, req.Min, req.Max, req.Min, req.Max, req.Min, req.Max).Scan(&results).Error
	
	if err == nil && len(results) > 0 {
		bestDigit := results[0]
		
		// Создаем запись об использовании
		drop := models.DigitDrop{
			DigitID: bestDigit.ID,
		}
		
		err := h.DB.Create(&drop).Error
		if err == nil {
			value = strconv.Itoa(bestDigit.Value)
		}
	}
	
	// Если ничего не нашли, генерируем случайное число
	if value == "" {
		min := req.Min
		max := req.Max
		
		if min > max {
			min, max = max, min
		}
		
		rangeSize := int64(max - min + 1)
		randomBig, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to generate random number",
			})
		}
		
		result := min + int(randomBig.Int64())
		value = strconv.Itoa(result)
	}

	tmp := schemas.TemplateAnswer{
		RandomNumber: value,
		Min:          req.Min,
		Max:          req.Max,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}

	return c.Render(http.StatusOK, "random_return.html", tmp)
}