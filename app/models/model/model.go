package model

import "net/http"

type ErrorRequest struct {
	Code   int    `json:"code" example:"400"`
	Status string `json:"status" example:"Bad Request"`
	Error  string `json:"error" example:"error"`
}

// Response400Error 代表 400 錯誤回應
type Response400Error struct {
	Code   int    `json:"code" example:"400"`
	Status string `json:"status" example:"Bad Request"`
	Error  string `json:"error" example:"Invalid input data"`
}

// Response401Error 代表 401 錯誤回應
type Response401Error struct {
	Code   int    `json:"code" example:"401"`
	Status string `json:"status" example:"Unauthorized"`
	Error  string `json:"error" example:"Invalid or missing JWT token"`
}

// Response403Error 代表 403 錯誤回應
type Response403Error struct {
	Code   int    `json:"code" example:"403"`
	Status string `json:"status" example:"Forbidden"`
	Error  string `json:"error" example:"Permission denied"`
}

// Response404Error 代表 404 錯誤回應
type Response404Error struct {
	Code   int    `json:"code" example:"404"`
	Status string `json:"status" example:"Not Found"`
	Error  string `json:"error" example:"Resource not found"`
}

// Response500Error 代表 500 錯誤回應
type Response500Error struct {
	Code   int    `json:"code" example:"500"`
	Status string `json:"status" example:"Internal Server Error"`
	Error  string `json:"error" example:"Internal server error"`
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
