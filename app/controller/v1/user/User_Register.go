package user

import (
	accountModel "Community_Notification_System/app/models/account"
	database "Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *UserController) UserRegister(ctx *gin.Context) {
	var registerModel accountModel.Register

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&registerModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user_info = &user_db.UserInfo{}

	user_info.Email = registerModel.Email
	user_info.Name = registerModel.Name
	user_info.Password = registerModel.Password

	if err := database.DB.Create(user_info).Error; err != nil {
		log.Fatalf("建立失敗")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Regiser error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Register successful",
		"token":   "example-jwt-token",
	})

}
