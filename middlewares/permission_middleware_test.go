package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMinimalPermissionMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userPermID     interface{}
		reqPermID      int
		expectedStatus int
	}{
		{
			name:           "權限足夠 (1 <= 2)",
			userPermID:     1,
			reqPermID:      2,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "權限相等 (2 <= 2)",
			userPermID:     2,
			reqPermID:      2,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "權限不足 (3 > 2)",
			userPermID:     3,
			reqPermID:      2,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "缺少權限資訊",
			userPermID:     nil,
			reqPermID:      2,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			r.Use(func(c *gin.Context) {
				if tt.userPermID != nil {
					c.Set("permission_id", tt.userPermID)
				}
				c.Next()
			})
			r.Use(MinimalPermissionMiddleware(tt.reqPermID))
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "Case: %s", tt.name)
		})
	}
}
