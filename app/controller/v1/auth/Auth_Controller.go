package auth

import (
	authModel "Community_Notification_System/app/models/auth"
	"Community_Notification_System/app/models/model"
	authUsecase "Community_Notification_System/app/usecase/auth"
	"Community_Notification_System/pkg/email"
	"net/http"

	utilsErr "Community_Notification_System/utils/errors"
	"github.com/gin-gonic/gin"
)

// AuthController 認證與安全控制器
type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// RefreshToken 刷新 Access Token
// @Summary 刷新 Access Token
// @Description 接收 Refresh Token，驗證有效性後核發新 Access Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authModel.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} authModel.TokenResponse "成功返回新的 Token"
// @Failure 401 {object} model.Response401Error "無效或過期的 Token"
// @Router /api/v1/auth/refresh [post]
func (a *AuthController) RefreshToken(ctx *gin.Context) {
	var req authModel.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "無效的輸入資料"))
		return
	}

	// 實作 Refresh Token 驗證邏輯
	resp, err := authUsecase.RefreshToken(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(ctx, http.StatusUnauthorized, utilsErr.ErrUnauthorized, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// SwitchCommunity 切換社區權限
// @Summary 切換社區權限
// @Description 驗證使用者對目標社區的權限，重新核發帶有新 community_id 的 Access Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authModel.SwitchCommunityRequest true "目標社區 ID"
// @Success 200 {object} authModel.TokenResponse "成功返回新的 Token"
// @Failure 401 {object} model.Response401Error "無權限存取該社區"
// @Security BearerAuth
// @Router /api/v1/auth/switch-community [post]
func (a *AuthController) SwitchCommunity(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "功能尚未實作"})
}

// ForgotPassword 忘記密碼發送 OTP
// @Summary 忘記密碼發送 OTP
// @Description 產生 6 位數 OTP，存入 Redis (TTL 10m)，並模擬 Email 發送
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authModel.ForgotPasswordRequest true "使用者 Email"
// @Success 200 "發送成功"
// @Router /api/v1/auth/forgot-password [post]
func (a *AuthController) ForgotPassword(ctx *gin.Context) {
	var req authModel.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(ctx, http.StatusBadRequest, utilsErr.ErrInvalidParams, "無效的輸入資料"))
		return
	}

	if err := authUsecase.ForgotPassword(req.Email, email.NewLogEmailSender()); err != nil {
		ctx.JSON(http.StatusInternalServerError, model.NewErrorResponse(ctx, http.StatusInternalServerError, utilsErr.ErrInternal, "系統錯誤"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "如該 Email 存在，驗證碼已寄出"})
}

// VerifyOTP 驗證密碼重設 OTP
// @Summary 驗證密碼重設 OTP
// @Description 驗證 OTP，成功後核發限重設密碼用的 ResetToken
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authModel.VerifyOTPRequest true "Email 與 6 位數 OTP"
// @Success 200 {object} authModel.ResetTokenResponse "驗證成功，返回 ResetToken"
// @Failure 401 {object} model.Response401Error "OTP 錯誤或過期"
// @Router /api/v1/auth/verify-otp [post]
func (a *AuthController) VerifyOTP(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "功能尚未實作"})
}

// ResetPassword 重設密碼
// @Summary 重設密碼
// @Description 驗證 ResetToken，更新密碼，並撤銷舊的 Refresh Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authModel.ResetPasswordRequest true "ResetToken 與 新密碼"
// @Success 200 "重設密碼成功"
// @Failure 401 {object} model.Response401Error "ResetToken 無效"
// @Router /api/v1/auth/reset-password [post]
func (a *AuthController) ResetPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "功能尚未實作"})
}
