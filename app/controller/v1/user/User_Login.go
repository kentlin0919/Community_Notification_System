package user

import (
	// "log"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	accountModel "Community_Notification_System/app/models/account"
	repository "Community_Notification_System/app/repositories/user"
)

// UserLogin 處理使用者登入
// @Summary 使用者登入
// @Description 使用者提供帳號與密碼後登入系統，並取得 JWT Token。
// @Tags User
// @Accept json
// @Produce json
// @Param login body account.Userlogin true "登入資料（Email & Password）"
// @Success 200 {object} account.UserloginRequest "登入成功訊息與 JWT Token"
// @Router /login [post]
func (u *UserController) UserLogin(ctx *gin.Context) {
	var loginData accountModel.Userlogin

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result := repository.LoginRepository(&loginData)

	if result.Statue.Error == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var request accountModel.UserloginRequest

	if result.Result.Password == loginData.Password {
		request.Message = "Login successful"
		request.Token = "11111111"
		ctx.JSON(http.StatusOK, request)
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
	}

}
