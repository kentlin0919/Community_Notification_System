package user

import (
	accountModel "Community_Notification_System/app/models/account"
	repository "Community_Notification_System/app/repositories/user"
	user_db "Community_Notification_System/database/User_DB"
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

// UserRegister 處理使用者註冊
// @Summary 使用者註冊
// @Description 使用註冊
// @Tags User
// @Accept json
// @Produce json
// @Param Register body account.Register true "註冊"
// @Success 200 {object} map[string]interface{} "註冊成功訊息與 JWT Token"
// @Router /register [post]
func (u *UserController) UserRegister(ctx *gin.Context) {
	var registerModel accountModel.Register

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&registerModel); err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var loginData accountModel.Userlogin

	loginData.Email = registerModel.Email
	loginData.Password = registerModel.Password
	checkaccountStatue := checkAccount(&loginData)

	if !checkaccountStatue {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "已經註冊過了"})
		return
	}

	var user_info = &user_db.UserInfo{}

	uuidString := uuid.New()
	user_info.ID = uuidString.String()
	user_info.Email = registerModel.Email
	user_info.Name = registerModel.Name
	user_info.Password = registerModel.Password
	user_info.Birthdaytime = registerModel.Bethday
	user_info.Registertime = time.Now()

	re := repository.RegisterRepository(user_info)

	if re.Result {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Register successful",
			"token":   "example-jwt-token",
		})
	} else {
		log.Fatalf("建立失敗")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Regiser error"})
	}

}

func checkAccount(loginData *accountModel.Userlogin) bool {
	result := repository.LoginRepository(loginData)
	return result.Statue.Error == gorm.ErrRecordNotFound
}
