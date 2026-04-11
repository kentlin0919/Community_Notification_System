# Commit 摘要 - 2026年04月

| SHA | 作者 | 日期 | 訊息 |
| --- | --- | --- | --- |
| 595ac16 | kentlin0919 | 2026-04-11 | docs: add project analysis and feature planning documents |
| 587abd5 | kentlin0919 | 2026-04-11 | fix(firebase): handle missing credentials and fix send message crash |
| 4b7e8d9 | kentlin0919 | 2026-04-11 | chore(docker): setup docker-compose and vscode remote debug |
| (Pending) | Antigravity | 2026-04-09 | fix(docs): 修正 Markdown 標題空行問題 (MD022) |

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
