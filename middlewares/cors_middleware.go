package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("🔥 CORS middleware triggered for:", c.Request.URL.Path)

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 前端執行的 domain
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // 這一行很關鍵！

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
