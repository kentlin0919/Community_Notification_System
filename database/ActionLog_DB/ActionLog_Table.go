package ActionLog_DB

import (
	"Community_Notification_System/pkg/databasepkg"

	"gorm.io/gorm"
)

type ActionLogController struct{}

func NewActionLogController() *ActionLogController {
	return &ActionLogController{}
}

func (c *ActionLogController) ActionLogTable(DB *gorm.DB) {
	databasepkg.NewCreateTableController().Base_Create_Table(DB, &ActionLog{}, "action_logs")
}
