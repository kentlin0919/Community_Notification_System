package user_db

import "time"

type UserSession struct {
	ID           string    `gorm:"primaryKey"`
	UserID       string    `json:"user_id" gorm:"index;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"uniqueIndex;not null"`
	DeviceID     string    `json:"device_id"`
	DeviceInfo   string    `json:"device_info"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"index;not null"`
	IsRevoked    bool      `json:"is_revoked" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
