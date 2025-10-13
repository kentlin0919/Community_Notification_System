# 2025-10 Commit 摘要

- 2025-10-10 `feat(platform): add platform listing endpoint`
  - `app/controller/v1/platform/*`: 新增平台控制器與取得平台列表的 API。
  - `app/repositories/platform/Platform_repository.go`: 建立平台資料存取邏輯，回傳平台清單資訊。
  - `routers/api/v1/v1.go`、`app/controller/v1/v1.go`: 註冊平台路由並調整社區列表路徑。
  - `app/repositories/community/community_repository.go`: 重構查詢條件與回傳結構以配合新的列表需求。
  - `docs/docs.go`、`docs/swagger.json`、`docs/swagger.yaml`: 更新 Swagger 文件描述平台與社區列表端點。
  - `__debug_bin892389873`: 移除不必要的除錯二進位檔案。

- 2025-10-10 `refactor(model): 統一錯誤回應結構`
  - `app/models/model/model.go`: ErrorRequest 新增 `code`、`status` 欄位並提供建構函式。
  - `app/controller/v1/user/*`、`app/controller/v1/message/Message_SendMessage.go`: 改用新錯誤格式回應並對應 HTTP 狀態碼。
  - `app/controller/v1/user/User_Login_test.go`: 更新測試斷言以驗證新的錯誤結構與狀態碼。

- 2025-10-10 `feat(router): 根路徑導向 Swagger UI`
  - `main.go`: 新增 `/` 路由自動重新導向至 Swagger 與引入 `net/http`。
  - `README.md`: 說明根路徑將自動開啟 Swagger 介面。
  - `middlewares/jwt_middleware.go`: 放行根路徑讓未登入使用者可轉址至 Swagger。

- 2025-10-10 `feat(database): 自動建庫並種子預設資料`
  - `database/db.go`: 拆分 DSN 建構並於缺少資料庫時自動建立，同步掛載權限與平台資料表初始化。
  - `database/User_DB/User_Table.go`: 建立預設系統管理員帳號並調整欄位類型與匯入。
  - `database/Permission_DB/*`、`database/Platform_DB/*`: 新增權限與平台資料表 schema 與預設資料。
  - `app/controller/v1/user/User_Register.go`、`app/models/account/User_Model.go`: 對應欄位使用整數型別並連結新權限 id。
  - `docs/router_flow.md`: 補充路由流程與登入時序圖說明。
  - `AGENTS.md`、`GEMINI.md`: 更新協作指引與軟體架構說明。
  - `README.md`: 增加 Docker 啟動 PostgreSQL 範例指令。
