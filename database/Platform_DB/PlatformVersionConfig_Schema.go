package platform_db

// PlatformVersionConfig 記錄各平台的最低與最新版本設定
type PlatformVersionConfig struct {
	ID             int    `gorm:"primaryKey;autoIncrement"`
	OS             string `gorm:"uniqueIndex;not null" json:"os"`
	MinVersion     string `gorm:"not null" json:"min_version"`
	LatestVersion  string `gorm:"not null" json:"latest_version"`
	ForceUpdateMsg string `json:"force_update_msg"`
}
