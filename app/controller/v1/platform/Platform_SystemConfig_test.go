package platform

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	platformModel "Community_Notification_System/app/models/platform"
	"Community_Notification_System/database"
	platform_db "Community_Notification_System/database/Platform_DB"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupSystemConfigTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}
	if err := db.AutoMigrate(&platform_db.PlatformVersionConfig{}); err != nil {
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

func TestSystemConfigReturnsForceUpdateTrueWhenBelowMinVersion(t *testing.T) {
	setupSystemConfigTestDB(t)
	database.DB.Create(&platform_db.PlatformVersionConfig{OS: "ios", MinVersion: "1.1.0", LatestVersion: "1.2.0"})

	controller := NewPlatformController()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/config?os=ios&version=1.0.0", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	controller.SystemConfig(ctx)

	if w.Code != http.StatusOK {
		t.Fatalf("預期狀態碼 200，實際為 %d", w.Code)
	}
	var resp platformModel.SystemConfigResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析回應失敗: %v", err)
	}
	if !resp.ForceUpdate {
		t.Fatal("預期 force_update 為 true")
	}
	if resp.MinVersion != "1.1.0" || resp.LatestVersion != "1.2.0" {
		t.Fatalf("版本資訊不符: %+v", resp)
	}
}

func TestSystemConfigReturnsForceUpdateFalseWhenUpToDate(t *testing.T) {
	setupSystemConfigTestDB(t)
	database.DB.Create(&platform_db.PlatformVersionConfig{OS: "ios", MinVersion: "1.1.0", LatestVersion: "1.2.0"})

	controller := NewPlatformController()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/config?os=ios&version=1.2.0", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	controller.SystemConfig(ctx)

	var resp platformModel.SystemConfigResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ForceUpdate {
		t.Fatal("預期 force_update 為 false")
	}
}
