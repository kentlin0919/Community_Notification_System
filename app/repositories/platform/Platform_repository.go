package platform

import (
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	platform_db "Community_Notification_System/database/Platform_DB"
)

type PlatformListResult struct {
	Total int64
	Items []string
}

func PlatformRepository() repositoryModels.RepositoryModel[PlatformListResult] {
	var result repositoryModels.RepositoryModel[PlatformListResult]
	result.Result.Items = []string{}
	if database.DB == nil {
		return result
	}

	var platformInfos []platform_db.PlatformInfo
	dbResult := database.DB.Model(&platform_db.PlatformInfo{}).Find(&platformInfos)
	result.Statue = *dbResult

	if dbResult.Error != nil {
		return result
	}

	result.Result.Items = make([]string, 0, len(platformInfos))
	for _, platform := range platformInfos {
		result.Result.Items = append(result.Result.Items, platform.Platform)
	}

	result.Result.Total = int64(len(platformInfos))

	return result
}
