# 2025-10 Commit 摘要

- 2025-10-10 `feat(router): 根路徑導向 Swagger UI`
  - `main.go`: 新增 `/` 路由自動重新導向至 Swagger 與引入 `net/http`。
  - `README.md`: 說明根路徑將自動開啟 Swagger 介面。

- 2025-10-10 `feat(database): 自動建庫並種子預設資料`
  - `database/db.go`: 拆分 DSN 建構並於缺少資料庫時自動建立，同步掛載權限與平台資料表初始化。
  - `database/User_DB/User_Table.go`: 建立預設系統管理員帳號並調整欄位類型與匯入。
  - `database/Permission_DB/*`、`database/Platform_DB/*`: 新增權限與平台資料表 schema 與預設資料。
  - `app/controller/v1/user/User_Register.go`、`app/models/account/User_Model.go`: 對應欄位使用整數型別並連結新權限 id。
  - `docs/router_flow.md`: 補充路由流程與登入時序圖說明。
  - `AGENTS.md`、`GEMINI.md`: 更新協作指引與軟體架構說明。
  - `README.md`: 增加 Docker 啟動 PostgreSQL 範例指令。
