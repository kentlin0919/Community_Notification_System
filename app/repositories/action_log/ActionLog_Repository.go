package action_log

import (
	"Community_Notification_System/database"
	"Community_Notification_System/database/ActionLog_DB"
)

type ActionLogRepository struct{}

func (r *ActionLogRepository) CreateLog(logEntry *ActionLog_DB.ActionLog) error {
	result := database.DB.Create(logEntry)
	return result.Error
}
