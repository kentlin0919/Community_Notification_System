# 認證實作細節 (Auth Implementation)

更新日期：2026-05-12
適用範圍：Community_Notification_System_backend

## 1. 目前實作總覽

目前系統使用 JWT + Session Cookie 混合模式：
- Login/Register 成功後回傳 JWT（response body）
- 同時設定 `session_id` Cookie（HttpOnly + Secure）
- API 認證主要透過 `Authorization: Bearer <token>`

對應程式：
- JWT 產生：`utils/Jwt_Token.go`
- JWT 驗證：`middlewares/jwt_middleware.go`
- 登入流程：`app/controller/v1/user/User_Login.go`
- 註冊流程：`app/controller/v1/user/User_Register.go`

## 2. 安全機制

### 2.1 密碼雜湊
- 使用 `golang.org/x/crypto/bcrypt`
- 雜湊：`bcrypt.GenerateFromPassword(..., bcrypt.DefaultCost)`
- 驗證：`bcrypt.CompareHashAndPassword(...)`

### 2.2 JWT
- 套件：`github.com/golang-jwt/jwt`
- 演算法：HS256
- 金鑰來源：環境變數 `JWTPASSWORD`
- Claims（目前）：
  - `username`
  - `user_id`
  - `permission_id`
  - `community_id`
  - `iat`
  - `exp`

### 2.3 Cookie
- Cookie 名稱：`session_id`
- 屬性：`Secure=true`, `HttpOnly=true`
- 期限：3600 秒（1 小時）

## 3. Token 生命週期

- Access Token：2 小時（`GenerateJWT`）
- Refresh Token：尚未實作（目前無 refresh endpoint）

## 4. 已知風險與待修正項

### 4.1 Middleware 控制流語意需明確化
- `JWTAuthMiddleware` 成功驗證後，應明確 `c.Next()`。
- 雖然 Gin 在 handler return 後通常會繼續 chain，但建議顯式呼叫以降低維護風險。

### 4.2 全域 JWT + skipPaths 白名單模式風險
- 目前在 middleware 內使用硬編碼 path 白名單（如 `/api/v1/login`, `/api/v1/register`）。
- 風險：路由異動或新公開 API 易漏改，造成誤攔截或授權漏洞。
- 建議：改為 public/private route group 分離。

### 4.3 認證策略混合（JWT + session）定義不夠清楚
- 目前同時回 token 並設 cookie，但沒有明確宣告主要授權來源與撤銷策略。
- 建議：
  - 定義單一路徑（純 JWT 或 session-based）；
  - 或完整定義雙軌策略（優先順序、過期、撤銷、登出）。

### 4.4 錯誤回應格式尚未完全統一
- 認證相關 API 目前存在 `model.NewErrorRequest`、`NewGlobalErrorRequestWithMsg`、`gin.H` 混用情況。
- 建議：統一 ErrorResponse contract（含 `request_id`）。

## 5. 建議修正順序（短期）

1. 修正 JWT middleware 成功分支流程（顯式 `c.Next()`）
2. 將路由改為 public/private 分組，移除 path 白名單邏輯
3. 統一 auth 錯誤回應格式
4. 補齊 middleware/login 測試案例（401/200/expired/invalid format）

## 6. 後續規劃（中期）

- 新增 Refresh Token 機制（7 天）
- 加入 token 撤銷策略（logout / 強制下線）
- 導入 login rate limit（IP + account）
- claims 型別轉換集中封裝（避免各處手動轉型）
