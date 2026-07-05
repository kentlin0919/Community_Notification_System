package auth

import (
	"errors"
	"time"

	accountModel "Community_Notification_System/app/models/account"
	repository "Community_Notification_System/app/repositories/user"
	authRepo "Community_Notification_System/app/repositories/auth"
	user_db "Community_Notification_System/database/User_DB"
	"Community_Notification_System/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credential")
	ErrDatabase          = errors.New("database error")
	ErrInternal          = errors.New("internal error")
)

type LoginResult struct {
	Token        string
	RefreshToken string
	SessionID    string
	UserInfo     accountModel.UserInfo
}

func ExecuteLogin(loginData *accountModel.User) (LoginResult, error) {
	result := repository.LoginRepository(loginData)
	if result.Statue.Error != nil {
		if result.Statue.Error == gorm.ErrRecordNotFound {
			return LoginResult{}, ErrUserNotFound
		}
		return LoginResult{}, ErrDatabase
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Result.Password), []byte(loginData.Password)); err != nil {
		return LoginResult{}, ErrInvalidCredential
	}

	token, err := utils.GenerateJWT(result.Result.Email, result.Result.ID, result.Result.PermissionId, result.Result.Community_id)
	if err != nil {
		return LoginResult{}, ErrInternal
	}

	sessionID := uuid.New().String()
	updateUserToken := repository.UpdateUserLoginStateRepository(result.Result.Email, token, sessionID, loginData.Fcmtoken)
	if updateUserToken.Statue.Error != nil {
		return LoginResult{}, ErrDatabase
	}

	logResult := repository.UserLogRepository(&result.Result)
	if logResult.Statue.Error != nil {
		return LoginResult{}, ErrDatabase
	}

	refreshToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return LoginResult{}, ErrInternal
	}
	session := &user_db.UserSession{
		ID:           uuid.New().String(),
		UserID:       result.Result.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	if err := authRepo.CreateSession(session).Statue.Error; err != nil {
		return LoginResult{}, ErrDatabase
	}

	return LoginResult{
		Token:        token,
		RefreshToken: refreshToken,
		SessionID:    sessionID,
		UserInfo: accountModel.UserInfo{
			PermissionId: result.Result.PermissionId,
			Name:         result.Result.Name,
			Email:        result.Result.Email,
			Home_id:      result.Result.Home_id,
			Birthdaytime: result.Result.Birthdaytime,
			PlatformID:   result.Result.Platform,
			Session_id:   sessionID,
			Community_id: result.Result.Community_id,
		},
	}, nil
}
