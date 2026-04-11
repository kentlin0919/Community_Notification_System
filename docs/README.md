# 開發文件索引

此目錄依主題分類專案文件，方便快速尋找對應資訊。Swagger 相關產出維持在目錄根層，請勿手動修改。

## 文件分類

### `architecture/`
- `router_flow.md`：說明路由註冊、middleware 鏈與登入請求流向。

### `analysis/`
- `project_analysis.md`：專案整體分析，包含架構盤點、資料模型、整體 activity diagram、整體 sequence diagram 與 class diagram。

### `features/`
- `README.md`：功能文件索引，說明功能資料夾命名與維護規範。
- `user_login/README.md`：使用者登入功能文件。
- `user_register/README.md`：使用者註冊功能文件。
- `user_delete/README.md`：使用者刪除功能文件。
- `community_management/`：社區管理大項功能資料夾。
- `community_management/community_getlist/README.md`：社區列表查詢功能文件。
- `community_management/community_register/README.md`：新增社區功能文件。
- `platform_getlist/README.md`：平台列表查詢功能文件。
- `message_send/README.md`：發送通知功能文件。
- `permission_management/README.md`：權限管理共用規格，定義 `Super admin`、`Admin` 與 `PermissionID 3 ~ 7` 的授權規則。
- `facility_booking/`：設施預約大項功能資料夾。
- `facility_booking/facility_create/README.md`：新增預約設施功能文件，聚焦設施主檔建立流程與驗證規則。
- `facility_booking/facility_reservation/README.md`：社區基本設施預約設計文件，包含流程圖、功能拆分、類別圖與資料表草案。
- `facility_booking/reservation_reschedule_request/README.md`：申請更改預約時間功能文件，說明改期流程、驗證規則與資料設計建議。
- `facility_booking/reservation_cancel/README.md`：取消預約功能文件，說明取消規則、違規處理與通知流程。

### `commit_summaries/`
- `commit_summary_2025_11.md`：2025 年 11 月 commit 摘要，依新到舊排序。
- `commit_summary_2025_10.md`：2025 年 10 月 commit 摘要，依新到舊排序。
- `commit_summary_2025_09.md`：2025 年 9 月 commit 摘要，依新到舊排序。

## Swagger 產出
- `docs.go`
- `swagger.json`
- `swagger.yaml`

以上三者由 `swag init -g main.go` 產生，不應手動編輯。

## 文件維護原則
- 架構分析放在 `analysis/`
- 功能說明放在 `features/`，並以一個功能一個資料夾的方式管理
- 路由或基礎設施流程放在 `architecture/`
- commit 歷史整理放在 `commit_summaries/`
- Swagger 相關檔案維持於 `docs/` 根目錄
- 涉及多租戶或社區範圍的規則，應同步更新主 README、架構文件與功能文件，確保 `community_id` 隔離規則描述一致
