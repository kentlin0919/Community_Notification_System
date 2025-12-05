package facility_db

import (
	"time"
)

type Facilities struct {
	ID               uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Community_id     uint64 `gorm:"not null" json:"community_id"`
	Name             string `gorm:"type:varchar(100);not null" json:"name"`
	Description      string `gorm:"type:text" json:"description"`
	OpenTime         string `gorm:"type:time" json:"open_time"` // Using string for TIME type in GORM usually works, or time.Time
	CloseTime        string `gorm:"type:time" json:"close_time"`
	MaxDuration      int    `json:"max_duration"` // Minutes
	Quota            int    `json:"quota"`
	Require_Approval bool   `json:"require_approval"`
	Status           string `gorm:"type:varchar(20)" json:"status"` // active / closed
}

type FacilityReservations struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FacilityID uint64    `gorm:"not null" json:"facility_id"`
	ResidentID string    `gorm:"type:varchar(64);not null" json:"resident_id"`
	Date       time.Time `gorm:"type:date" json:"date"`
	Start      string    `gorm:"type:time" json:"start"`
	End        string    `gorm:"type:time" json:"end"`
	Status     string    `gorm:"type:varchar(20)" json:"status"` // pending / approved / canceled
	CreateTime time.Time `json:"create_time"`
}
