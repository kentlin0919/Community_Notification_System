package utils

import "github.com/golang-jwt/jwt"

type Claims struct {
	Username     string `json:"username"`
	PermissionID uint   `json:"permission_id"`
	UserID       string `json:"user_id"`
	jwt.StandardClaims
}
