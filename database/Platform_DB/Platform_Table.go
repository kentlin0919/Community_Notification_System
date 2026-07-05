package platform_db

import (
	"Community_Notification_System/pkg/common"
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
	common.NewCreateTableController().Base_Create_Table(DB, &PlatformInfo{}, "platform_info")
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
			log.Printf("新增預設權限：%d - %s", perm.ID, perm.Platform)
		default:
			return fmt.Errorf("查詢權限 %s 失敗: %w", perm.Platform, err)
		}
	}
	return nil
}

func (u *PlatformInfoController) PlatformVersionConfigTable(DB *gorm.DB) {
	common.NewCreateTableController().Base_Create_Table(DB, &PlatformVersionConfig{}, "platform_version_config")
	if err := seedDefaultPlatformVersionConfig(DB); err != nil {
		log.Printf("初始化 platform_version_config 預設資料失敗: %v", err)
	}
}

func seedDefaultPlatformVersionConfig(db *gorm.DB) error {
	defaults := []PlatformVersionConfig{
		{OS: "ios", MinVersion: "1.0.0", LatestVersion: "1.0.0", ForceUpdateMsg: "請更新至最新版本以獲得最佳體驗"},
		{OS: "android", MinVersion: "1.0.0", LatestVersion: "1.0.0", ForceUpdateMsg: "請更新至最新版本以獲得最佳體驗"},
	}
	for _, cfg := range defaults {
		var existing PlatformVersionConfig
		err := db.Where("os = ?", cfg.OS).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if createErr := db.Create(&cfg).Error; createErr != nil {
				return fmt.Errorf("新增預設版本設定 %s 失敗: %w", cfg.OS, createErr)
			}
			log.Printf("新增預設版本設定：%s", cfg.OS)
		}
	}
	return nil
}
