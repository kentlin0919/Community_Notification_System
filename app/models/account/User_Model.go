package account

import "time"

type User struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"yourpassword"`
	Platform string `json:"platform" example:"App"`
}

type Register struct {
	Email      string    `json:"Email" example:"user@example.com"`
	Name       string    `json:"name" example:"kent"`
	Password   string    `json:"password" example:"yourpassword"`
	Bethday    time.Time `json:"bethday" time_format:"2006-01-02T15:04:05-00:00" example:"2025-03-23T15:04:05-00:00"`
	Permission int       `json:"permission" example:"1"`
	Platform   int       `json:"platform" example:"App"`
}

type UserRequest struct {
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token" example:"example-jwt-token"`
}
