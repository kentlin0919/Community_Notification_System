package user

import (
	"net/http"

	accountModel "Community_Notification_System/app/models/account"
	"Community_Notification_System/app/models/model"
	authUsecase "Community_Notification_System/app/usecase/auth"
	"github.com/gin-gonic/gin"

	utilsErr "Community_Notification_System/utils/errors"
)

// UserLogin 處理使用者登入
// @Summary 使用者登入
// @Description 使用者提供帳號與密碼後登入系統，並取得 JWT Token。登入成功後會設置 session cookie 並返回 JWT token。
// @Tags User
// @Accept json
// @Produce json
// @Param login body accountModel.User true "登入資料（Email, Password, Fcmtoken）"
// @Success 200 {object} accountModel.UserRequest "登入成功，返回 JWT Token 和使用者詳細資訊"
// @Failure 400 {object} model.Response400Error "無效的輸入資料"
// @Failure 401 {object} model.Response401Error "帳號或密碼錯誤"
// @Failure 404 {object} model.Response404Error "使用者不存在"
// @Failure 500 {object} model.Response500Error "系統錯誤、JWT 簽發失敗或更新登入狀態失敗"
// @Router /api/v1/login [post]
func (u *UserController) UserLogin(ctx *gin.Context) {
	var loginData accountModel.User

	// 綁定 JSON 資料並驗證輸入格式
	// 使用 ShouldBindJSON 可以自動驗證 JSON 格式是否符合結構體定義
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		errorModel := model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "無效的輸入資料")
		ctx.JSON(http.StatusBadRequest, errorModel)
		return
	}

	loginResult, err := authUsecase.ExecuteLogin(&loginData)
	if err != nil {
		switch err {
		case authUsecase.ErrUserNotFound:
			errorModel := model.NewErrorResponse(ctx, http.StatusNotFound, utilsErr.ErrUserNotFound, "使用者不存在")
			ctx.JSON(http.StatusNotFound, errorModel)
			return
		case authUsecase.ErrInvalidCredential:
			errorModel := model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrPasswordWrong, "帳號或密碼錯誤")
			ctx.JSON(http.StatusUnauthorized, errorModel)
			return
		case authUsecase.ErrDatabase:
			errorModel := model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrDatabase, "系統錯誤")
			ctx.JSON(http.StatusInternalServerError, errorModel)
			return
		default:
			errorModel := model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "JWT 簽發失敗")
			ctx.JSON(http.StatusInternalServerError, errorModel)
			return
		}
	}

	// 設置安全的 Cookie
	// 1. 使用 HTTPS only
	// 2. 防止 JavaScript 訪問
	// 3. 設置合理的過期時間
	ctx.SetCookie(
		"session_id",
		loginResult.SessionID,
		3600, // 1小時過期
		"/",
		"",
		true, // 只在 HTTPS 下傳輸
		true, // 防止 JavaScript 訪問
	)

	// 返回登入成功響應
	request := accountModel.UserRequest{
		Message:      "登入成功",
		Token:        loginResult.Token,
		RefreshToken: loginResult.RefreshToken,
		UserInfo:     loginResult.UserInfo,
	}
	ctx.JSON(http.StatusOK, request)
}
