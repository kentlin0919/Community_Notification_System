package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var JwtKey = []byte(os.Getenv("JWTPASSWORD"))

func GenerateJWT(email string) (string, error) {
	// payload 欄位
	claims := jwt.MapClaims{
		"username": email,
		"exp":      time.Now().Add(2 * time.Hour).Unix(), // 過期時間：2小時
		"iat":      time.Now().Unix(),                    // 簽發時間
	}

	// 建立 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密鑰簽名
	return token.SignedString(JwtKey)
}
