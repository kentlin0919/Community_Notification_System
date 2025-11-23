package ActionLog_DB

import (
	"Community_Notification_System/pkg/common"
	"gorm.io/gorm"
)

type ActionLogController struct{}

func NewActionLogController() *ActionLogController {
	return &ActionLogController{}
}

func (c *ActionLogController) ActionLogTable(DB *gorm.DB) {
	common.NewCreateTableController().Base_Create_Table(DB, &ActionLog{}, "action_logs")
}
