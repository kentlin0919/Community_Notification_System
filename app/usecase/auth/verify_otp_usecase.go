package auth

import (
	"errors"
	"time"

	repository "Community_Notification_System/app/repositories/auth"
	"Community_Notification_System/utils"
)

var (
	ErrOtpInvalid = errors.New("otp invalid or expired")
	ErrOtpLocked  = errors.New("otp locked due to too many attempts")
)

const maxOtpAttempts = 5

// VerifyOTP 驗證 OTP，成功後核發短效期 ResetToken 並回傳明文（僅此次回傳，DB 只存 hash）
func VerifyOTP(userEmail string, otp string) (string, error) {
	token, err := repository.GetLatestPasswordResetTokenByEmail(userEmail)
	if err != nil {
		return "", err
	}
	if token == nil || token.Verified || time.Now().After(token.ExpiresAt) {
		return "", ErrOtpInvalid
	}
	if token.AttemptCount >= maxOtpAttempts {
		return "", ErrOtpLocked
	}

	if hashToken(otp) != token.OtpHash {
		token.AttemptCount++
		_ = repository.SavePasswordResetToken(token)
		return "", ErrOtpInvalid
	}

	resetToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	token.Verified = true
	token.ResetTokenHash = hashToken(resetToken)
	token.ExpiresAt = time.Now().Add(5 * time.Minute)
	if err := repository.SavePasswordResetToken(token); err != nil {
		return "", err
	}

	return resetToken, nil
}
