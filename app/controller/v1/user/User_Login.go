package user

import (
	// "log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	accountModel "Community_Notification_System/app/models/account"
	"Community_Notification_System/app/models/model"
	repository "Community_Notification_System/app/repositories/user"
	"Community_Notification_System/utils"
)

// UserLogin 處理使用者登入
// @Summary 使用者登入
// @Description 使用者提供帳號與密碼後登入系統，並取得 JWT Token。登入成功後會設置 session cookie 並返回 JWT token。
// @Tags User
// @Accept json
// @Produce json
// @Param login body accountModel.User true "登入資料（Email & Password）"
// @Success 200 {object} accountModel.UserRequest "登入成功，返回 JWT Token 和成功訊息"
// @Failure 400 {object} model.ErrorRequest "無效的輸入資料"
// @Failure 401 {object} model.ErrorRequest "密碼錯誤"
// @Failure 404 {object} model.ErrorRequest "使用者不存在"
// @Failure 500 {object} model.ErrorRequest "系統錯誤或 JWT 簽發失敗"
// @Router /api/v1/login [post]
func (u *UserController) UserLogin(ctx *gin.Context) {
	var loginData accountModel.User

	// 綁定 JSON 資料並驗證輸入格式
	// 使用 ShouldBindJSON 可以自動驗證 JSON 格式是否符合結構體定義
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "無效的輸入資料"
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorModel})
		return
	}

	// 查詢用戶資料
	// 使用 repository 模式分離資料庫操作邏輯
	result := repository.LoginRepository(&loginData)

	// 處理資料庫錯誤
	// 1. 檢查是否有任何資料庫錯誤
	// 2. 特別處理用戶不存在的情況
	// 3. 處理其他可能的資料庫錯誤
	if result.Statue.Error != nil {
		if result.Statue.Error == gorm.ErrRecordNotFound {
			var errorModel model.ErrorRequest
			errorModel.Error = "使用者不存在"
			ctx.JSON(http.StatusUnauthorized, errorModel)
			return
		}
		var errorModel model.ErrorRequest
		errorModel.Error = "系統錯誤"
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	// 使用 bcrypt 進行密碼驗證
	// 1. 避免明文密碼比較
	// 2. 防止時序攻擊
	// 3. 符合安全最佳實踐
	if err := bcrypt.CompareHashAndPassword([]byte(result.Result.Password), []byte(loginData.Password)); err != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "帳號或密碼錯誤"
		ctx.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	// 生成 JWT Token
	// 用於後續的身份驗證和授權
	token, err := utils.GenerateJWT(result.Result.Email)
	if err != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "JWT 簽發失敗"
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	// 生成唯一的 Session ID
	// 使用 UUID 確保 session ID 的唯一性和安全性
	sessionID := uuid.New().String()

	// 設置安全的 Cookie
	// 1. 使用 HTTPS only
	// 2. 防止 JavaScript 訪問
	// 3. 設置合理的過期時間
	ctx.SetCookie(
		"session_id",
		sessionID,
		3600, // 1小時過期
		"/",
		"",
		true, // 只在 HTTPS 下傳輸
		true, // 防止 JavaScript 訪問
	)

	// 更新用戶最後登入時間
	// 用於追蹤用戶活動和安全性監控
	logResult := repository.UserLogRepository(&result.Result)
	if logResult.Statue.Error != nil {
		var errorModel model.ErrorRequest
		errorModel.Error = "更新用戶最後登入時間失敗"
		ctx.JSON(http.StatusInternalServerError, errorModel)
		return
	}

	// 返回登入成功響應
	var request accountModel.UserRequest
	request.Message = "登入成功"
	request.Token = token
	ctx.JSON(http.StatusOK, request)
}
