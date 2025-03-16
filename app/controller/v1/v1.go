package v1

import (
	"Community_Notification_System/app/controller/v1/user"
)

// User 回傳 UserController
func User() *user.UserController {
	return user.NewUserController()
}
