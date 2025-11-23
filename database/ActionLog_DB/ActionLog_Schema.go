package ActionLog_DB

import (
	"time"

	"gorm.io/gorm"
)

type ActionLog struct {
	gorm.Model
	Timestamp   time.Time `gorm:"comment:日誌記錄時間"`
	Module      string    `gorm:"type:varchar(50);comment:模組名稱"`
	APIPath     string    `gorm:"type:varchar(255);comment:API路徑"`
	ErrorCode   int       `gorm:"comment:錯誤碼或狀態碼"`
	UserID      uint      `gorm:"comment:操作者ID"`
	Description string    `gorm:"type:text;comment:操作描述"`
}

func (ActionLog) TableName() string {
	return "action_logs"
}
