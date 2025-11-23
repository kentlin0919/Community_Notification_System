package platform_db

import (
	"Community_Notification_System/pkg/databasepkg"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type PlatformInfoController struct{}

func NewPlatformInfoController() *PlatformInfoController {
	return &PlatformInfoController{}
}

func (u *PlatformInfoController) PlatformInfoTable(DB *gorm.DB) {
	// 檢查是否存在 PlatformInfo 表
	databasepkg.NewCreateTableController().Base_Create_Table(DB, &PlatformInfo{}, "platform_info")
	if err := seedDefaultPlatformInfo(DB); err != nil {
		log.Printf("初始化 platform_info 預設資料失敗: %v", err)
	}
}

func seedDefaultPlatformInfo(db *gorm.DB) error {
	defaultPlatformInfo := []PlatformInfo{
		{Platform: "web"},
		{Platform: "App"},
		{Platform: "Desktop"},
	}

	for _, perm := range defaultPlatformInfo {
		var existing PlatformInfo
		err := db.Where("Platform = ?", perm.Platform).First(&existing).Error
		switch {
		case err == nil:
			continue
		case errors.Is(err, gorm.ErrRecordNotFound):
			if createErr := db.Create(&perm).Error; createErr != nil {
				return fmt.Errorf("新增預設權限 %s 失敗: %w", perm.Platform, createErr)
			}
			log.Printf("新增預設權限：%s - %s", perm.ID, perm.Platform)
		default:
			return fmt.Errorf("查詢權限 %s 失敗: %w", perm.Platform, err)
		}
	}
	return nil
}
