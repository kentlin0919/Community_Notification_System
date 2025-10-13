package community

import (
	communityModel "Community_Notification_System/app/models/community"
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	communitydb "Community_Notification_System/database/Community_DB"
	"fmt"

	"gorm.io/gorm"
)

// CommunityListResult 封裝社區列表查詢結果
type CommunityListResult struct {
	Total int64
	Items []communitydb.CommunityInfo
}

func CommunityOneRepository(communityInfo communitydb.CommunityInfo) repositoryModels.RepositoryModel[CommunityListResult] {

	baseQuery := database.DB.Model(&communitydb.CommunityInfo{})

	applyFilters := func(db *gorm.DB) *gorm.DB {

		if communityInfo.Municipality != "" {
			db = db.Where("municipality = ?", communityInfo.Municipality)
		}
		if communityInfo.District != "" {
			db = db.Where("district = ?", communityInfo.District)
		}
		if communityInfo.Community_name != "" {
			keyword := fmt.Sprintf("%%%s%%", communityInfo.Community_name)
			db = db.Where("(community_name ILIKE ? OR address ILIKE ?)", keyword, keyword)
		}
		return db
	}

	filteredForCount := applyFilters(baseQuery.Session(&gorm.Session{}))
	var result repositoryModels.RepositoryModel[CommunityListResult]
	countResult := filteredForCount.Count(&result.Result.Total)
	if countResult.Error != nil {
		return result
	}

	listResult := filteredForCount.Order("community_id asc").Find(&result.Result.Items)
	result.Statue = *listResult

	return result

}

// CommunityListRepository 依據查詢條件取得社區列表
func CommunityListRepository(query *communityModel.CommunityListQuery) repositoryModels.RepositoryModel[CommunityListResult] {
	baseQuery := database.DB.Model(&communitydb.CommunityInfo{})

	applyFilters := func(db *gorm.DB) *gorm.DB {
		if query == nil {
			return db
		}
		if query.Municipality != "" {
			db = db.Where("municipality = ?", query.Municipality)
		}
		if query.District != "" {
			db = db.Where("district = ?", query.District)
		}
		if query.PostalCode != nil {
			db = db.Where("postal_code = ?", *query.PostalCode)
		}
		if query.Keyword != "" {
			keyword := fmt.Sprintf("%%%s%%", query.Keyword)
			db = db.Where("(community_name ILIKE ? OR address ILIKE ?)", keyword, keyword)
		}
		return db
	}

	filteredForCount := applyFilters(baseQuery.Session(&gorm.Session{}))

	var result repositoryModels.RepositoryModel[CommunityListResult]
	countResult := filteredForCount.Count(&result.Result.Total)
	result.Statue = *countResult
	if countResult.Error != nil {
		return result
	}

	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	filteredForList := applyFilters(baseQuery.Session(&gorm.Session{}))
	listResult := filteredForList.Order("community_id asc").Offset(offset).Limit(pageSize).Find(&result.Result.Items)
	result.Statue = *listResult

	return result
}

// RegisterRepository 新增社區基本資料
func RegisterRepository(community_info *communitydb.CommunityInfo) repositoryModels.RepositoryModel[communitydb.CommunityInfo] {
	var result repositoryModels.RepositoryModel[communitydb.CommunityInfo]

	createResult := database.DB.Create(community_info)
	result.Statue = *createResult
	if createResult.Error != nil {
		return result
	}

	result.Result = *community_info
	return result
}
