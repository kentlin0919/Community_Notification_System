package permission_db

type PermissionInfo struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	PermissionID string `gorm:"uniqueIndex;not null;type:varchar(100)" json:"permission_id"`
	Name         string `gorm:"not null;type:varchar(100)" json:"name"`
}
