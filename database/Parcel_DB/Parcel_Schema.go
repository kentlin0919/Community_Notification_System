package parcel_db

import (
	"time"
)

// ParcelInfo 包裹代收紀錄表
type ParcelInfo struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CommunityID    uint64     `gorm:"index;not null" json:"community_id"`
	HomeID         uint64     `gorm:"index;not null" json:"home_id"`
	CourierCompany string     `gorm:"type:varchar(100);not null" json:"courier_company"`
	TrackingNumber string     `gorm:"type:varchar(100);not null" json:"tracking_number"`
	Status         int        `gorm:"not null;default:1" json:"status"` // 1=待領取, 2=已領取
	Remark         string     `gorm:"type:varchar(255)" json:"remark"`
	ReceivedAt     time.Time  `gorm:"autoCreateTime" json:"received_at"`
	PickedUpAt     *time.Time `json:"picked_up_at"`
	CreatedBy      string     `gorm:"type:varchar(100)" json:"created_by"` // 登錄者 user_id
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
