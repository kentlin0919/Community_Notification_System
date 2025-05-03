package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// CookieMiddleware 處理 cookie 相關的中間件
func CookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err == nil {
			fmt.Println("Session ID from cookie:", sessionID)
			// 可加入上下文 c.Set(...) 做後續使用
		}
		c.Next()
	}
}
