package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	accountModel "Community_Notification_System/app/models/account"
	model "Community_Notification_System/app/models/model"
	"Community_Notification_System/database"
	userlog_db "Community_Notification_System/database/UserLog_DB"
	user_db "Community_Notification_System/database/User_DB"
	"Community_Notification_System/utils"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestJWT(t *testing.T) {
	t.Helper()

	previousKey := make([]byte, len(utils.JwtKey))
	copy(previousKey, utils.JwtKey)
	previousPassword, hasPrevious := os.LookupEnv("JWTPASSWORD")

	if err := os.Setenv("JWTPASSWORD", "testsecret"); err != nil {
		t.Fatalf("設定 JWTPASSWORD 環境變數失敗: %v", err)
	}
	utils.JwtKey = []byte("testsecret")

	t.Cleanup(func() {
		utils.JwtKey = previousKey
		if hasPrevious {
			if err := os.Setenv("JWTPASSWORD", previousPassword); err != nil {
				t.Fatalf("恢復 JWTPASSWORD 失敗: %v", err)
			}
		} else {
			if err := os.Unsetenv("JWTPASSWORD"); err != nil {
				t.Fatalf("清除 JWTPASSWORD 失敗: %v", err)
			}
		}
	})
}

func setupTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}

	if err := db.AutoMigrate(&user_db.UserInfo{}, &userlog_db.UserLog{}); err != nil {
		t.Fatalf("自動遷移資料表失敗: %v", err)
	}

	database.DB = db

	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		database.DB = nil
	})
}

func TestUserLoginSuccess(t *testing.T) {
	setupTestJWT(t)
	setupTestDB(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("雜湊密碼失敗: %v", err)
	}

	userRecord := user_db.UserInfo{
		Email:    "user@example.com",
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&userRecord).Error; err != nil {
		t.Fatalf("建立測試使用者失敗: %v", err)
	}

	controller := NewUserController()
	body := `{"email":"user@example.com","password":"password123","platform":"App"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	controller.UserLogin(ctx)

	if w.Code != http.StatusOK {
		t.Fatalf("預期回傳狀態碼 %d，實際為 %d", http.StatusOK, w.Code)
	}

	var response accountModel.UserRequest
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析回應 JSON 失敗: %v", err)
	}

	if response.Message != "登入成功" {
		t.Fatalf("預期成功訊息 '登入成功'，實際為 '%s'", response.Message)
	}

	if response.Token == "" {
		t.Fatal("預期 JWT Token 不為空")
	}

	var hasSessionCookie bool
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == "session_id" {
			hasSessionCookie = true
			if cookie.Value == "" {
				t.Fatal("Session cookie 未設置值")
			}
		}
	}

	if !hasSessionCookie {
		t.Fatal("預期設定 session_id cookie")
	}

	var logs []userlog_db.UserLog
	if err := database.DB.Find(&logs).Error; err != nil {
		t.Fatalf("查詢使用者紀錄失敗: %v", err)
	}

	if len(logs) != 1 {
		t.Fatalf("預期 1 筆登入紀錄，實際為 %d", len(logs))
	}

	if logs[0].Email != userRecord.Email {
		t.Fatalf("預期登入紀錄 email 為 %s，實際為 %s", userRecord.Email, logs[0].Email)
	}
}

func TestUserLoginInvalidInput(t *testing.T) {
	setupTestJWT(t)
	setupTestDB(t)

	controller := NewUserController()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	controller.UserLogin(ctx)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("預期回傳狀態碼 %d，實際為 %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]model.ErrorRequest
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析錯誤回應失敗: %v", err)
	}

	errorResponse, ok := response["error"]
	if !ok {
		t.Fatalf("預期回應包含 error 欄位，實際為 %v", response)
	}

	if errorResponse.Error != "無效的輸入資料" {
		t.Fatalf("預期錯誤訊息 '無效的輸入資料'，實際為 '%s'", errorResponse.Error)
	}
}

func TestUserLoginWrongPassword(t *testing.T) {
	setupTestJWT(t)
	setupTestDB(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("雜湊密碼失敗: %v", err)
	}

	userRecord := user_db.UserInfo{
		Email:    "user@example.com",
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&userRecord).Error; err != nil {
		t.Fatalf("建立測試使用者失敗: %v", err)
	}

	controller := NewUserController()
	body := `{"email":"user@example.com","password":"wrong","platform":"App"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	controller.UserLogin(ctx)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("預期回傳狀態碼 %d，實際為 %d", http.StatusUnauthorized, w.Code)
	}

	var response model.ErrorRequest
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析錯誤回應失敗: %v", err)
	}

	if response.Error != "帳號或密碼錯誤" {
		t.Fatalf("預期錯誤訊息 '帳號或密碼錯誤'，實際為 '%s'", response.Error)
	}
}

func TestUserLoginUserNotFound(t *testing.T) {
	setupTestJWT(t)
	setupTestDB(t)

	controller := NewUserController()
	body := `{"email":"nouser@example.com","password":"password123","platform":"App"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	controller.UserLogin(ctx)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("預期回傳狀態碼 %d，實際為 %d", http.StatusUnauthorized, w.Code)
	}

	var response model.ErrorRequest
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("解析錯誤回應失敗: %v", err)
	}

	if response.Error != "使用者不存在" {
		t.Fatalf("預期錯誤訊息 '使用者不存在'，實際為 '%s'", response.Error)
	}
}
