package parcel

import (
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	parcel_db "Community_Notification_System/database/Parcel_DB"
	user_db "Community_Notification_System/database/User_DB"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CreateParcelRepository 建立包裹代收紀錄，並檢查同社區是否有相同單號的待領包裹
func CreateParcelRepository(parcel *parcel_db.ParcelInfo) repositoryModels.RepositoryModel[parcel_db.ParcelInfo] {
	var result repositoryModels.RepositoryModel[parcel_db.ParcelInfo]

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 檢查同社區是否已有相同單號且處於「待領取」狀態
		var exist parcel_db.ParcelInfo
		if err := tx.Where("community_id = ? AND tracking_number = ? AND status = 1", parcel.CommunityID, parcel.TrackingNumber).First(&exist).Error; err == nil {
			return errors.New("該社區已存在相同單號的待領包裹")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 建立包裹紀錄
		if err := tx.Create(parcel).Error; err != nil {
			return errors.New("建立包裹紀錄失敗: " + err.Error())
		}

		return nil
	})

	if err != nil {
		result.Statue.Error = err
	} else {
		result.Statue.RowsAffected = 1
		result.Result = *parcel
	}

	return result
}

// GetParcelListRepository 查詢包裹列表，支援 homeID 與 status 篩選
func GetParcelListRepository(communityID uint64, homeID uint64, status int) repositoryModels.RepositoryModel[[]parcel_db.ParcelInfo] {
	var result repositoryModels.RepositoryModel[[]parcel_db.ParcelInfo]
	var parcels []parcel_db.ParcelInfo

	query := database.DB.Where("community_id = ?", communityID)

	if homeID > 0 {
		query = query.Where("home_id = ?", homeID)
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	res := query.Order("received_at DESC").Find(&parcels)
	result.Statue = *res
	result.Result = parcels

	return result
}

// PickupParcelRepository 標記包裹為已領取
func PickupParcelRepository(communityID uint64, parcelID uint64) repositoryModels.RepositoryModel[parcel_db.ParcelInfo] {
	var result repositoryModels.RepositoryModel[parcel_db.ParcelInfo]

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var parcel parcel_db.ParcelInfo
		if err := tx.Where("id = ? AND community_id = ?", parcelID, communityID).First(&parcel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("找不到該包裹紀錄")
			}
			return err
		}

		if parcel.Status == 2 {
			return errors.New("該包裹已被領取")
		}

		now := time.Now()
		parcel.Status = 2
		parcel.PickedUpAt = &now

		if err := tx.Save(&parcel).Error; err != nil {
			return errors.New("更新包裹狀態失敗: " + err.Error())
		}

		result.Result = parcel
		return nil
	})

	if err != nil {
		result.Statue.Error = err
	} else {
		result.Statue.RowsAffected = 1
	}

	return result
}

// GetFcmTokensByHomeID 查詢指定住戶所有成員的 FCM Token
func GetFcmTokensByHomeID(homeID uint64) []string {
	var tokens []string
	var users []user_db.UserInfo

	database.DB.Where("home_id = ? AND fcmtoken != ''", homeID).Find(&users)

	for _, u := range users {
		if u.Fcmtoken != "" {
			tokens = append(tokens, u.Fcmtoken)
		}
	}

	return tokens
}
