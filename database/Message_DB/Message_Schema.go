package message_db

import "time"

type MessageInfo struct {
	ID         string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	UserID     string    `gorm:"type:varchar(64);not null" json:"user_id"`
	Title      string    `gorm:"type:varchar(255)" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	Type       string    `gorm:"type:varchar(20)" json:"type"` // mail / reminder / delivery
	ImageURL   string    `gorm:"type:text" json:"image_url"`
	Status     string    `gorm:"type:varchar(10)" json:"status"` // unread / read
	CreateTime time.Time `json:"create_time"`
}
