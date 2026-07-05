package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	repository "Community_Notification_System/app/repositories/auth"
	user_db "Community_Notification_System/database/User_DB"
	"Community_Notification_System/pkg/email"

	"github.com/google/uuid"
)

// ForgotPassword 產生 OTP 並寄送，帳號是否存在皆回傳成功以避免帳號枚舉
func ForgotPassword(userEmail string, sender email.EmailSender) error {
	userResult := repository.GetUserByEmail(userEmail)
	if userResult.Statue.Error != nil || userResult.Result == nil {
		// 帳號不存在：不寄信、不建立紀錄，但對外仍視為成功
		return nil
	}

	otp, err := generateNumericOTP(6)
	if err != nil {
		return err
	}
	otpHash := hashToken(otp)

	token := &user_db.PasswordResetToken{
		ID:        uuid.New().String(),
		Email:     userEmail,
		OtpHash:   otpHash,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := repository.CreatePasswordResetToken(token); err != nil {
		return err
	}

	return sender.SendOTP(userEmail, otp)
}

func generateNumericOTP(length int) (string, error) {
	digits := make([]byte, length)
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i, b := range buf {
		digits[i] = '0' + (b % 10)
	}
	return string(digits), nil
}

func hashToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
