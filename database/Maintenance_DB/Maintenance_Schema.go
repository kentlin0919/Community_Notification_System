package maintenance_db

import (
	"time"
)

// AssetTag 設備/設施標籤 (資產管理)
type AssetTag struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CommunityID uint      `gorm:"not null;index" json:"community_id"`
	TagCode     string    `gorm:"type:varchar(50);uniqueIndex" json:"tag_code"` // NFC 或 QR Code ID
	AssetName   string    `gorm:"type:varchar(100);not null" json:"asset_name"`
	Location    string    `gorm:"type:text" json:"location"`
	AssetType   string    `gorm:"type:varchar(50)" json:"asset_type"` // Elevator, FireExtinguisher, Pump
	LastCheck   time.Time `json:"last_check"`
	NextCheck   time.Time `json:"next_check"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MaintenanceTicket 維修單 (物業閉環)
type MaintenanceTicket struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CommunityID uint       `gorm:"not null;index" json:"community_id"`
	ReporterID  uint       `gorm:"not null;index" json:"reporter_id"` // 報修人
	AssetTagID  *uint      `json:"asset_tag_id"`                      // 關聯設備 (可選)
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Priority    int        `gorm:"default:0" json:"priority"`                        // 0: Normal, 1: Urgent
	Status      string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, assigned, in_progress, completed, closed
	AssigneeID  *uint      `json:"assignee_id"`                                      // 維修人員
	PhotoBefore string     `gorm:"type:text" json:"photo_before"`
	PhotoAfter  string     `gorm:"type:text" json:"photo_after"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
