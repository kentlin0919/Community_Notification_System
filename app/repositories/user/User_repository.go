package repository

import (
	"Community_Notification_System/app/models/account"
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	userlog_db "Community_Notification_System/database/UserLog_DB"
	user_db "Community_Notification_System/database/User_DB"
	"log"
	"time"

	"gorm.io/gorm"
)

func LoginRepository(user *account.User) repositoryModels.RepositoryModel[user_db.UserInfo] {
	var user_info user_db.UserInfo
	var repositoryModel repositoryModels.RepositoryModel[user_db.UserInfo]

	result := database.DB.Where("email = ?", user.Email).First(&user_info)

	repositoryModel.Statue = *result

	repositoryModel.Result = user_info

	return repositoryModel
}

func RegisterRepository(user_info *user_db.UserInfo) repositoryModels.RepositoryModel[bool] {

	var registerModel repositoryModels.RepositoryModel[bool]

	err := database.DB.Create(user_info)

	registerModel.Statue = *err

	log.Print(err.Error)

	registerModel.Result = err.Error == nil

	return registerModel
}

func UserLogRepository(user_info *user_db.UserInfo) repositoryModels.RepositoryModel[bool] {

	var user_log userlog_db.UserLog

	user_log.Email = user_info.Email
	user_log.Action = "登入"
	user_log.Timestamp = time.Now()

	err := database.DB.Create(&user_log)

	var repositoryModel repositoryModels.RepositoryModel[bool]

	repositoryModel.Statue = *err

	repositoryModel.Result = err.Error == nil

	return repositoryModel
}

func UserDeleteRepository(user_info *user_db.UserInfo) repositoryModels.RepositoryModel[bool] {

	var repositoryModel repositoryModels.RepositoryModel[bool]

	err := database.DB.Delete(&user_info)

	repositoryModel.Statue = *err

	repositoryModel.Result = err.Error == nil

	return repositoryModel
}

func UserInfoListRepository(userList []string) repositoryModels.RepositoryModel[[]*user_db.UserInfo] {
	var repositoryModels repositoryModels.RepositoryModel[[]*user_db.UserInfo]

	var user_infoList []*user_db.UserInfo

	var lastErr *gorm.DB

	for _, user := range userList {
		var user_info *user_db.UserInfo
		result := database.DB.Where("email = ?", user).First(&user_info)

		if result.Error == nil {
			lastErr = result
			user_infoList = append(user_infoList, user_info)
		}
	}

	repositoryModels.Statue = *lastErr
	repositoryModels.Result = user_infoList

	return repositoryModels
}

func UpdateUserInfoRepository(user_info *user_db.UserInfo) repositoryModels.RepositoryModel[bool] {

	var repositoryModel repositoryModels.RepositoryModel[bool]

	err := database.DB.Save(&user_info)

	repositoryModel.Statue = *err

	repositoryModel.Result = err.Error == nil

	return repositoryModel
}
