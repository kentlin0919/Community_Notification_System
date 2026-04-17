package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CommunityContextMiddleware 從上下文中確保 CommunityID 存在，並限制非平台員的跨社區操作
func CommunityContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		communityIDStr := c.Param("community_id")

		// 如果 JWT Token 中沒有 community_id，那就代表未驗證身份或非預設社區成員。
		jwtCommunityID, exists := c.Get("community_id")
		permissionID, pExists := c.Get("permission_id")

		if !exists || !pExists {
			// 對於某些需要社區隔離的 API，若沒有 JWT 提供 community_id 表示異常
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無法識別所屬社區或權限，請重新登入"})
			c.Abort()
			return
		}

		permID := permissionID.(int)

		// 權限 1 (Super admin) 擁有跨社區權限，不用檢查 Param 中的 community_id
		if permID == 1 {
			c.Next()
			return
		}

		// 若使用者想透過 URL 操作他人的社區 (例如 /api/v1/communities/:community_id/...)
		// 需驗證 URL 的 community_id 是否跟他的 jwt 中的 community_id 相等
		if communityIDStr != "" {
			jwtCommIDStr := fmt.Sprintf("%v", jwtCommunityID)
			if communityIDStr != jwtCommIDStr {
				c.JSON(http.StatusForbidden, gin.H{"error": "您無權操作或查詢其他社區的資料"})
				c.Abort()
				return
			}
		}

		// 讓後續路由繼續處理，後續 Repository 都可以統一用 c.Get("community_id") 來隔離資料
		c.Next()
	}
}
