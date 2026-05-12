# 權限與資料結構 (Permission & Data Structure)

> [!NOTE]
> 本文件描述後端內部的權限校驗邏輯與資料結構，屬於系統層面之實作細節，不對應至 App/Desktop 之功能分組。

## 1. 權限設計
系統採用 RBAC (Role-Based Access Control) 與社區隔離機制：
- **Super Admin**: 擁有 `*` 權限，跨社區存取。
- **Community Admin**: 擁有該 `community_id` 下的所有權限。
- **Resident**: 擁有特定功能（如包裹領取、設施預約）的 `view` 與 `create` 權限。

## 2. 資料結構
- **Permission Table**: 定義功能關鍵字 (e.g., `parcel:read`, `parcel:write`).
- **Role Table**: 關聯權限集合。
- **UserRole Table**: 關聯使用者、角色與社區 ID。

## 3. 中間件校驗
後端 `middlewares/` 下的 `PermissionMiddleware` 會在 API 進入點進行校驗。
