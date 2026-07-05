package auth

import (
	repository "Community_Notification_System/app/repositories/auth"
)

// Logout 撤銷指定 Refresh Token 對應的 session；找不到視為已登出，不回傳錯誤
func Logout(refreshToken string) error {
	result := repository.GetSessionByToken(refreshToken)
	if result.Statue.Error != nil || result.Result == nil {
		return nil
	}

	result.Result.IsRevoked = true
	updateResult := repository.UpdateSession(result.Result)
	return updateResult.Statue.Error
}
