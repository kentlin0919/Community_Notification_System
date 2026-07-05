package auth

import (
	"errors"
	"time"

	repository "Community_Notification_System/app/repositories/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrResetTokenInvalid = errors.New("reset token invalid or expired")

// ResetPassword 驗證 ResetToken、更新密碼並撤銷該使用者所有既有 session
func ResetPassword(resetToken string, newPassword string) error {
	tokenRecord, err := repository.GetPasswordResetTokenByResetHash(hashToken(resetToken))
	if err != nil {
		return err
	}
	if tokenRecord == nil || !tokenRecord.Verified || time.Now().After(tokenRecord.ExpiresAt) {
		return ErrResetTokenInvalid
	}

	userResult := repository.GetUserByEmail(tokenRecord.Email)
	if userResult.Statue.Error != nil || userResult.Result == nil {
		return ErrResetTokenInvalid
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return repository.TransactionWrapper(func(tx *gorm.DB) error {
		if err := tx.Model(userResult.Result).Update("password", string(hashedPassword)).Error; err != nil {
			return err
		}
		if err := repository.RevokeAllUserSessions(tx, userResult.Result.ID); err != nil {
			return err
		}
		tokenRecord.Used = true
		return tx.Save(tokenRecord).Error
	})
}
