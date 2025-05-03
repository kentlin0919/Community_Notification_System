package userlog_db

import (
	"time"

	"gorm.io/gorm"
)

type UserLog struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"not null"`
	Action    string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
}
