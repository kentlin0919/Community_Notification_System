package permission_db

type PermissionInfo struct {
	ID           string `gorm:"primaryKey;autoIncrement"`
	PermissionID string `gorm:"uniqueIndex;not null" json:"permission_id"`
	Name         string `gorm:"not null" json:"name"`
}
