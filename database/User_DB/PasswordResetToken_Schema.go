package user_db

import "time"

// PasswordResetToken 記錄忘記密碼流程中的 OTP 與 ResetToken 狀態
type PasswordResetToken struct {
	ID             string    `gorm:"primaryKey"`
	Email          string    `json:"email" gorm:"index;not null"`
	OtpHash        string    `json:"otp_hash" gorm:"not null"`
	ResetTokenHash string    `json:"reset_token_hash" gorm:"uniqueIndex"`
	Verified       bool      `json:"verified" gorm:"default:false"`
	Used           bool      `json:"used" gorm:"default:false"`
	AttemptCount   int       `json:"attempt_count" gorm:"default:0"`
	ExpiresAt      time.Time `json:"expires_at" gorm:"index;not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
