package community

import (
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	communitydb "Community_Notification_System/database/Community_DB"
	userdb "Community_Notification_System/database/User_DB"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CreateApplicationRepository 新增社區申請單
func CreateApplicationRepository(application *communitydb.CommunityRegisterApplication) repositoryModels.RepositoryModel[communitydb.CommunityRegisterApplication] {
	var result repositoryModels.RepositoryModel[communitydb.CommunityRegisterApplication]

	createResult := database.DB.Create(application)
	result.Statue = *createResult
	if createResult.Error != nil {
		return result
	}

	result.Result = *application
	return result
}

// GetApplicationByIDRepository 依據 ID 取得申請單
func GetApplicationByIDRepository(id uint64) repositoryModels.RepositoryModel[communitydb.CommunityRegisterApplication] {
	var result repositoryModels.RepositoryModel[communitydb.CommunityRegisterApplication]
	var application communitydb.CommunityRegisterApplication

	findResult := database.DB.First(&application, id)
	result.Statue = *findResult
	if findResult.Error != nil {
		return result
	}

	result.Result = application
	return result
}

// UpdateApplicationStatusRepository 更新申請單狀態 (Reject用)
func UpdateApplicationStatusRepository(id uint64, status string, by string, reason string) repositoryModels.RepositoryModel[communitydb.CommunityRegisterApplication] {
	var result repositoryModels.RepositoryModel[communitydb.CommunityRegisterApplication]
	
	now := time.Now()
	updateResult := database.DB.Model(&communitydb.CommunityRegisterApplication{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"reviewed_by":   by,
		"reviewed_at":   now,
		"reject_reason": reason,
	})
	
	result.Statue = *updateResult
	if updateResult.Error != nil {
		return result
	}
	
	return GetApplicationByIDRepository(id)
}

// ApproveApplicationTransactionRepository 使用 Transaction 核可申請單
func ApproveApplicationTransactionRepository(application *communitydb.CommunityRegisterApplication, communityInfo *communitydb.CommunityInfo, userInfo *userdb.UserInfo, by string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 建立 community_info
		if err := tx.Create(communityInfo).Error; err != nil {
			return errors.New("建立社區資料失敗")
		}

		// 為了讓 admin 使用者綁定正確的 community_id
		userInfo.Community_id = communityInfo.Community_id

		// 2. 建立 admin 使用者
		if err := tx.Create(userInfo).Error; err != nil {
			return errors.New("建立管理員帳號失敗")
		}

		// 3. 更新申請單狀態
		now := time.Now()
		if err := tx.Model(&communitydb.CommunityRegisterApplication{}).Where("id = ?", application.ID).Updates(map[string]interface{}{
			"status":      "approved",
			"reviewed_by": by,
			"reviewed_at": now,
		}).Error; err != nil {
			return errors.New("更新申請單狀態失敗")
		}

		return nil
	})
}
