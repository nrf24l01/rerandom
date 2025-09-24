package schemas

import "net/http"

type Message struct {
	Status string `json:"status"`
}

type Error struct {
	Error string `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var DefaultInternalErrorResponse = ErrorResponse{
	Message: "Internal Server Error",
	Code:    http.StatusInternalServerError,
}