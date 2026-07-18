# Commit Summary 2025-07

## 2026-07-18 14:10 - `feat(community): 新增社區申請列表與更新文件路徑`

### 變更原因

- 前端 Super Admin 需要以 `/api/v1/super-admin/community-applications` 取得社區申請列表，並依申請狀態篩選。
- 專案文件路徑已搬到 iCloud Obsidian 目錄，需要同步更新 AGENTS、CLAUDE、GEMINI 與 README 中的文件中心位置。
- CORS 需要允許前端常用的 `PUT`、`PATCH`、`DELETE` 與 `Authorization` header。
- Swagger 文件需要反映目前 Auth、Platform 與 Super Admin 相關 API。

### 逐行分析

- `app/controller/v1/communityManager/CommunityManager_GetApplicationList.go`：新增社區申請列表 controller，讀取 query `status` 作為篩選條件，回傳前端需要的申請摘要欄位。
- `app/repositories/community/CommunityRegisterApplication.go`：新增 `GetApplicationsRepository`，可依狀態查詢社區註冊申請資料。
- `routers/api/v1/v1.go`：新增 Super Admin 對齊前端命名的社區與社區申請路由，包含列表、核准與拒絕。
- `middlewares/cors_middleware.go`：允許 `GET, POST, PUT, PATCH, DELETE, OPTIONS`，並允許 `Authorization` header。
- `docs/docs.go`、`docs/swagger.json`、`docs/swagger.yaml`：同步 Swagger 產物，補上 Auth、Platform 與 Super Admin community applications API 文件。
- `AGENTS.md`、`CLAUDE.md`、`GEMINI.md`、`README.md`：將開發文件中心位置更新為 `/Users/kent/Library/Mobile Documents/iCloud~md~obsidian/Documents/Community_Notification_System_docs`。
- 其他 Go 檔案：套用 gofmt/import 排序與空白整理，不改變既有業務邏輯。

### 驗證

- 已執行 `gofmt` 套用 Go 格式化。
- 已執行 `go test ./...` 驗證全部 Go 測試。

## 2026-07-18 13:55 - `fix(debug): 修正 Docker 遠端偵錯設定`

### 變更原因

- 登入 API 在 Docker 環境中因 `JWTPASSWORD` 預設值長度不足，導致 JWT 簽發失敗並回傳 500。
- VS Code attach Delve 時因 Delve 版本過舊，曾出現遠端偵錯連線中斷。
- Swagger 自動開啟條件與 Gin 實際啟動訊息不一致，導致偵錯啟動後沒有跳出網頁。
- 自動 `postDebugTask` 會在 attach 失敗或偵錯中斷時刪除容器，造成反覆重建。

### 逐行分析

- `Dockerfile`：將 `DLV_VERSION` 由 `v1.24.2` 升級為 `v1.27.0`，讓 Delve 支援目前 Go 1.26.1 的遠端偵錯流程。
- `docker-compose.yml`：將 `JWTPASSWORD=your_jwt_secret` 改為 `${JWTPASSWORD:-community_dev_jwt_secret_2026}`，優先讀取外部環境變數，並提供長度足夠的開發預設值。
- `.vscode/launch.json`：將 `serverReadyAction.pattern` 從 `listening on` 改為 `Listening and serving HTTP on`，匹配 Gin 實際啟動輸出。
- `README.md`：更新 macOS、Linux、Windows 的 `.env` 範例，避免文件繼續提供過短 JWT secret。
- `README.md`：更新 Docker Debug 說明，記錄 Swagger 會在服務啟動後自動開啟。
- `README.md`：補充偵錯中斷或 attach 失敗時不會自動刪除容器，需手動執行 VS Code Task `docker-compose: down` 清理。

### 驗證

- 已驗證 `.vscode/launch.json` 與 `.vscode/tasks.json` JSON 格式正確。
- 已驗證 Docker 容器內 `JWTPASSWORD` 長度大於 16。
- 已驗證 Delve 版本為 `1.27.0`。
- 已驗證 Swagger `http://localhost:9080/swagger/index.html` 回傳 200。
- 已驗證登入 API 可成功回傳 access token 與 refresh token。
