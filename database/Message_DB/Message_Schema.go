package message_db

import "time"

type MessageInfo struct {
	ID           string    `gorm:"primaryKey;autoIncrement"`
	UserID       string    `json:"name"`
	Email        string    `json:"email"`
	SenderID     string    `json:"sender_id"`
	CommunityID  uint64    `json:"community_id"`
	BatchID      string    `json:"batch_id"`
	TargetType   string    `json:"target_type"`
	Title        string    `json:"Title"`
	Subtile      string    `json:"Subtile" example:"Subtile"`
	Detail       string    `json:"Detail" example:"Detail"`
	Category     string    `json:"category"`
	EntityType   string    `json:"entity_type"`
	EntityID     string    `json:"entity_id"`
	FcmStatus    string    `json:"fcm_status"`
	FcmMessageID string    `json:"fcm_message_id"`
	FcmError     string    `json:"fcm_error"`
	IsRead       bool      `json:"is_read"`
	CreateTime   time.Time `json:"CreateTime" example:"2025-03-23T15:04:05Z"`
}
