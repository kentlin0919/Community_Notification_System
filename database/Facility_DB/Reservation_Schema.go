package facility_db

import (
	"time"

	"gorm.io/gorm"
)

// FacilityReservation 住戶預約主單
type FacilityReservation struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	FacilityID      uint64         `gorm:"index;not null" json:"facility_id"`
	UserID          string         `gorm:"index;type:varchar(100);not null" json:"user_id"`
	CommunityID     uint64         `gorm:"index;not null" json:"community_id"`
	HomeID          uint64         `gorm:"index" json:"home_id"`
	ReservationDate string         `gorm:"type:date;not null" json:"reservation_date"` // YYYY-MM-DD
	StartTime       string         `gorm:"type:varchar(5);not null" json:"start_time"` // HH:mm
	EndTime         string         `gorm:"type:varchar(5);not null" json:"end_time"`   // HH:mm
	PeopleCount     int            `gorm:"not null" json:"people_count"`
	Status          string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, approved, rejected, cancelled, completed, no_show
	Remark          string         `gorm:"type:text" json:"remark"`
	CancelReason    string         `gorm:"type:text" json:"cancel_reason"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// RescheduleRequest 預約改期申請單
type RescheduleRequest struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ReservationID uint64    `gorm:"index;not null" json:"reservation_id"`
	NewDate       string    `gorm:"type:date;not null" json:"new_date"`
	NewStartTime  string    `gorm:"type:varchar(5);not null" json:"new_start_time"`
	NewEndTime    string    `gorm:"type:varchar(5);not null" json:"new_end_time"`
	Status        string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, approved, rejected
	Reason        string    `gorm:"type:text" json:"reason"`
	AdminComment  string    `gorm:"type:text" json:"admin_comment"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
