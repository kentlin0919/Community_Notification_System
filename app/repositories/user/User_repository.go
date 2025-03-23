package repository

import (
	"Community_Notification_System/app/models/account"
	"Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"

	"gorm.io/gorm"
)

func LoginRepository(loginData *account.Userlogin) *user_db.UserInfo {
	var user_info user_db.UserInfo

	result := database.DB.Where("email = ?", loginData.Email).First(&user_info)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {

		} else {

		}

	}

	return &user_info
}
