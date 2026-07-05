package platform

import (
	"testing"

	"Community_Notification_System/database"
	platform_db "Community_Notification_System/database/Platform_DB"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupPlatformVersionTestDB(t *testing.T) {
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

func TestGetPlatformVersionConfig(t *testing.T) {
	setupPlatformVersionTestDB(t)
	database.DB.Create(&platform_db.PlatformVersionConfig{OS: "ios", MinVersion: "1.1.0", LatestVersion: "1.2.0"})

	result := GetPlatformVersionConfig("ios")
	if result.Statue.Error != nil {
		t.Fatalf("查詢版本設定失敗: %v", result.Statue.Error)
	}
	if result.Result == nil || result.Result.MinVersion != "1.1.0" {
		t.Fatalf("預期 MinVersion 為 1.1.0，實際為 %+v", result.Result)
	}
}
