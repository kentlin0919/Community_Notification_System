package middlewares

import (
	"net/http"
	"strings"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var JwtKey = []byte(os.Getenv("JWTPASSWORD"))

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//需要跳過的路由
		skipPaths := map[string]bool{

			"/":                        true,
			"/api/v1/login":            true,
			"/api/v1/platform/getlist": true,
			"/api/v1/register":         true,
			"/swagger/*any":            true,
		}

		if skipPaths[c.FullPath()] {
			c.Next()
			return
		}

		//從 map 中撈出對應的路徑確定是否跳過
		if skipPaths[c.Request.RequestURI] {
			c.Next()
			return
		}

		//從key中撈出value
		authHeader := c.GetHeader("Authorization")

		//確定是否要包含
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少 Authorization header",
			})
			c.Abort()
			return
		}

		//撈出token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "格式錯誤，應為 Bearer token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 token"})
			c.Abort()
			return
		}

		// 使用 claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("username", claims["username"])
		}

	}
}
