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

// GetFacilityListRepository 依社區取得設施列表；住戶視角只回傳 active 設施。
func GetFacilityListRepository(communityID uint64, onlyActive bool) repositoryModels.RepositoryModel[[]facilitydb.FacilityInfo] {
	var result repositoryModels.RepositoryModel[[]facilitydb.FacilityInfo]
	var facilities []facilitydb.FacilityInfo

	query := database.DB.Where("community_id = ?", communityID)
	if onlyActive {
		query = query.Where("status = ?", "active")
	}

	if err := query.Order("id ASC").Find(&facilities).Error; err != nil {
		result.Statue.Error = err
		return result
	}

	result.Result = facilities
	result.Statue.RowsAffected = int64(len(facilities))
	return result
}

// GetFacilityDetailRepository 取得同社區單一設施詳情。
func GetFacilityDetailRepository(id uint64, communityID uint64) repositoryModels.RepositoryModel[facilitydb.FacilityInfo] {
	var result repositoryModels.RepositoryModel[facilitydb.FacilityInfo]
	var facility facilitydb.FacilityInfo

	if err := database.DB.Where("id = ? AND community_id = ?", id, communityID).First(&facility).Error; err != nil {
		result.Statue.Error = err
		return result
	}

	result.Result = facility
	result.Statue.RowsAffected = 1
	return result
}

// UpdateFacilityWithRuleRepository 更新設施主檔並同步更新對應預約規則。
func UpdateFacilityWithRuleRepository(facility *facilitydb.FacilityInfo, rule *facilitydb.FacilityRule) repositoryModels.RepositoryModel[facilitydb.FacilityInfo] {
	var result repositoryModels.RepositoryModel[facilitydb.FacilityInfo]

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var exist facilitydb.FacilityInfo
		err := tx.Where("community_id = ? AND name = ? AND id <> ?", facility.CommunityID, facility.Name, facility.ID).First(&exist).Error
		if err == nil {
			return errors.New("該社區已存在相同名稱的設施")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := tx.Save(facility).Error; err != nil {
			return err
		}

		if rule != nil {
			rule.FacilityID = facility.ID
			var existingRule facilitydb.FacilityRule
			err := tx.Where("facility_id = ?", facility.ID).First(&existingRule).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.Create(rule).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				rule.ID = existingRule.ID
				if err := tx.Model(&existingRule).Updates(map[string]interface{}{
					"slot_minutes":        rule.SlotMinutes,
					"max_advance_days":    rule.MaxAdvanceDays,
					"cancel_before_hours": rule.CancelBeforeHours,
					"auto_approve":        rule.AutoApprove,
					"is_active":           rule.IsActive,
				}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		result.Statue.Error = err
		return result
	}

	result.Result = *facility
	result.Statue.RowsAffected = 1
	return result
}

// DeleteFacilityRepository 軟刪除同社區設施。
func DeleteFacilityRepository(id uint64, communityID uint64) error {
	res := database.DB.Where("id = ? AND community_id = ?", id, communityID).Delete(&facilitydb.FacilityInfo{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
