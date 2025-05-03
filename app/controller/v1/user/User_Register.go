package user

import (
	accountModel "Community_Notification_System/app/models/account"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/user"
	user_db "Community_Notification_System/database/User_DB"
	"Community_Notification_System/utils"
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

// UserRegister 處理使用者註冊
// @Summary 使用者註冊
// @Description 新用戶註冊功能，包含基本資料驗證、密碼加密、JWT Token 生成和 Session 設置
// @Tags User
// @Accept json
// @Produce json
// @Param Register body accountModel.Register true "註冊資料（Email、Password、Name、Birthday、Permission、Platform）"
// @Success 200 {object} accountModel.UserRequest "註冊成功，返回 JWT Token 和成功訊息"
// @Failure 400 {object} model.ErrorRequest "無效的輸入資料或帳號已存在"
// @Failure 500 {object} model.ErrorRequest "系統錯誤或 JWT 簽發失敗"
// @Router /api/v1/register [post]
func (u *UserController) UserRegister(ctx *gin.Context) {
	var registerModel accountModel.Register
	var user_info = &user_db.UserInfo{}
	var loginData accountModel.Userlogin

	// 綁定 JSON 資料
	if err := ctx.ShouldBindJSON(&registerModel); err != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "Invalid input"
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	loginData.Email = registerModel.Email
	loginData.Password = registerModel.Password
	checkaccountStatue := checkAccount(&loginData)

	if !checkaccountStatue {
		var errorModel model.ErrorRequest
		errorModel.Error = "已經註冊過了"
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerModel.Password), bcrypt.DefaultCost)
	if err != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "密碼加密失敗"
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	uuidString := uuid.New()
	user_info.ID = uuidString.String()
	user_info.Email = registerModel.Email
	user_info.Name = registerModel.Name
	user_info.Password = string(hashedPassword) // 存儲加密後的密碼
	user_info.Birthdaytime = registerModel.Bethday
	user_info.Registertime = time.Now()
	user_info.Permission = registerModel.Permission
	user_info.Platform = registerModel.Platform
	user_info.Session_id = uuid.New().String()
	token, err := utils.GenerateJWT(user_info.Email)

	//確認jwt簽發
	if err != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "JWT 簽發失敗"
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	user_info.Token = token

	re := repository.RegisterRepository(user_info)

	if !re.Result {
		log.Fatalf("建立失敗")
		var errorModel model.ErrorRequest
		errorModel.Error = "Regiser error"
		ctx.JSON(http.StatusBadRequest, errorModel)
	}

	var userRequest accountModel.UserRequest
	userRequest.Message = "Register successful"
	userRequest.Token = token
	ctx.SetCookie(
		"session_id",
		user_info.Session_id,
		3600, // 1小時過期
		"/",
		"",
		true, // 只在 HTTPS 下傳輸
		true, // 防止 JavaScript 訪問
	)
	ctx.JSON(http.StatusOK, userRequest)

}

func checkAccount(loginData *accountModel.Userlogin) bool {
	result := repository.LoginRepository(loginData)
	return result.Statue.Error == gorm.ErrRecordNotFound
}
