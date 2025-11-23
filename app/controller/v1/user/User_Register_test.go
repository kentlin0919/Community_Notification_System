package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	model "Community_Notification_System/app/models/model"

	"github.com/gin-gonic/gin"
)

func TestUserRegisterPasswordTooShort(t *testing.T) {
	setupTestJWT(t)
	setupTestDB(t)

	controller := NewUserController()
	body := `{"Email":"short@example.com","name":"Tester","password":"short","bethday":"2025-03-23T15:04:05-00:00","permission":1,"platform":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	controller.UserRegister(ctx)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("預期回傳狀態碼 %d，實際為 %d", http.StatusBadRequest, w.Code)
	}

	var response model.ErrorRequest
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析錯誤回應失敗: %v", err)
	}

	if response.Error != "密碼長度至少需 8 碼" {
		t.Fatalf("預期錯誤訊息 '密碼長度至少需 8 碼'，實際為 '%s'", response.Error)
	}
}
