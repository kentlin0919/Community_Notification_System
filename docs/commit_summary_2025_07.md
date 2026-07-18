# Commit Summary 2025-07

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
