package middlewares

import (
	"context"
	"net/http"
	"strings"

	"Community_Notification_System/app/models/model"
	"Community_Notification_System/pkg/firebase"
	"Community_Notification_System/utils/errors"

	"github.com/gin-gonic/gin"
)

// FirebaseAuthMiddleware 驗證 Firebase JWT Token
func FirebaseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Header 取得 Authorization 內容
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			errorModel := model.NewErrorResponse(c, http.StatusUnauthorized, errors.ErrUnauthorized, "缺少 Authorization header")
			c.JSON(http.StatusUnauthorized, errorModel)
			c.Abort()
			return
		}

		// 解析 Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			errorModel := model.NewErrorResponse(c, http.StatusUnauthorized, errors.ErrTokenInvalid, "格式錯誤，應為 Bearer token")
			c.JSON(http.StatusUnauthorized, errorModel)
			c.Abort()
			return
		}

		// 確認 Firebase AuthClient 已經初始化
		if firebase.AuthClient == nil {
			errorModel := model.NewErrorResponse(c, http.StatusInternalServerError, errors.ErrInternal, "Firebase 認證服務未初始化")
			c.JSON(http.StatusInternalServerError, errorModel)
			c.Abort()
			return
		}

		// 呼叫 Firebase Admin SDK 驗證 ID Token
		token, err := firebase.AuthClient.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			errorModel := model.NewErrorResponse(c, http.StatusUnauthorized, errors.ErrTokenInvalid, "Firebase token 驗證失敗")
			c.JSON(http.StatusUnauthorized, errorModel)
			c.Abort()
			return
		}

		// 將驗證後的資訊放入 gin.Context 中供後續 controller 使用
		c.Set("firebase_uid", token.UID)
		c.Set("firebase_token", token)

		c.Next()
	}
}
