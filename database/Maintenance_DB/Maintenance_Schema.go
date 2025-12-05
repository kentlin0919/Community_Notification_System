package maintenance_db

import (
	"time"
)

type MaintenanceReports struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ResidentID  string    `gorm:"type:varchar(64);not null" json:"resident_id"`
	Category    string    `gorm:"type:varchar(100)" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	Photos      string    `gorm:"type:json" json:"photos"`        // Storing JSON as string or using specific JSON type if supported by DB/GORM
	Status      string    `gorm:"type:varchar(20)" json:"status"` // pending / processing / completed
	AssignedTo  string    `gorm:"type:varchar(64)" json:"assigned_to"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}
