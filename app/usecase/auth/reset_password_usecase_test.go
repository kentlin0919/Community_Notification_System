package auth

import (
	"testing"
	"time"

	repository "Community_Notification_System/app/repositories/auth"
	"Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupResetPasswordTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("建立測試資料庫失敗: %v", err)
	}
	if err := db.AutoMigrate(&user_db.UserInfo{}, &user_db.PasswordResetToken{}, &user_db.UserSession{}); err != nil {
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

func TestResetPasswordUpdatesPasswordAndRevokesSessions(t *testing.T) {
	setupResetPasswordTestDB(t)

	user := user_db.UserInfo{ID: "user-1", Email: "user@example.com", Password: "old-hash"}
	database.DB.Create(&user)
	database.DB.Create(&user_db.UserSession{ID: "session-1", UserID: "user-1", RefreshToken: "rt-1", ExpiresAt: time.Now().Add(time.Hour)})

	resetTokenHash := hashToken("plain-reset-token")
	repository.CreatePasswordResetToken(&user_db.PasswordResetToken{
		ID: "token-1", Email: "user@example.com", Verified: true,
		ResetTokenHash: resetTokenHash, ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	if err := ResetPassword("plain-reset-token", "newpassword123"); err != nil {
		t.Fatalf("預期無錯誤，實際為: %v", err)
	}

	var updatedUser user_db.UserInfo
	database.DB.Where("id = ?", "user-1").First(&updatedUser)
	if bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpassword123")) != nil {
		t.Fatal("預期密碼已更新為 newpassword123")
	}

	var session user_db.UserSession
	database.DB.Where("id = ?", "session-1").First(&session)
	if !session.IsRevoked {
		t.Fatal("預期重設密碼後既有 session 被撤銷")
	}

	var usedToken user_db.PasswordResetToken
	database.DB.Where("id = ?", "token-1").First(&usedToken)
	if !usedToken.Used {
		t.Fatal("預期 ResetToken 標記為已使用")
	}
}

func TestResetPasswordFailsWithInvalidToken(t *testing.T) {
	setupResetPasswordTestDB(t)

	err := ResetPassword("bad-token", "newpassword123")
	if err != ErrResetTokenInvalid {
		t.Fatalf("預期 ErrResetTokenInvalid，實際為: %v", err)
	}
}
