package user_db

import (
	"Community_Notification_System/pkg/common"

	"gorm.io/gorm"
)

type UserTablesController struct{}

func NewUserDBController() *UserTablesController {
	return &UserTablesController{}
}

func (u *UserTablesController) UserTable(DB *gorm.DB) {
	// 檢查是否存在 UserInfo 表

	common.NewCreateTableController().Base_Create_Table(DB, &UserInfo{}, "user_info")

}
