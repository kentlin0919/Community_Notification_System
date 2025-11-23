package message_db

import (
	"Community_Notification_System/pkg/databasepkg"

	"gorm.io/gorm"
)

type MessageInfoTablesController struct{}

func NewUserDBController() *MessageInfoTablesController {
	return &MessageInfoTablesController{}
}

func (u *MessageInfoTablesController) MessageInfoTable(DB *gorm.DB) {
	// 檢查是否存在 UserInfo 表
	databasepkg.NewCreateTableController().Base_Create_Table(DB, &MessageInfo{}, "message_info")

}
