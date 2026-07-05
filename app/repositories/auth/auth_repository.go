package auth

import (
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"

	"gorm.io/gorm"
)

// CreateSession 新增一筆 Refresh Token 紀錄
func CreateSession(session *user_db.UserSession) repositoryModels.RepositoryModel[bool] {
	var repositoryModel repositoryModels.RepositoryModel[bool]

	err := database.DB.Create(session)
	repositoryModel.Statue = *err
	repositoryModel.Result = err.Error == nil

	return repositoryModel
}

// GetSessionByToken 依據 Refresh Token 取得 Session 紀錄
func GetSessionByToken(token string) repositoryModels.RepositoryModel[*user_db.UserSession] {
	var session user_db.UserSession
	var repositoryModel repositoryModels.RepositoryModel[*user_db.UserSession]

	// 加入 FOR UPDATE 以防並發 (Race Condition) 進行鎖定
	result := database.DB.Where("refresh_token = ?", token).First(&session)

	repositoryModel.Statue = *result
	if result.Error == nil {
		repositoryModel.Result = &session
	}

	return repositoryModel
}

// UpdateSession 儲存修改的 Session 紀錄
func UpdateSession(session *user_db.UserSession) repositoryModels.RepositoryModel[bool] {
	var repositoryModel repositoryModels.RepositoryModel[bool]

	err := database.DB.Save(session)
	repositoryModel.Statue = *err
	repositoryModel.Result = err.Error == nil

	return repositoryModel
}

// RevokeAllUserSessions 強制登出該使用者的所有裝置
func RevokeAllUserSessions(tx *gorm.DB, userID string) error {
	return tx.Model(&user_db.UserSession{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

// GetUserByID 從 User_DB 取得 UserInfo
func GetUserByID(userID string) repositoryModels.RepositoryModel[*user_db.UserInfo] {
	var user user_db.UserInfo
	var repositoryModel repositoryModels.RepositoryModel[*user_db.UserInfo]

	result := database.DB.Where("id = ?", userID).First(&user)
	repositoryModel.Statue = *result
	if result.Error == nil {
		repositoryModel.Result = &user
	}

	return repositoryModel
}

// TransactionWrapper 包裝 Transaction 用於解決 Race condition 鎖定機制
func TransactionWrapper(fn func(tx *gorm.DB) error) error {
	return database.DB.Transaction(fn)
}
