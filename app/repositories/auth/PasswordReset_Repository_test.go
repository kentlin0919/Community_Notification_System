package auth

import (
	"testing"
	"time"

	"Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupPasswordResetTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}
	if err := db.AutoMigrate(&user_db.UserInfo{}, &user_db.PasswordResetToken{}); err != nil {
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

func TestCreateAndGetLatestPasswordResetTokenByEmail(t *testing.T) {
	setupPasswordResetTestDB(t)

	token := &user_db.PasswordResetToken{
		ID:        "token-1",
		Email:     "user@example.com",
		OtpHash:   "hash1",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := CreatePasswordResetToken(token); err != nil {
		t.Fatalf("建立密碼重設紀錄失敗: %v", err)
	}

	found, err := GetLatestPasswordResetTokenByEmail("user@example.com")
	if err != nil {
		t.Fatalf("查詢密碼重設紀錄失敗: %v", err)
	}
	if found == nil || found.ID != "token-1" {
		t.Fatalf("預期找到 token-1，實際為 %+v", found)
	}
}

func TestGetLatestPasswordResetTokenByEmailReturnsNilWhenNotFound(t *testing.T) {
	setupPasswordResetTestDB(t)

	found, err := GetLatestPasswordResetTokenByEmail("nouser@example.com")
	if err != nil {
		t.Fatalf("預期無錯誤，實際為: %v", err)
	}
	if found != nil {
		t.Fatalf("預期找不到紀錄，實際為 %+v", found)
	}
}

func TestGetUserByEmail(t *testing.T) {
	setupPasswordResetTestDB(t)

	user := user_db.UserInfo{ID: "user-1", Email: "user@example.com", Password: "hashed"}
	if err := database.DB.Create(&user).Error; err != nil {
		t.Fatalf("建立測試使用者失敗: %v", err)
	}

	result := GetUserByEmail("user@example.com")
	if result.Statue.Error != nil {
		t.Fatalf("查詢使用者失敗: %v", result.Statue.Error)
	}
	if result.Result == nil || result.Result.ID != "user-1" {
		t.Fatalf("預期找到 user-1，實際為 %+v", result.Result)
	}
}
