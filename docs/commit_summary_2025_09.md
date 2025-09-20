# 2025-09 Commit 摘要

- 2025-09-20 `feat(api): 擴增 v2 路由並補齊測試與文件`
  - `README.md`: 改寫為完整的專案導覽，新增版本資訊、架構說明、時序圖與平台安裝流程。
  - `routers/router.go`、`routers/api/v2/v2.go`: 開啟 `/api/v2` 路徑並沿用 v1 控制器，讓新版路由可登入、註冊、刪除與發送訊息。
  - `app/controller/v1/user/User_Login_test.go`: 建立登入流程的單元測試，覆蓋成功、輸入錯誤、密碼錯誤與找不到使用者等情境。
  - `scripts/install_dependencies.sh`: 提供 macOS/Linux 開發環境安裝腳本，涵蓋 Go、PostgreSQL、Air 與 Swag CLI。
  - `AGENTS.md`、`GEMINI.md`: 補充協作說明與快速導讀。
  - `go.mod`、`go.sum`: 納入測試所需的 SQLite 依賴並調整套件引用範疇。
