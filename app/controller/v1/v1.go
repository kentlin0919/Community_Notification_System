package v1

import (
	"Community_Notification_System/app/controller/v1/communityManager"
	"Community_Notification_System/app/controller/v1/facility"
	"Community_Notification_System/app/controller/v1/message"
	"Community_Notification_System/app/controller/v1/permission"
	"Community_Notification_System/app/controller/v1/reservation"
	platform "Community_Notification_System/app/controller/v1/platform"
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

// Platform 回傳 PlatformController
func Platform() *platform.PlatformController {
	return platform.NewPlatformController()
}

// Permission 回傳權限設定控制器
func Permission() *permission.CommunityPermissionController {
	return permission.NewCommunityPermissionController()
}

// Facility 回傳設施控制器
func Facility() *facility.FacilityController {
	return facility.NewFacilityController()
}

// Reservation 回傳設施預約控制器
func Reservation() *reservation.ReservationController {
	return reservation.NewReservationController()
}
