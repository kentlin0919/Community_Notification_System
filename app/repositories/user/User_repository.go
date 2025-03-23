package repository

import (
	"Community_Notification_System/app/models/account"
	repositoryModels "Community_Notification_System/app/models/repository"
	"Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"
	"log"
)

func LoginRepository(loginData *account.Userlogin) repositoryModels.RepositoryModel[user_db.UserInfo] {
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
