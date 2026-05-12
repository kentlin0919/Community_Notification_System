package message_db

import (
	"Community_Notification_System/pkg/common"

	"gorm.io/gorm"
)

type MessageInfoTablesController struct{}

func NewUserDBController() *MessageInfoTablesController {
	return &MessageInfoTablesController{}
}

func (u *MessageInfoTablesController) MessageInfoTable(DB *gorm.DB) {
	// 確保 message_info 表存在，並同步新增欄位
	if err := DB.AutoMigrate(&MessageInfo{}); err != nil {
		return
	}

	common.NewCreateTableController().Base_Create_Table(DB, &MessageInfo{}, "message_info")

}
