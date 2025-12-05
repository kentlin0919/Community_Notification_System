package Message_Repository

import (
	message_model "Community_Notification_System/app/models/message"
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	message_db "Community_Notification_System/database/Message_DB"
	user_db "Community_Notification_System/database/User_DB"
	"time"

	"gorm.io/gorm"
)

func MessageRepository(userInfoList []*user_db.UserInfo, messageModel *message_model.MessageData) repositoryModels.RepositoryModel[bool] {
	// 初始化返回結果結構
	var repositoryModel repositoryModels.RepositoryModel[bool]
	// 用於追蹤資料庫操作的最後一個錯誤
	var lastErr *gorm.DB

	// 遍歷所有接收訊息的用戶
	for _, userInfo := range userInfoList {
		// 為每個用戶創建新的訊息記錄
		var message message_db.MessageInfo
		message.UserID = userInfo.ID          // 設置用戶ID
		message.Title = messageModel.Title    // 設置訊息標題
		message.Content = messageModel.Detail // 設置訊息詳細內容
		message.Type = "mail"                 // 預設類型
		message.Status = "unread"             // 預設狀態
		message.CreateTime = time.Now()       // 設置訊息創建時間

		// 將訊息保存到資料庫
		lastErr = database.DB.Create(&message)
	}

	// 設置返回結果的錯誤信息
	repositoryModel.Statue.Error = lastErr.Error
	// 根據是否有錯誤設置操作結果
	repositoryModel.Result = lastErr.Error == nil

	return repositoryModel
}
