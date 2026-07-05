package platform

import (
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	platform_db "Community_Notification_System/database/Platform_DB"
)

// GetPlatformVersionConfig 依 os 取得版本設定
func GetPlatformVersionConfig(os string) repositoryModels.RepositoryModel[*platform_db.PlatformVersionConfig] {
	var config platform_db.PlatformVersionConfig
	var repositoryModel repositoryModels.RepositoryModel[*platform_db.PlatformVersionConfig]

	result := database.DB.Where("os = ?", os).First(&config)

	repositoryModel.Statue = *result
	if result.Error == nil {
		repositoryModel.Result = &config
	}

	return repositoryModel
}
