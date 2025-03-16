package user

import (
	// "log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	database "Community_Notification_System/database"
	user_db "Community_Notification_System/database/User_DB"

	accountModel "Community_Notification_System/app/models/account"
)

func (u *UserController) Login(ctx *gin.Context) {
	var loginData accountModel.Userlogin

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user_info user_db.UserInfo

	result := database.DB.Where("email = ?", loginData.Email).First(&user_info)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// 簡單的帳號密碼驗證（範例）
	if user_info.Password == loginData.Password {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   "example-jwt-token",
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
	}

}

func (u *UserController) DeleteUser(ctx *gin.Context) {
	var registerModel accountModel.Register

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&registerModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

}
