package v1

import (
	"Community_Notification_System/app/controller/v1/communityManager"
	"Community_Notification_System/app/controller/v1/message"
	"Community_Notification_System/app/controller/v1/user"
)

// User 回傳 UserController
func User() *user.UserController {
	return user.NewUserController()
}

func Message() *message.MessageController {
	return message.NewMessageController()
}

// CommunityManager 回傳 CommunityManagerController
func CommunityManager() *communityManager.CommunityManagerController {
	return communityManager.NewCommunityTableController()
}
