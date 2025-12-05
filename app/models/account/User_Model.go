package account

import "time"

type User struct {
	Email    string `json:"email" example:"kent900919@gmail.com"`
	Password string `json:"password" example:"09190919"`
	Platform string `json:"platform" example:"App"`
	Fcmtoken string `json:"fcmtoken" example:"example-fcm-token"`
}

type Register struct {
	Email       string    `json:"Email" example:"user@example.com"`
	Name        string    `json:"name" example:"kent"`
	Password    string    `json:"password" example:"yourpassword"`
	Bethday     time.Time `json:"bethday" time_format:"2006-01-02T15:04:05-00:00" example:"2025-03-23T15:04:05-00:00"`
	Permission  int       `json:"permission" example:"1"`
	Platform    int       `json:"platform" example:"1"`
	CommunityID uint64    `json:"community_id" example:"1"`
	Home_id     uint64    `json:"home_id" example:"1"`
}

type UserRequest struct {
	Message  string   `json:"message" example:"Login successful"`
	Token    string   `json:"token" example:"example-jwt-token"`
	UserInfo UserInfo `json:"user_info"`
}

type UserInfo struct {
	PermissionId int       `json:"PermissionId"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Home_id      uint64    `json:"Home_id"`
	Birthdaytime time.Time `json:"Birthdaytime" example:"2025-03-23T15:04:05Z"`
	PlatformID   int       `json:"Platform"`
	Session_id   string    `json:"Session_id"`
	Community_id uint64    `json:"community_id"`
}
