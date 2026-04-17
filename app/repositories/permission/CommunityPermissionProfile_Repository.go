package permission

import (
	repositoryModels "Community_Notification_System/app/models/repository"
	permissiondb "Community_Notification_System/database/Permission_DB"
	"Community_Notification_System/database"
	"gorm.io/gorm/clause"
)

// UpsertPermissionProfileRepository 新增或更新某社區的自訂權限名稱
func UpsertPermissionProfileRepository(profile *permissiondb.CommunityPermissionProfile) repositoryModels.RepositoryModel[permissiondb.CommunityPermissionProfile] {
	var result repositoryModels.RepositoryModel[permissiondb.CommunityPermissionProfile]

	// 透過 Upsert (OnConflict) 來處理如果已經存在相同的 community_id + permission_id 就做更新
	dbResult := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "community_id"}, {Name: "permission_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name", "description", "updated_by", "updated_at"}),
	}).Create(profile)

	result.Statue.Error = dbResult.Error
	result.Statue.RowsAffected = dbResult.RowsAffected
	result.Result = *profile
	return result
}

// GetCommunityPermissionProfilesRepository 取得特定社區所有的自訂權限列表
func GetCommunityPermissionProfilesRepository(communityID uint64) repositoryModels.RepositoryModel[[]permissiondb.CommunityPermissionProfile] {
	var result repositoryModels.RepositoryModel[[]permissiondb.CommunityPermissionProfile]
	var profiles []permissiondb.CommunityPermissionProfile

	dbResult := database.DB.Where("community_id = ?", communityID).Order("permission_id asc").Find(&profiles)
	
	result.Statue.Error = dbResult.Error
	result.Statue.RowsAffected = dbResult.RowsAffected
	result.Result = profiles
	return result
}
