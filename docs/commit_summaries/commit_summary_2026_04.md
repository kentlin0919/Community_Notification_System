# Commit 摘要 - 2026年04月

| SHA | 作者 | 日期 | 訊息 |
| --- | --- | --- | --- |
| (Session) | Antigravity | 2026-04-20 | doc: 維護專案文件 (README, GEMINI, Analysis, Router Flow) 以符合實作現況 |
| 1176630 | Gemini CLI | 2026-04-17 | feat(community): implement community registration audit, permission profiles, and facility reservation system |
| 595ac16 | kentlin0919 | 2026-04-11 | docs: add project analysis and feature planning documents |
| 587abd5 | kentlin0919 | 2026-04-11 | fix(firebase): handle missing credentials and fix send message crash |
| 4b7e8d9 | kentlin0919 | 2026-04-11 | chore(docker): setup docker-compose and vscode remote debug |
| (Pending) | Antigravity | 2026-04-09 | fix(docs): 修正 Markdown 標題空行問題 (MD022) |

## 逐行分析 - doc: 維護專案文件 (README, GEMINI, Analysis, Router Flow) (Session)

### 變更檔案：
- README.md
- GEMINI.md
- docs/analysis/project_analysis.md
- docs/architecture/router_flow.md
- docs/commit_summaries/commit_summary_2026_04.md

### 變更描述：
同步功能開發現況至各項文件，包含檔案結構、核心功能時序圖、API 總覽與成熟度觀察。

### 分析細節：
1. **README.md**:
   - 加入設施預約、社區審核、預約改期與權限 Profile 的 Mermaid 時序圖。
   - 修正檔案結構樹與核心功能描述。
2. **GEMINI.md**:
   - 修正檔案結構樹，移除「預留」標籤，同步最新目錄結構。
3. **project_analysis.md**:
   - 更新 API 總覽表。
   - 擴充類別圖 (Class Diagram) 以納入設施與預約模型。
   - 更新「功能成熟度觀察」，確認社區審核與設施系統已實作完成。
4. **router_flow.md**:
   - 更新中介層列表（加入 Permission 與 CommunityContext 中介層）。
   - 更新路由表，移除過時的範例。

## 逐行分析 - feat(community): implement community registration audit, permission profiles, and facility reservation system (1176630)

### 變更檔案：
- app/controller/v1/communityManager/ (CommunityManager_Approve.go, CommunityManager_Reject.go, CommunityManager_Controller.go, CommunityManager_add.go)
- app/controller/v1/facility/ (Facility_Controller.go, Facility_Create.go)
- app/controller/v1/permission/ (Permission_Controller.go, Permission_Update.go)
- app/controller/v1/reservation/ (Reservation_Controller.go, Reservation_Create.go, Reservation_Cancel.go, Reservation_Reschedule.go)
- app/models/ (community, facility, permission, reservation models)
- app/repositories/ (community, facility, permission, reservation repositories)
- database/ (CommunityRegisterApplication_Schema.go, CommunityPermissionProfile_Schema.go, Facility_DB, etc.)
- middlewares/ (community_context_middleware.go, jwt_middleware.go, permission_middleware.go)
- routers/api/v1/v1.go
- utils/Jwt_Token.go
- docs/ (swagger files and implementation docs)

### 變更描述：
實現完整的社區申請審核流程、社區專屬權限角色名稱定義，以及公用設施預約管理系統，並擴充 JWT 以支援多租戶隔離。

### 分析細節：
1. **社區管理 (Community Management)**:
   - 重構 `CommunityManager_Register` 為申請單模式，狀態預設為 `pending`。
   - 新增 `Approve` 與 `Reject` API，由 Super Admin 審核申請並自動建立正式社區與首位管理員。
2. **權限管理 (Permission Management)**:
   - 支援社區管理員自定義角色名稱（針對等級 3-7），並新增對應的查詢與更新 API。
   - 擴充 `jwt_middleware.go` 以從 Token 提取 `user_id`, `permission_id` 與 `community_id` 並存入 Context。
3. **設施與預約 (Facility & Reservation)**:
   - 實作設施主檔建立功能（包含封面圖、位置、狀態等）。
   - 實作預約系統核心邏輯：新增預約、取消預約、變更時段申請以及管理員核准變更時段。
4. **JWT 與 安全 (Security)**:
   - 擴充 `GenerateJWT` 函數，將關鍵身分與租戶資訊（CommunityID）加密至 Token 中，強化多租戶數據隔離。
5. **API 文件與開發指南**:
   - 更新 Swagger 定義以符合最新 API 規格。
   - 針對新功能撰寫詳細的實作說明文件 (`Implementation_Doc.md`)。

## 逐行分析 - docs: add project analysis and feature planning documents (595ac16)

### 變更檔案：
- README.md
- docs/README.md
- docs/architecture/router_flow.md
- docs/analysis/project_analysis.md
- docs/features/ (多個文件)

### 變更描述：
新增專案分析、架構流程、權限管理與設施預約等多項核心功能的規劃文件。

### 分析細節：
1. **README.md**:
   - 更新技術棧資訊（Go 1.26.1, Docker debug 流程）。
   - 擴充系統模組概觀與核心功能時序圖。
   - 新增文件索引與 Docker Debug 指引。
2. **docs/architecture/router_flow.md**:
   - 補充 `CommunityContextMiddleware` 的規劃，說明社區上下文隔離規則。
3. **docs/analysis/project_analysis.md**:
   - 提供專案整體 Activity, Sequence 與 Class Diagram。
4. **docs/features/**:
   - 建立使用者、社區管理、設施預約、發送通知等細部功能規格書。

## 逐行分析 - fix(firebase): handle missing credentials and fix send message crash (587abd5)

### 變更檔案：
- app/controller/v1/message/Message_SendMessage.go
- configs/config.go
- pkg/firebase/firebase.go
- .gitignore

### 變更描述：
提升系統對於遺失設定檔（.env, Firebase key）的容錯能力，並修正 Firebase 未初始化時發送訊息導致的當機。

### 分析細節：
1. **pkg/firebase/firebase.go**:
   - 啟動時檢查 `serviceAccountKey.json`，若遺失則記錄日誌而非 Fatal。
2. **app/controller/v1/message/Message_SendMessage.go**:
   - 在 `SendMessage` 中加入 `firebase.FcmClient` 的 nil 檢查，回傳 `503 Service Unavailable`。
3. **configs/config.go**:
   - 當 `.env` 遺失時僅顯示提示日誌，允許使用系統環境變數啟動。
4. **.gitignore**:
   - 排除 `tmp/`, `firebase-debug.log`, `main` 與 `*.log`。

## 逐行分析 - chore(docker): setup docker-compose and vscode remote debug (4b7e8d9)

### 變更檔案：
- .air-debug.toml / .air.toml
- .vscode/launch.json / .vscode/tasks.json
- Dockerfile / docker-compose.yml

### 變更描述：
整合 Docker 容器化開發環境，並配置 VS Code 與容器內 Delve 的遠端偵錯連結。

### 分析細節：
1. **Dockerfile / docker-compose.yml**:
   - 基於 Go 1.26.1 建立開發映像，安裝 Air 與 Delve。
   - 建立 API 服務與 PostgreSQL 容器關聯。
2. **.vscode/**:
   - 配置 `Docker: Remote Debug` launch 設定，透過 40000 埠連結容器。
   - 新增 `docker-compose: up/down/wait-dlv` 自動化工作流。
3. **Air 配置**:
   - 更新 `.air-debug.toml` 的 `full_bin` 以正確調用 Delve 進行偵錯。
