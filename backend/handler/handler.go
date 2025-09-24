package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/nrf24l01/rerandom/backend/core"
	"github.com/nrf24l01/rerandom/backend/redis"
)

type Handler struct {
	DB     *gorm.DB
	Config *core.Config
	Redis  *redis.RedisClient
}

type IntegersRequest struct {
	Num    int    `form:"num" binding:"required,min=1,max=10000"`
	Min    int    `form:"min" binding:"required"`
	Max    int    `form:"max" binding:"required"`
	Col    int    `form:"col,default=1" binding:"min=1"`
	Base   int    `form:"base,default=10" binding:"oneof=2 8 10 16"`
	Format string `form:"format,default=plain" binding:"oneof=plain html json"`
	Rnd    string `form:"rnd,default=new"`
	Cl     string `form:"cl,default=w" binding:"oneof=w b"`
}

type SetRequest struct {
	Num         int    `json:"num" form:"num" binding:"required"`
	Key         string `json:"key,omitempty" form:"key,omitempty"`
	Description string `json:"description,omitempty" form:"description,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SetResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Value   int  `json:"value"`
}

type IntegersJSONResponse struct {
	Numbers    []int                  `json:"numbers"`
	Parameters map[string]interface{} `json:"parameters"`
}

// GenerateIntegers обрабатывает GET /integers/
func (h *Handler) GenerateIntegers(c *gin.Context) {
	var req IntegersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_parameter",
			Message: "Invalid parameter value",
			Details: err.Error(),
		})
		return
	}

	// Проверяем, что min <= max
	if req.Min > req.Max {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_parameter",
			Message: "Min value cannot be greater than max value",
		})
		return
	}

	var numbers []int
	
	// Получаем числа из Redis или генерируем случайные
	for i := 0; i < req.Num; i++ {
		var num int
		
		// Пытаемся получить из Redis
		if redisValue, err := h.Redis.GetFirstFromSet(); err == nil {
			if parsedValue, parseErr := strconv.Atoi(redisValue); parseErr == nil {
				// Проверяем, что значение в диапазоне
				if parsedValue >= req.Min && parsedValue <= req.Max {
					num = parsedValue
				} else {
					// Если не в диапазоне, генерируем случайное
					num = rand.Intn(req.Max-req.Min+1) + req.Min
				}
			} else {
				// Если не удалось распарсить, генерируем случайное
				num = rand.Intn(req.Max-req.Min+1) + req.Min
			}
		} else {
			// Если Redis пуст, генерируем случайное число
			num = rand.Intn(req.Max-req.Min+1) + req.Min
		}
		
		numbers = append(numbers, num)
	}

	// Форматируем ответ в зависимости от запрошенного формата
	switch req.Format {
	case "json":
		c.JSON(http.StatusOK, IntegersJSONResponse{
			Numbers: numbers,
			Parameters: map[string]interface{}{
				"num":  req.Num,
				"min":  req.Min,
				"max":  req.Max,
				"base": req.Base,
			},
		})
	case "html":
		// Форматируем число в соответствии с системой счисления
		numStr := formatNumberInBase(numbers[0], req.Base)
		now := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")
		htmlContent := fmt.Sprintf(
			"<center><span style='font-size:100%%;font-weight:bold;'>%s</span><span style='font-size:70%%;'>Min:&nbsp;%d, Max:&nbsp;%d<br>%s</span></center>",
			numStr, req.Min, req.Max, now,
		)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	default: // plain
		numStr := formatNumberInBase(numbers[0], req.Base)
		c.String(http.StatusOK, numStr)
	}
}

// SetValue обрабатывает POST /set
func (h *Handler) SetValue(c *gin.Context) {
	var req SetRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_parameter",
			Message: "Invalid parameter value",
			Details: err.Error(),
		})
		return
	}

	// Сохраняем значение в Redis
	if err := h.Redis.AddToSet(strconv.Itoa(req.Num)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "redis_error",
			Message: "Failed to save value to Redis",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SetResponse{
		Success: true,
		Message: "Value set successfully",
		Value:   req.Num,
	})
}

// formatNumberInBase форматирует число в заданной системе счисления
func formatNumberInBase(num, base int) string {
	switch base {
	case 2:
		return fmt.Sprintf("%b", num)
	case 8:
		return fmt.Sprintf("%o", num)
	case 16:
		return fmt.Sprintf("%x", num)
	default:
		return fmt.Sprintf("%d", num)
	}
}