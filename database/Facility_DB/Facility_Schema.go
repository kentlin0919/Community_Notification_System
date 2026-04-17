package facility_db

import (
	"time"
)

// FacilityInfo 可供住戶預約的公共設施主檔
type FacilityInfo struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CommunityID  uint64    `gorm:"index;not null" json:"community_id"`
	Name         string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_community_facility_name" json:"name"`
	FacilityType string    `gorm:"type:varchar(50);not null" json:"facility_type"` // e.g., gym, pool, meeting_room
	Location     string    `gorm:"type:varchar(255)" json:"location"`
	Description  string    `gorm:"type:text" json:"description"`
	Status       string    `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` // draft, active, inactive, maintenance
	CoverImage   string    `gorm:"type:varchar(255)" json:"cover_image"`
	CreatedBy    string    `gorm:"type:varchar(100)" json:"created_by"` // user_id UUID
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
