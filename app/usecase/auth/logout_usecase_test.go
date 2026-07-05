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

func setupLogoutTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}
	if err := db.AutoMigrate(&user_db.UserSession{}); err != nil {
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

func TestLogoutRevokesSession(t *testing.T) {
	setupLogoutTestDB(t)
	database.DB.Create(&user_db.UserSession{ID: "session-1", UserID: "user-1", RefreshToken: "rt-1", ExpiresAt: time.Now().Add(time.Hour)})

	if err := Logout("rt-1"); err != nil {
		t.Fatalf("預期無錯誤，實際為: %v", err)
	}

	var session user_db.UserSession
	database.DB.Where("id = ?", "session-1").First(&session)
	if !session.IsRevoked {
		t.Fatal("預期登出後 session 被撤銷")
	}
}

func TestLogoutIsNoopForUnknownToken(t *testing.T) {
	setupLogoutTestDB(t)

	if err := Logout("unknown-token"); err != nil {
		t.Fatalf("預期無錯誤，實際為: %v", err)
	}
}
