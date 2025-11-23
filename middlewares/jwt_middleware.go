package middlewares

import (
	"Community_Notification_System/app/models/model"
	"Community_Notification_System/app/repositories/api_route"
	"Community_Notification_System/utils"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

		if skipPaths[c.FullPath()] {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 Authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "格式錯誤，應為 Bearer token"})
			c.Abort()
			return
		}

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 token"})
			c.Abort()
			return
		}

		// 新的權限檢查邏輯
		apiRouteRepo := api_route.NewApiRouteRepository()
		requiredPermID, err := apiRouteRepo.GetRequiredPermission(c.FullPath(), c.Request.Method)

		// 如果路由未在資料庫中定義，為安全起見，預設拒絕存取
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorModel := model.NewErrorRequest(http.StatusForbidden, "此 API 路徑未定義權限")
			c.JSON(http.StatusForbidden, errorModel)
			c.Abort()
			return
		}

		// 假設 0 代表公開路由，無需權限
		if requiredPermID != 0 {
			if claims.PermissionID != requiredPermID {
				errorModel := model.NewErrorRequest(http.StatusForbidden, "沒有權限存取此功能")
				c.JSON(http.StatusForbidden, errorModel)
				c.Abort()
				return
			}
		}

		c.Set("claims", claims)
		c.Next()
	}
}
