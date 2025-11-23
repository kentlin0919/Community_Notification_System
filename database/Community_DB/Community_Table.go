package communitydb

import (
	"Community_Notification_System/pkg/databasepkg"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type CommunityInfoController struct{}

func NewCommunityInfoController() *CommunityInfoController {
	return &CommunityInfoController{}
}

func (c *CommunityInfoController) CommunityInfoTable(DB *gorm.DB) {
	// 檢查是否存在 CommunityInfo 表
	databasepkg.NewCreateTableController().Base_Create_Table(DB, &CommunityInfo{}, "community_info")

	if err := seedDefaultCommunityInfo(DB); err != nil {
		log.Printf("初始化 CommunityInfo 預設資料失敗: %v", err)
	}
}

func seedDefaultCommunityInfo(db *gorm.DB) error {

	defaultCommunity := CommunityInfo{
		Community_id:   1,
		PostalCode:     251,
		Municipality:   "新北市",
		District:       "淡水區",
		RoadName:       "濱海路一段",
		LaneNumber:     306,
		AlleyNumber:    0,
		Community_name: "甜水郡社區",
		Address:        "251新北市淡水區濱海路一段306巷",
	}

	var existing CommunityInfo
	err := db.
		Where(map[string]interface{}{
			"postal_code":  defaultCommunity.PostalCode,
			"municipality": defaultCommunity.Municipality,
			"district":     defaultCommunity.District,
			"road_name":    defaultCommunity.RoadName,
			"lane_number":  defaultCommunity.LaneNumber,
			"alley_number": defaultCommunity.AlleyNumber, // 0 也會被包含
		}).
		First(&existing).Error
	switch {
	case err == nil:

	case errors.Is(err, gorm.ErrRecordNotFound):

		if createErr := db.Create(&defaultCommunity).Error; createErr != nil {
			return fmt.Errorf("新增預設社區 %s 失敗: %w", defaultCommunity.Community_name, createErr)
		}
		log.Printf("新增預設社區%s - %s", defaultCommunity.Community_name, defaultCommunity.Address)
	default:
		return fmt.Errorf("查詢預設社區 %s 失敗: %w", defaultCommunity.Community_name, err)
	}

	return nil
}
