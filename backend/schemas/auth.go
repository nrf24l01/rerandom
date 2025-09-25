package schemas

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
