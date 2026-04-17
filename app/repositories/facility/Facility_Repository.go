package facility

import (
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	facilitydb "Community_Notification_System/database/Facility_DB"
	"errors"

	"gorm.io/gorm"
)

// CreateFacilityWithRuleTransactionRepository 以 Transaction 建立設施及初始預約規則
func CreateFacilityWithRuleTransactionRepository(facility *facilitydb.FacilityInfo) repositoryModels.RepositoryModel[facilitydb.FacilityInfo] {
	var result repositoryModels.RepositoryModel[facilitydb.FacilityInfo]

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 檢查同社區下是否已有相同名稱之設施
		var exist facilitydb.FacilityInfo
		if err := tx.Where("community_id = ? AND name = ?", facility.CommunityID, facility.Name).First(&exist).Error; err == nil {
			return errors.New("該社區已存在相同名稱的設施")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 1. 建立 Facility Info
		if err := tx.Create(facility).Error; err != nil {
			return errors.New("建立設施基本檔失敗: " + err.Error())
		}

		// 2. 建立預設 Facility Rule (參照文檔預設配置)
		rule := &facilitydb.FacilityRule{
			FacilityID:        facility.ID,
			SlotMinutes:       60,
			MaxAdvanceDays:    7,
			CancelBeforeHours: 24,
			AutoApprove:       false,
			IsActive:          false, // 需等待管理員開放
		}

		if err := tx.Create(rule).Error; err != nil {
			return errors.New("建立設施預設規則失敗: " + err.Error())
		}

		return nil
	})

	if err != nil {
		result.Statue.Error = err
	} else {
		result.Statue.RowsAffected = 1
		result.Result = *facility
	}

	return result
}
