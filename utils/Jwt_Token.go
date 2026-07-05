package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(email string, userID string, permissionID int, communityID uint64) (string, error) {
	secret := os.Getenv("JWTPASSWORD")
	if len(secret) < 16 {
		return "", fmt.Errorf("JWT secret is not configured or too short")
	}
	JwtKey := []byte(secret)

	// payload 欄位
	claims := jwt.MapClaims{
		"username":      email,
		"user_id":       userID,
		"permission_id": permissionID,
		"community_id":  fmt.Sprintf("%d", communityID), // 轉為字串避免精度流失
		"exp":           time.Now().Add(30 * time.Minute).Unix(),
		"nbf":           time.Now().Unix(),
		"iat":           time.Now().Unix(),
	}

	// 建立 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密鑰簽名
	return token.SignedString(JwtKey)
}

// GenerateSecureToken 生成長度足夠且安全的隨機字串，適用於 Refresh Token 或 OTP
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
