package middlewares

import (
	"net/http"
	"strings"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

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

		if skipPaths[c.Request.URL.Path] || skipPaths[c.FullPath()] {
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
			return []byte(os.Getenv("JWTPASSWORD")), nil
		})

		if err != nil || !token.Valid {
			// fmt.Printf("Token validation failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 token", "details": err.Error()})
			c.Abort()
			return
		}

		// 使用 claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("username", claims["username"])
			c.Set("user_id", claims["user_id"])

			// 支援從 float64/int 轉換
			if pid, ok := claims["permission_id"].(float64); ok {
				c.Set("permission_id", int(pid))
			} else {
				c.Set("permission_id", claims["permission_id"])
			}

			if cid, ok := claims["community_id"].(float64); ok {
				c.Set("community_id", uint64(cid))
			} else {
				c.Set("community_id", claims["community_id"])
			}
		}

	}
}
