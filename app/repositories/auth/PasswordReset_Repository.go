package auth

import (
	"errors"

	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"

	"gorm.io/gorm"
)

// CreatePasswordResetToken 新增一筆密碼重設紀錄
func CreatePasswordResetToken(token *user_db.PasswordResetToken) error {
	return database.DB.Create(token).Error
}

// GetLatestPasswordResetTokenByEmail 依 email 取得最新一筆尚未使用的密碼重設紀錄
func GetLatestPasswordResetTokenByEmail(email string) (*user_db.PasswordResetToken, error) {
	var token user_db.PasswordResetToken
	err := database.DB.Where("email = ? AND used = ?", email, false).
		Order("created_at DESC").
		First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// GetPasswordResetTokenByResetHash 依 ResetToken hash 取得紀錄
func GetPasswordResetTokenByResetHash(hash string) (*user_db.PasswordResetToken, error) {
	var token user_db.PasswordResetToken
	err := database.DB.Where("reset_token_hash = ? AND used = ?", hash, false).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// SavePasswordResetToken 更新既有的密碼重設紀錄
func SavePasswordResetToken(token *user_db.PasswordResetToken) error {
	return database.DB.Save(token).Error
}

// GetUserByEmail 依 email 取得使用者資訊
func GetUserByEmail(email string) repositoryModels.RepositoryModel[*user_db.UserInfo] {
	var user user_db.UserInfo
	var repositoryModel repositoryModels.RepositoryModel[*user_db.UserInfo]

	result := database.DB.Where("email = ?", email).First(&user)

	repositoryModel.Statue = *result
	if result.Error == nil {
		repositoryModel.Result = &user
	}

	return repositoryModel
}
