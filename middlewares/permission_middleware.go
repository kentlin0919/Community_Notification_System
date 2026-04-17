package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MinimalPermissionMiddleware 檢查使用者的 PermissionID 是否 小於等於 (即權限大於等於) 所需的最小權限
// 例如 reqPermissionID = 2 (Admin)，則 PermissionID 1 或是 2 就可以通過。 3 則會被擋下 (403)。
func MinimalPermissionMiddleware(reqPermissionID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissionID, exists := c.Get("permission_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無法驗證權限層級"})
			c.Abort()
			return
		}

		userPermID := permissionID.(int)

		// 數值越小，權限越高
		if userPermID > reqPermissionID {
			c.JSON(http.StatusForbidden, gin.H{"error": "您的權限不足，無法執行此操作"})
			c.Abort()
			return
		}

		c.Next()
	}
}
