package ActionLog_DB

import (
	"time"
)

type ActionLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"type:varchar(64);comment:操作者ID"`
	APIPath   string    `gorm:"type:varchar(255);comment:API路徑"`
	Module    string    `gorm:"type:varchar(100);comment:模組名稱"`
	Timestamp time.Time `gorm:"type:datetime;comment:日誌記錄時間"`
}

func (ActionLog) TableName() string {
	return "action_logs"
}
