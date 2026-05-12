# Login/Auth 結構修正實作計畫

> For Hermes: 使用 subagent-driven-development 逐任務執行，任務完成後跑測試與 pre-commit review。

目標：修正目前登入/認證流程的高風險設計，讓行為一致、可測試、可維護，並與 clean architecture 方向對齊。

架構策略：先做最小風險修正（middleware 與 route 分組），再做錯誤回應一致化，最後把登入流程中 controller 直連 DB 拆到 usecase/repository。

技術範圍：Gin, Gorm, JWT, middleware, controller/repository 分層。

---

## Task 1: 明確化 JWT middleware 成功分支流程

Objective: 避免 middleware chain 行為依賴隱性機制。

Files:
- Modify: `middlewares/jwt_middleware.go`
- Test: `middlewares/jwt_middleware_test.go`

Step 1: 在 token 驗證成功且 claims 寫入 context 後，顯式加入 `c.Next()` 與 `return`。
Step 2: 確認失敗分支（header 缺失、格式錯誤、token 無效）皆為 `Abort()`。
Step 3: 補測試（成功請求可通過到 handler）。
Step 4: 執行測試
- `go test ./middlewares -v`

完成定義:
- middleware 測試通過
- 成功與失敗路徑行為可預期

---

## Task 2: 將認證改為 route group 管理（public/private）

Objective: 移除 path 白名單硬編碼，降低授權漏洞風險。

Files:
- Modify: `main.go`
- Modify: `routers/api/v1/v1.go`（若需拆分註冊函式）
- Modify: `middlewares/jwt_middleware.go`

Step 1: 在 router 設計 public/private group。
- public: `/login`, `/register`, swagger, root
- private: 其他需授權 API
Step 2: private group 才掛 JWT middleware；移除/極小化 `skipPaths`。
Step 3: 跑 smoke test，確認公開路由不需 token、私有路由需 token。

驗證:
- `go test ./... -v`
- 手動驗證 4 條路由：
  - POST /api/v1/login (無 token 應可呼叫)
  - POST /api/v1/register (無 token 應可呼叫)
  - GET /api/v1/messages (無 token 應 401)
  - GET /api/v1/messages (合法 token 應可呼叫)

---

## Task 3: 統一認證錯誤回應格式

Objective: API contract 一致，方便前端與監控。

Files:
- Modify: `app/models/model/model.go`（若需補 helper）
- Modify: `middlewares/jwt_middleware.go`
- Modify: `app/controller/v1/user/User_Login.go`
- Modify: `app/controller/v1/user/User_Register.go`

Step 1: 定義 auth 相關固定錯誤回應形式（code/internal_code/status/error）。
Step 2: middleware 內將 `gin.H` 全部改為統一 model。
Step 3: login/register 也套同一回應工廠。
Step 4: 補/調整 swagger 註解與錯誤碼對應。

驗證:
- `go test ./... -v`
- 失敗案例回應 JSON key 應一致

---

## Task 4: 將 UserLogin 的 DB 直連移出 controller

Objective: 收斂商業邏輯到 usecase/service 層，貼近 clean architecture。

Files:
- Create: `app/usecase/auth/login_usecase.go`（或專案慣用層）
- Modify: `app/controller/v1/user/User_Login.go`
- Modify: `app/repositories/user/*`（新增更新 token/session/fcm 的 repository 方法）

Step 1: 建立 LoginUsecase，封裝：
- 比對密碼
- 產生 JWT/session
- 更新 token/session/fcm
- 更新最後登入時間
Step 2: controller 僅保留：bind request、呼叫 usecase、回應 mapping。
Step 3: 補 usecase 單元測試（可 mock repository）。

驗證:
- `go test ./... -v`
- 功能與既有 API 行為一致

---

## Task 5: 補齊登入/認證最小測試矩陣

Objective: 防止回歸並保障重構安全。

Files:
- Modify/Create: `app/controller/v1/user/User_Login_test.go`
- Modify/Create: `middlewares/jwt_middleware_test.go`

必備測試:
1) login 帳號不存在 -> 404
2) login 密碼錯誤 -> 401
3) login DB 錯誤 -> 500
4) login 成功 -> 200 且回 token
5) JWT 缺 header -> 401
6) JWT Bearer 格式錯 -> 401
7) JWT 過期/無效 -> 401
8) JWT 有效 -> 可進 handler

驗證:
- `go test ./... -v`

---

## Task 6: 文件同步更新

Objective: 文件與實作一致，避免認知落差。

Files:
- Modify: `docs/technical/auth_implementation.md`
- Modify: `README.md`（認證章節）
- Modify: `/docs` 下相關索引

Step 1: 更新架構圖與 auth 流程說明（public/private route + JWT 行為）。
Step 2: 更新錯誤回應範例。
Step 3: 如 API contract 變更，重跑 swagger。
- `swag init -g main.go`

---

## 風險控管與回滾

- 每個 task 完成都獨立 commit，避免大批次風險。
- 任何 task 若造成認證全面失效，立即回滾到上一 commit。
- 先保證測試矩陣綠燈，再進下一 task。

建議 commit 範例:
- `fix(auth): 明確 JWT middleware 成功分支 c.Next 控制流`
- `ref(auth): 拆分 public/private routes 移除 path 白名單`
- `ref(auth): 抽離 login usecase 移除 controller 直連 DB`
- `test(auth): 補齊 login 與 JWT middleware 測試矩陣`
- `docs(auth): 同步更新認證實作與 README`

---

## 驗收標準（Definition of Done）

1. 認證流程可由測試矩陣完整覆蓋，`go test ./... -v` 通過。
2. 不再依賴 middleware path 白名單做公開路由控制。
3. controller 不直接更新 DB（登入主要流程由 usecase/repository 接手）。
4. auth 錯誤回應 JSON 格式一致。
5. docs/README 與程式實作一致。
