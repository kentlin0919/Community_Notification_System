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
	"unicode/utf8"

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
	// registerModel 充當請求本文的映射容器，對應 Swagger 中宣告的註冊欄位
	var registerModel accountModel.Register
	// user_info 代表待寫入資料庫的核心實體，整合帳號、個資與登入憑證
	var user_info = &user_db.UserInfo{}
	// loginData 只保留登入所需欄位，用來檢查帳號是否已存在
	var loginData accountModel.User

	// 綁定 JSON 資料，若結構無法對應或缺欄位則直接回傳 400
	if err := ctx.ShouldBindJSON(&registerModel); err != nil {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "Invalid input")
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	if utf8.RuneCountInString(registerModel.Password) < 8 {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "密碼長度至少需 8 碼")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	// 將輸入的 Email 與 Password 複製到 loginData，以便後續查詢帳號是否已註冊
	loginData.Email = registerModel.Email
	loginData.Password = registerModel.Password
	checkaccountStatue := checkAccount(&loginData)

	if !checkaccountStatue {
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "已經註冊過了")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	// 加密密碼：即使資料庫外洩，亦可藉由單向雜湊降低密碼被破解的風險
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerModel.Password), bcrypt.DefaultCost)
	if err != nil {
		errorModel := model.NewErrorRequest(http.StatusInternalServerError, "密碼加密失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	// 建立唯一識別資訊與使用者基礎欄位，確保跨系統追蹤時的唯一性與一致性
	uuidString := uuid.New()
	user_info.ID = uuidString.String()
	user_info.Email = registerModel.Email
	user_info.Name = registerModel.Name
	user_info.Password = string(hashedPassword) // 存儲加密後的密碼
	user_info.BirthdayTime = registerModel.Birthday
	user_info.RegisterTime = time.Now()
	user_info.PermissionId = registerModel.Permission
	user_info.Platform = registerModel.Platform
	user_info.Session_id = uuid.New().String()
	// 為新使用者簽發 JWT，後續前端登入流程可直接沿用此 Token
	token, err := utils.GenerateJWT(user_info.Email, uint(user_info.PermissionId), user_info.ID)

	// 確認 JWT 是否成功簽發，避免回傳未簽名的憑證造成安全風險
	if err != nil {
		errorModel := model.NewErrorRequest(http.StatusInternalServerError, "JWT 簽發失敗")
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	user_info.Token = token

	// 透過 repository 實際寫入資料庫，統一封裝資料存取邏輯
	re := repository.RegisterRepository(user_info)

	if !re.Result {
		log.Fatalf("建立失敗")
		errorModel := model.NewErrorRequest(http.StatusBadRequest, "Register error")
		ctx.JSON(http.StatusBadRequest, errorModel)
	}

	var userRequest accountModel.UserRequest
	// 回應包含成功訊息與 JWT Token，供前端儲存並進行後續 API 驗證
	userRequest.Message = "Register successful"
	userRequest.Token = token
	// 設置 session_id Cookie，以配合前端瀏覽器維持短期會話狀態
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

func checkAccount(loginData *accountModel.User) bool {
	result := repository.LoginRepository(loginData)
	return result.Statue.Error == gorm.ErrRecordNotFound
}
