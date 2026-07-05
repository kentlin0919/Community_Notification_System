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

type stubEmailSender struct {
	sentTo  string
	sentOTP string
}

func (s *stubEmailSender) SendOTP(toEmail string, otp string) error {
	s.sentTo = toEmail
	s.sentOTP = otp
	return nil
}

func setupForgotPasswordTestDB(t *testing.T) {
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

func TestForgotPasswordSendsOTPForExistingUser(t *testing.T) {
	setupForgotPasswordTestDB(t)
	database.DB.Create(&user_db.UserInfo{ID: "user-1", Email: "user@example.com"})

	sender := &stubEmailSender{}
	if err := ForgotPassword("user@example.com", sender); err != nil {
		t.Fatalf("預期無錯誤，實際為: %v", err)
	}
	if sender.sentTo != "user@example.com" {
		t.Fatalf("預期寄送對象為 user@example.com，實際為 %s", sender.sentTo)
	}
	if len(sender.sentOTP) != 6 {
		t.Fatalf("預期 OTP 長度為 6，實際為 %d", len(sender.sentOTP))
	}

	var stored user_db.PasswordResetToken
	if err := database.DB.Where("email = ?", "user@example.com").First(&stored).Error; err != nil {
		t.Fatalf("預期儲存密碼重設紀錄，實際查詢失敗: %v", err)
	}
	if stored.ExpiresAt.Before(time.Now()) {
		t.Fatal("預期 ExpiresAt 在未來")
	}
}

func TestForgotPasswordDoesNotErrorForUnknownEmail(t *testing.T) {
	setupForgotPasswordTestDB(t)

	sender := &stubEmailSender{}
	if err := ForgotPassword("nouser@example.com", sender); err != nil {
		t.Fatalf("為避免帳號枚舉，預期無錯誤，實際為: %v", err)
	}
	if sender.sentTo != "" {
		t.Fatal("預期未知帳號不會實際寄送 Email")
	}
}
