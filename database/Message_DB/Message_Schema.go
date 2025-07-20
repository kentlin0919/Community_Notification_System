package message_db

import "time"

type MessageInfo struct {
	ID         string    `gorm:"primaryKey;autoIncrement"`
	UserID     string    `json:"name"`
	Email      string    `json:"email"`
	Title      string    `json:"Title"`
	Subtile    string    `json:"Subtile" example:"Subtile"`
	Detail     string    `json:"Detail" example:"Detail"`
	CreateTime time.Time `json:"CreateTime" example:"2025-03-23T15:04:05Z"`
}
