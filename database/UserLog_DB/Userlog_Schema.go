package userlog_db

import (
	"time"
)

type UserLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"type:varchar(150);not null"`
	Action    string    `gorm:"type:text;not null"`
	Timestamp time.Time `gorm:"type:datetime;not null"`
}
