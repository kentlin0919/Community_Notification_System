package repository

import (
	"Community_Notification_System/app/models/account"
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	userlog_db "Community_Notification_System/database/UserLog_DB"
	user_db "Community_Notification_System/database/User_DB"
	"log"
	"time"
)

func LoginRepository(loginData *account.User) repositoryModels.RepositoryModel[user_db.UserInfo] {
	var user_info user_db.UserInfo
	var repositoryModel repositoryModels.RepositoryModel[user_db.UserInfo]

	result := database.DB.Where("email = ?", loginData.Email).First(&user_info)

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

	err := database.DB.Delete(&user_info)

	var repositoryModel repositoryModels.RepositoryModel[bool]

	repositoryModel.Statue = *err

	repositoryModel.Result = err.Error == nil

	return repositoryModel
}
