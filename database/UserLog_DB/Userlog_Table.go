package userlog_db

import (
	// "Community_Notification_System/pkg/common"

	"gorm.io/gorm"
)

type UserLogTableController struct{}

func NewUserLogTableController() *UserLogTableController {
	return &UserLogTableController{}
}

func (u *UserLogTableController) UserLogTable(DB *gorm.DB) {
	DB.AutoMigrate(&UserLog{})
}
