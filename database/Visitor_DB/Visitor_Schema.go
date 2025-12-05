package visitor_db

import (
	"time"
)

type VisitorInfo struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ResidentID  string    `gorm:"type:varchar(64);not null" json:"resident_id"`
	VisitorName string    `gorm:"type:varchar(100);not null" json:"visitor_name"`
	VisitDate   time.Time `gorm:"type:date" json:"visit_date"`
	QRCode      string    `gorm:"type:text" json:"qrcode"` // Base64 or URL
	ValidUntil  time.Time `json:"valid_until"`
	Status      string    `gorm:"type:varchar(20)" json:"status"` // active / expired
	CreatedAt   time.Time `json:"created_at"`
}
