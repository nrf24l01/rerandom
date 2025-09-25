package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/go-web-utils/jwtutil"
	"github.com/nrf24l01/rerandom/backend/models"
	"github.com/nrf24l01/rerandom/backend/schemas"
)


func (h *Handler) Login(c echo.Context) error {
	user_data := c.Get("validatedBody").(*schemas.LoginRequest)

	var user models.User
	if err := h.DB.Where("username = ?", user_data.Username).First(&user).Error; err != nil {
		return echo.ErrUnauthorized
	}

	if !user.CheckPassword(user_data.Password) {
		return echo.ErrUnauthorized
	}

	accessToken, refreshToken, err := jwtutil.GenerateTokenPair(user.ID.String(), user.Username, []byte(h.Config.JWTAccessSecret), []byte(h.Config.JWTRefreshSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = refreshToken
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, schemas.LoginResponse{
		AccessToken: accessToken,
	})
}

func (h *Handler) Refresh(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return echo.ErrUnauthorized
	}

	claims, err := jwtutil.ValidateToken(refreshToken.Value, []byte(h.Config.JWTRefreshSecret))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, schemas.Message{Status: "Unauthorized"})
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, schemas.Message{Status: "Unauthorized"})
	}
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return c.JSON(http.StatusUnauthorized, schemas.Message{Status: "Unauthorized"})
	}

	accessToken, newRefreshToken, err := jwtutil.GenerateTokenPair(userID, username, []byte(h.Config.JWTAccessSecret), []byte(h.Config.JWTRefreshSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = newRefreshToken
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, schemas.RefreshResponse{
		AccessToken: accessToken,
	})
}