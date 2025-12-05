package platform_db

type PlatformInfo struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Platform string `gorm:"not null;type:varchar(50)" json:"platform"`
}
