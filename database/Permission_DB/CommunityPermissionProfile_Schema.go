package permission_db

import (
	"time"
)

// CommunityPermissionProfile 社區角色顯示設定表
// 讓每個社區自訂 3 ~ 7 的顯示名稱，但不改變權限大小
type CommunityPermissionProfile struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CommunityID  uint64    `gorm:"index:idx_community_permission,unique;not null" json:"community_id"`
	PermissionID int       `gorm:"index:idx_community_permission,unique;not null;check:permission_id >= 3 AND permission_id <= 7" json:"permission_id"`
	DisplayName  string    `gorm:"type:varchar(50);not null" json:"display_name"`
	Description  string    `gorm:"type:text" json:"description"`
	UpdatedBy    string    `gorm:"type:varchar(100)" json:"updated_by"` // user_id (uuid)
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
