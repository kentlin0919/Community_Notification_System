package iot_db

import (
	"time"
)

// IoTDevice 物聯網設備主檔
type IoTDevice struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CommunityID uint      `gorm:"not null;index" json:"community_id"`
	DeviceName  string    `gorm:"type:varchar(100);not null" json:"device_name"`
	DeviceType  string    `gorm:"type:varchar(50);not null" json:"device_type"` // e.g., Gate, Elevator, Intercom
	MACAddress  string    `gorm:"type:varchar(17);uniqueIndex" json:"mac_address"`
	IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address"`
	Status      string    `gorm:"type:varchar(20);default:'online'" json:"status"` // online, offline, maintenance
	LastOnline  time.Time `json:"last_online"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// VisitorPass 訪客通行證 (動態碼)
type VisitorPass struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CommunityID uint      `gorm:"not null;index" json:"community_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"` // 申請人
	VisitorName string    `gorm:"type:varchar(50)" json:"visitor_name"`
	PassCode    string    `gorm:"type:varchar(100);uniqueIndex" json:"pass_code"` // 加密或雜湊後的通行碼
	StartTime   time.Time `gorm:"not null" json:"start_time"`
	EndTime     time.Time `gorm:"not null" json:"end_time"`
	UseCount    int       `gorm:"default:0" json:"use_count"`
	MaxUse      int       `gorm:"default:1" json:"max_use"`
	Status      string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active, expired, used_up
	CreatedAt   time.Time `json:"created_at"`
}

// AccessLog 通行紀錄
type AccessLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CommunityID uint      `gorm:"not null;index" json:"community_id"`
	DeviceID    uint      `gorm:"not null;index" json:"device_id"`
	UserID      *uint     `json:"user_id"`         // 可為空 (若為訪客通行證)
	PassID      *uint     `json:"pass_id"`         // 可為空 (若為住戶感應)
	AccessMethod string   `gorm:"type:varchar(20)" json:"access_method"` // qrcode, card, app, remote
	AccessTime  time.Time `gorm:"index" json:"access_time"`
	Result      string    `gorm:"type:varchar(20)" json:"result"` // success, denied
}
