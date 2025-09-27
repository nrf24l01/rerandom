package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/rerandom/backend/models"
	"github.com/nrf24l01/rerandom/backend/schemas"
)

func (h *Handler) AddPredict(c echo.Context) error {
	predict_data := c.Get("validatedBody").(*schemas.PredictAddRequest)

	var predict models.Digit
	predict.Value = predict_data.Answ
	if predict_data.Min != nil {
		predict.Min = predict_data.Min
	}
	if predict_data.Max != nil {
		predict.Max = predict_data.Max
	}
	if predict_data.DropCount != nil {
		predict.MaxDrops = *predict_data.DropCount
	}

	if err := h.DB.Create(&predict).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	return c.JSON(http.StatusOK, schemas.PredictResponse{
		UUID:        predict.ID.String(),
		Answ:        predict.Value,
		Min:         predict.Min,
		Max:         predict.Max,
		Dropped:     0,
		MaxDrops:    predict.MaxDrops,
		Finished:    false,
		Added:       predict.CreatedAt.Unix(),
		LastDropped: nil,
	})
}

func (h *Handler) PredictList(c echo.Context) error {
	var predicts []schemas.PredictResponse
	
	err := h.DB.Table("digits").
		Where("digits.deleted_at IS NULL").
		Select(`digits.id::text as uuid,
			digits.value as answ,
			CASE WHEN digits.min = 0 THEN NULL ELSE digits.min END as min,
			CASE WHEN digits.max = 0 THEN NULL ELSE digits.max END as max,
			COALESCE(COUNT(digit_drops.id), 0) as dropped,
			digits.max_drops,
			CASE WHEN COALESCE(COUNT(digit_drops.id), 0) >= digits.max_drops THEN true ELSE false END as finished,
			EXTRACT(EPOCH FROM digits.created_at)::bigint as added,
			EXTRACT(EPOCH FROM MAX(digit_drops.created_at))::bigint as last_dropped`).
		Joins("LEFT JOIN digit_drops ON digits.id = digit_drops.digit_id").
		Group("digits.id, digits.value, digits.min, digits.max, digits.max_drops, digits.created_at").
		Scan(&predicts).Error
		
	if err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	return c.JSON(http.StatusOK, predicts)
}

func (h *Handler) DeletePredict(c echo.Context) error {
	predict_uuid := c.Param("uuid")

	if _, err := uuid.Parse(predict_uuid); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid UUID"})
	}

	var predict models.Digit
	if err := h.DB.Where("id = ?", predict_uuid).First(&predict).Error; err != nil {
		return c.JSON(http.StatusNotFound, schemas.DefaultNotFoundResponse)
	}

	if err := h.DB.Where("digit_id = ?", predict.ID).Delete(&models.DigitDrop{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	if err := h.DB.Delete(&predict).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Predict deleted successfully"})
}