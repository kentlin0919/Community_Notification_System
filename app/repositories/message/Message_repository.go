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

type MessageListResult struct {
	Total int64
	Items []message_db.MessageInfo
}

func MessageRepository(userInfoList []*user_db.UserInfo, messageModel *message_model.MessageData) repositoryModels.RepositoryModel[bool] {
	// 初始化返回結果結構
	var repositoryModel repositoryModels.RepositoryModel[bool]
	// 用於追蹤資料庫操作的最後一個錯誤
	var lastErr *gorm.DB

	// 遍歷所有接收訊息的用戶
	for _, userInfo := range userInfoList {
		// 為每個用戶創建新的訊息記錄
		var message message_db.MessageInfo
		message.UserID = userInfo.ID           // 設置用戶ID
		message.Email = userInfo.Email         // 設置用戶郵箱
		message.Detail = messageModel.Detail   // 設置訊息詳細內容
		message.Subtile = messageModel.Subtile // 設置訊息標題
		message.Title = messageModel.Title
		message.IsRead = false
		message.CreateTime = time.Now()

		// 將訊息保存到資料庫
		lastErr = database.DB.Create(&message)
	}

	if lastErr == nil {
		repositoryModel.Result = true
		return repositoryModel
	}

	// 設置返回結果的錯誤信息
	repositoryModel.Statue.Error = lastErr.Error
	// 根據是否有錯誤設置操作結果
	repositoryModel.Result = lastErr.Error == nil

	return repositoryModel
}

func GetUserMessagesRepository(userID string, query *message_model.MessageListQuery) repositoryModels.RepositoryModel[MessageListResult] {
	var result repositoryModels.RepositoryModel[MessageListResult]

	base := database.DB.Model(&message_db.MessageInfo{}).Where("user_id = ?", userID)
	if query != nil && query.IsRead != nil {
		base = base.Where("is_read = ?", *query.IsRead)
	}

	countResult := base.Count(&result.Result.Total)
	result.Statue = *countResult
	if countResult.Error != nil {
		return result
	}

	page := 1
	pageSize := 20
	if query != nil {
		if query.Page > 0 {
			page = query.Page
		}
		if query.PageSize > 0 {
			pageSize = query.PageSize
		}
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	listResult := base.Order("create_time desc").Offset(offset).Limit(pageSize).Find(&result.Result.Items)
	result.Statue = *listResult
	return result
}

func MarkMessageReadRepository(userID, messageID string) repositoryModels.RepositoryModel[bool] {
	var result repositoryModels.RepositoryModel[bool]
	updateResult := database.DB.Model(&message_db.MessageInfo{}).Where("id = ? AND user_id = ?", messageID, userID).Update("is_read", true)
	result.Statue = *updateResult
	result.Result = updateResult.Error == nil && updateResult.RowsAffected > 0
	return result
}

func MarkAllMessagesReadRepository(userID string) repositoryModels.RepositoryModel[int64] {
	var result repositoryModels.RepositoryModel[int64]
	updateResult := database.DB.Model(&message_db.MessageInfo{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true)
	result.Statue = *updateResult
	result.Result = updateResult.RowsAffected
	return result
}
