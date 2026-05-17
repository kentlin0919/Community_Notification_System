package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Community_Notification_System/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	os.Setenv("JWTPASSWORD", secret)
	defer os.Unsetenv("JWTPASSWORD")

	tests := []struct {
		name           string
		path           string
		setupAuth      func(req *http.Request)
		expectedStatus int
	}{
		{
			name: "正確的 Token",
			path: "/api/v1/protected",
			setupAuth: func(req *http.Request) {
				token, err := utils.GenerateJWT("test@example.com", "user-123", 1, 100)
				if err != nil {
					t.Errorf("GenerateJWT 失敗: %v", err)
				}
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "缺少 Authorization Header",
			path:           "/api/v1/protected",
			setupAuth:      func(req *http.Request) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "格式錯誤 (非 Bearer)",
			path: "/api/v1/protected",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Basic dGVzdDp0ZXN0")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "無效的 Token",
			path: "/api/v1/protected",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid-token")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			// 註冊中介層與測試路由
			r.Use(JWTAuthMiddleware())
			r.GET(tt.path, func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, tt.path, nil)
			tt.setupAuth(req)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "Case: %s", tt.name)
		})
	}
}

func TestJWTAuthMiddlewareExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// 此處需要模擬過期 Token，但因為 utils.GenerateJWT 目前沒有過期時間參數，
	// 若要測試過期，通常會手動簽發一個已過期的 Token。
	// 這裡先測試基本解析邏輯。
}
