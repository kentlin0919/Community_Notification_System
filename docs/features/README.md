# 功能文件索引

本目錄依功能拆分，每個功能使用獨立資料夾管理，方便後續持續補充流程圖、欄位說明、測試案例與設計決策。

## 目錄結構
- `user_login/`：使用者登入
- `user_register/`：使用者註冊
- `user_delete/`：使用者刪除
- `community_management/`：社區管理大項功能資料夾
- `community_management/community_getlist/`：社區列表查詢
- `community_management/community_register/`：新增社區
- `platform_getlist/`：平台列表查詢
- `message_send/`：發送通知
- `parcel_management/`：包裹管理（代收與領取）
- `permission_management/`：權限管理共用規格
- `facility_booking/`：設施預約大項功能資料夾
- `facility_booking/facility_create/`：新增預約設施
- `facility_booking/facility_reservation/`：社區基本設施預約設計文件
- `facility_booking/reservation_reschedule_request/`：申請更改預約時間
- `facility_booking/reservation_cancel/`：取消預約

## 文件規範
- 每個功能一個資料夾
- 每個資料夾至少包含一份 `README.md`
- 內容應包含：功能目標、API、涉及模組、Activity Diagram、Sequence Diagram、限制與備註
