package auth

import (
	"testing"
	"time"

	repository "Community_Notification_System/app/repositories/auth"
	"Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupVerifyOTPTestDB(t *testing.T) {
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

func TestVerifyOTPSucceedsAndIssuesResetToken(t *testing.T) {
	setupVerifyOTPTestDB(t)
	otpHash := hashToken("123456")
	repository.CreatePasswordResetToken(&user_db.PasswordResetToken{
		ID: "token-1", Email: "user@example.com", OtpHash: otpHash,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	})

	resetToken, err := VerifyOTP("user@example.com", "123456")
	if err != nil {
		t.Fatalf("預期無錯誤，實際為: %v", err)
	}
	if resetToken == "" {
		t.Fatal("預期回傳非空的 ResetToken")
	}
}

func TestVerifyOTPFailsWithWrongCode(t *testing.T) {
	setupVerifyOTPTestDB(t)
	otpHash := hashToken("123456")
	repository.CreatePasswordResetToken(&user_db.PasswordResetToken{
		ID: "token-1", Email: "user@example.com", OtpHash: otpHash,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	})

	_, err := VerifyOTP("user@example.com", "000000")
	if err != ErrOtpInvalid {
		t.Fatalf("預期 ErrOtpInvalid，實際為: %v", err)
	}
}

func TestVerifyOTPFailsWhenExpired(t *testing.T) {
	setupVerifyOTPTestDB(t)
	otpHash := hashToken("123456")
	repository.CreatePasswordResetToken(&user_db.PasswordResetToken{
		ID: "token-1", Email: "user@example.com", OtpHash: otpHash,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	})

	_, err := VerifyOTP("user@example.com", "123456")
	if err != ErrOtpInvalid {
		t.Fatalf("預期 ErrOtpInvalid，實際為: %v", err)
	}
}

func TestVerifyOTPLocksAfterFiveAttempts(t *testing.T) {
	setupVerifyOTPTestDB(t)
	otpHash := hashToken("123456")
	repository.CreatePasswordResetToken(&user_db.PasswordResetToken{
		ID: "token-1", Email: "user@example.com", OtpHash: otpHash,
		ExpiresAt: time.Now().Add(10 * time.Minute), AttemptCount: 5,
	})

	_, err := VerifyOTP("user@example.com", "123456")
	if err != ErrOtpLocked {
		t.Fatalf("預期 ErrOtpLocked，實際為: %v", err)
	}
}
