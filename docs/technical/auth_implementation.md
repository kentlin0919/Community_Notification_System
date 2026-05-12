# 認證實作細節 (Auth Implementation)

## 1. 安全機制
- **Bcrypt**: 使用 `golang.org/x/crypto/bcrypt` 進行密碼雜湊，Cost 設定為 10。
- **JWT**: 採用 `github.com/golang-jwt/jwt` 簽發 Token。
  - Payload 包含 `user_id`, `community_id`, `role_id`, `exp`。
- **HTTP-Only Cookie**: 用於 Web 端，防止 XSS 攻擊。

## 2. Token 生命周期
- **Access Token**: 有效期 2 小時。
- **Refresh Token**: 有效期 7 天（規劃中）。
