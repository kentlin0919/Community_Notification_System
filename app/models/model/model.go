package model

import "net/http"

type ErrorRequest struct {
	Code   int    `json:"code" example:"400"`
	Status string `json:"status" example:"Bad Request"`
	Error  string `json:"error" example:"error"`
}

func NewErrorRequest(code int, message string) ErrorRequest {
	return ErrorRequest{
		Code:   code,
		Status: http.StatusText(code),
		Error:  message,
	}
}

type RequestMessage struct {
	Message string `json:"Message" example:"message"`
}
