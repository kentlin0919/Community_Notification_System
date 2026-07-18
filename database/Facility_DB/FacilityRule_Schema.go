package facility_db

import (
	"time"
)

// FacilityRule 設施預約規則（防呆、時間配置等）
type FacilityRule struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FacilityID        uint64    `gorm:"uniqueIndex;not null" json:"facility_id"`        // 對應 FacilityInfo.ID
	SlotMinutes       int       `gorm:"not null;default:60" json:"slot_minutes"`        // 每個預約單位的時間 (分鐘)
	MaxAdvanceDays    int       `gorm:"not null;default:7" json:"max_advance_days"`     // 可提前預約的最大天數
	CancelBeforeHours int       `gorm:"not null;default:24" json:"cancel_before_hours"` // 最晚需提前幾小時取消
	AutoApprove       bool      `gorm:"not null;default:false" json:"auto_approve"`     // 預約是否自動審核通過
	IsActive          bool      `gorm:"not null;default:false" json:"is_active"`        // 此規則是否已正式生效 (對應前端的開放預約)
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
