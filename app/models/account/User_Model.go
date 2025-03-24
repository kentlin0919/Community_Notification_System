package account

import "time"

type Userlogin struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"yourpassword"`
}

type Register struct {
	Email    string    `json:"Email" example:"user@example.com"`
	Name     string    `json:"name" example:"kent"`
	Password string    `json:"password" example:"yourpassword"`
	Bethday  time.Time `json:"bethday" time_format:"2006-01-02T15:04:05-00:00" example:"2025-03-23T15:04:05-00:00"`
}

type UserloginRequest struct {
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token" example:"example-jwt-token"`
}
