package user_db

import (
	"time"
)

type UserInfo struct {
	ID           string    `gorm:"primaryKey;type:varchar(64)"`
	PermissionId int       `json:"PermissionId"`
	Name         string    `json:"name" gorm:"type:varchar(100)"`
	Email        string    `json:"email" gorm:"type:varchar(150);unique"`
	Password     string    `json:"password" gorm:"type:varchar(255)"`
	Ssotoken     string    `json:"Ssotoken"`
	Fcmtoken     string    `json:"Fcmtoken"`
	Home_id      uint64    `json:"Home_id"`
	Registertime time.Time `json:"Registertime" example:"2025-03-23T15:04:05Z"`
	Birthdaytime time.Time `json:"Birthdaytime" example:"2025-03-23T15:04:05Z"`
	Token        string    `json:"Token"`
	Platform     int       `json:"Platform"`
	Session_id   string    `json:"Session_id"`
	Community_id uint64    `json:"community_id"`
}
