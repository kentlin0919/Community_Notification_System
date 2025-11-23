package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var JwtKey = []byte(os.Getenv("JWTPASSWORD"))

func GenerateJWT(email string, permissionID uint, userID string) (string, error) {
	// 建立 Claims
	claims := &Claims{
		Username:     email,
		PermissionID: permissionID,
		UserID:       userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(), // 過期時間：2小時
			IssuedAt:  time.Now().Unix(),                    // 簽發時間
		},
	}

	// 建立 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密鑰簽名
	return token.SignedString(JwtKey)
}

func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
