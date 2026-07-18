package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware 生成並注入 X-Request-ID 到 Context 與 Response Header 中
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Request Header 嘗試獲取 X-Request-ID (用於微服務之間的傳遞)
		requestID := c.GetHeader("X-Request-ID")

		// 如果沒有，則生成一組新的 UUID 作為 Request ID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 寫入 Context 供後續的日誌或邏輯使用
		c.Set("X-Request-ID", requestID)

		// 寫入 Response Header 讓客戶端能夠追蹤
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
