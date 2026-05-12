package model

import (
	"Community_Notification_System/utils/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorRequest struct {
	RequestID    string `json:"request_id,omitempty" example:"uuid-string"`
	Code         int    `json:"code" example:"400"`
	InternalCode int    `json:"internal_code,omitempty" example:"1001"`
	Status       string `json:"status" example:"Bad Request"`
	Error        string `json:"error" example:"error"`
}

// Response400Error 代表 400 錯誤回應
type Response400Error struct {
	RequestID    string `json:"request_id,omitempty" example:"uuid-string"`
	Code         int    `json:"code" example:"400"`
	InternalCode int    `json:"internal_code,omitempty" example:"1001"`
	Status       string `json:"status" example:"Bad Request"`
	Error        string `json:"error" example:"Invalid input data"`
}

// Response401Error 代表 401 錯誤回應
type Response401Error struct {
	RequestID    string `json:"request_id,omitempty" example:"uuid-string"`
	Code         int    `json:"code" example:"401"`
	InternalCode int    `json:"internal_code,omitempty" example:"1002"`
	Status       string `json:"status" example:"Unauthorized"`
	Error        string `json:"error" example:"Invalid or missing JWT token"`
}

// Response403Error 代表 403 錯誤回應
type Response403Error struct {
	RequestID    string `json:"request_id,omitempty" example:"uuid-string"`
	Code         int    `json:"code" example:"403"`
	InternalCode int    `json:"internal_code,omitempty" example:"1003"`
	Status       string `json:"status" example:"Forbidden"`
	Error        string `json:"error" example:"Permission denied"`
}

// Response404Error 代表 404 錯誤回應
type Response404Error struct {
	RequestID    string `json:"request_id,omitempty" example:"uuid-string"`
	Code         int    `json:"code" example:"404"`
	InternalCode int    `json:"internal_code,omitempty" example:"1004"`
	Status       string `json:"status" example:"Not Found"`
	Error        string `json:"error" example:"Resource not found"`
}

// Response500Error 代表 500 錯誤回應
type Response500Error struct {
	RequestID    string `json:"request_id,omitempty" example:"uuid-string"`
	Code         int    `json:"code" example:"500"`
	InternalCode int    `json:"internal_code,omitempty" example:"1000"`
	Status       string `json:"status" example:"Internal Server Error"`
	Error        string `json:"error" example:"Internal server error"`
}

func getRequestID(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, exists := ctx.Get("X-Request-ID"); exists {
		if reqIDStr, ok := reqID.(string); ok {
			return reqIDStr
		}
	}
	return ""
}

// NewErrorRequest 舊版建構子 (建議棄用，改用 NewErrorResponse)
func NewErrorRequest(code int, message string) ErrorRequest {
	return ErrorRequest{
		Code:   code,
		Status: http.StatusText(code),
		Error:  message,
	}
}

// NewErrorResponse 建立統一的錯誤回應
func NewErrorResponse(ctx *gin.Context, httpCode int, internalCode int, message string) ErrorRequest {
	return ErrorRequest{
		RequestID:    getRequestID(ctx),
		Code:         httpCode,
		InternalCode: internalCode,
		Status:       http.StatusText(httpCode),
		Error:        message,
	}
}

// NewGlobalErrorRequest 建立一個帶有自定義 Global Error Code 的錯誤回應
func NewGlobalErrorRequest(httpCode int, internalCode int) ErrorRequest {
	return ErrorRequest{
		Code:         httpCode,
		InternalCode: internalCode,
		Status:       http.StatusText(httpCode),
		Error:        errors.GetMsg(internalCode),
	}
}

// NewGlobalErrorRequestWithMsg 允許自定義錯誤訊息並帶上 Global Error Code
func NewGlobalErrorRequestWithMsg(httpCode int, internalCode int, message string) ErrorRequest {
	return ErrorRequest{
		Code:         httpCode,
		InternalCode: internalCode,
		Status:       http.StatusText(httpCode),
		Error:        message,
	}
}

type RequestMessage struct {
	Message string `json:"Message" example:"message"`
}
