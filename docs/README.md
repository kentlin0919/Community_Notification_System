# 開發文件索引

此目錄依主題分類專案文件，方便快速尋找對應資訊。Swagger 相關產出維持在目錄根層，請勿手動修改。

## 文件分類

### `architecture/`
- `router_flow.md`：說明路由註冊、middleware 鏈與登入請求流向。

### `skills/`
- [Skills 總覽](skills/README.md)：定義專案中的 AI 代理人工作流與標準化任務指令集。
- [自動化審核](skills/skill_auto_approval.md)：定義社區申請單的自動檢核與決策邏輯。
- [系統維護](skills/skill_maintenance_cron.md)：定義清理過期數據與系統檢查的維護任務。
- [訊息推播](skills/skill_notification_dispatch.md)：標準化訊息分發與 Firebase FCM 整合工作流。

### `analysis/`
- `project_analysis.md`：專案整體分析，包含架構盤點 ... (含 IOT 未來規劃)。
- `PRD_Family_Group.md`：家庭群組與包裹通知優化規格。

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
- `parcel_management/README.md`：包裹管理功能文件總覽，涵蓋代收錄入與住戶領取。
- `parcel_management/parcel_receive/README.md`：管理員代收包裹流程。
- `parcel_management/parcel_pickup/README.md`：住戶領取包裹流程。
- `permission_management/README.md`：權限管理共用規格，定義 `Super admin`、`Admin` 與 `PermissionID 3 ~ 7` 的授權規則。
- `facility_booking/`：設施預約大項功能資料夾。
- `facility_booking/facility_create/README.md`：新增預約設施功能文件，聚焦設施主檔建立流程與驗證規則。
- `facility_booking/facility_reservation/README.md`：社區基本設施預約設計文件，包含流程圖、功能拆分、類別圖與資料表草案。
- `facility_booking/detailed_design_plan.md`：預約系統詳細設計規劃，包含時段衝突檢查、容量控管、違規停權與核銷流程。
- `facility_booking/ui_specification.md`：預約系統 UI/UX 詳細規劃，定義住戶端與管理端各頁面的欄位、組件與互動邏輯。
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
