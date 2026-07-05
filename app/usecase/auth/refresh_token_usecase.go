package auth

import (
	auth_model "Community_Notification_System/app/models/auth"
	auth_repo "Community_Notification_System/app/repositories/auth"
	user_db "Community_Notification_System/database/User_DB"
	"Community_Notification_System/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken 執行 Refresh Token 的核心業務邏輯
func RefreshToken(req auth_model.RefreshTokenRequest) (*auth_model.TokenResponse, error) {
	var resp *auth_model.TokenResponse

	err := auth_repo.TransactionWrapper(func(tx *gorm.DB) error {
		// 1. 取得 session 紀錄 (包含 FOR UPDATE 行級鎖)
		var session user_db.UserSession
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("refresh_token = ?", req.RefreshToken).First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("invalid refresh token")
			}
			return err
		}

		// 2. 檢查 Reuse Detection 與 Grace Period
		if session.IsRevoked {
			gracePeriod := 30 * time.Second
			if time.Since(session.UpdatedAt) < gracePeriod {
				return errors.New("token already used, waiting for client sync")
			} else {
				auth_repo.RevokeAllUserSessions(tx, session.UserID)
				return errors.New("token reuse detected, all sessions revoked")
			}
		}

		// 3. 檢查過期時間
		if time.Now().After(session.ExpiresAt) {
			return errors.New("refresh token expired")
		}

		// 4. 將舊 Session 標記為 revoked
		if err := tx.Model(&session).Update("is_revoked", true).Error; err != nil {
			return err
		}

		// 5. 取得使用者資訊以簽發新 Access Token
		var user user_db.UserInfo
		if err := tx.Where("id = ?", session.UserID).First(&user).Error; err != nil {
			return errors.New("user not found")
		}

		newAccessToken, err := utils.GenerateJWT(user.Email, user.ID, user.PermissionId, user.Community_id)
		if err != nil {
			return err
		}

		newRefreshTokenString, err := utils.GenerateSecureToken(32)
		if err != nil {
			return err
		}

		// 6. 建立新的 Session
		newSession := user_db.UserSession{
			ID:           uuid.New().String(),
			UserID:       user.ID,
			RefreshToken: newRefreshTokenString,
			DeviceID:     session.DeviceID,
			DeviceInfo:   session.DeviceInfo,
			ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // RT 效期 30 天
			IsRevoked:    false,
		}

		if err := tx.Create(&newSession).Error; err != nil {
			return err
		}

		// 7. 準備回傳結構
		resp = &auth_model.TokenResponse{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshTokenString,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
