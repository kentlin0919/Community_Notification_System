package user_db

import "time"

type UserInfo struct {
	ID           string    `gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Ssotoken     string    `json:"Ssotoken"`
	Home_id      string    `json:"Home_id"`
	Registertime time.Time `json:"Registertime" example:"2025-03-23T15:04:05Z"`
	Birthdaytime time.Time `json:"Birthdaytime" example:"2025-03-23T15:04:05Z"`
	Token        string    `json:"Token"`
	Platform     string    `json:"Platform"`
	Permission   int       `json:"Permission"`
	Session_id   string    `json:"Session_id"`
}
