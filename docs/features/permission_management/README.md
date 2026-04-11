# 權限管理共用規格

## 功能目標
定義系統共用的權限模型、角色階層、社區範圍與授權判斷規則，作為 JWT、middleware、controller、repository 與後台角色管理的統一依據。

## 設計原則
- 權限大小一律以 `PermissionID` 判斷
- `PermissionID` 數字越小，權限越高
- `Super admin` 為全平台角色，可管理所有社區
- `Admin` 為社區內角色，只能管理所屬社區
- `3 ~ 7` 為社區內可自訂角色名稱，但權限高低仍由 ID 決定，不由名稱決定

## 權限階層
| PermissionID | 系統角色 | 角色範圍 | 說明 |
| --- | --- | --- | --- |
| `1` | `Super admin` | 全平台 | 可管理所有社區、所有社區 admin、平台權限設定 |
| `2` | `Admin` | 單一社區 | 社區最高管理角色，可管理自己社區的帳號、角色與資料 |
| `3` | 自訂角色 A | 單一社區 | 由社區 `Admin` 自行命名，例如保全、主任、櫃台 |
| `4` | 自訂角色 B | 單一社區 | 由社區 `Admin` 自行命名 |
| `5` | 自訂角色 C | 單一社區 | 由社區 `Admin` 自行命名 |
| `6` | 自訂角色 D | 單一社區 | 由社區 `Admin` 自行命名 |
| `7` | 自訂角色 E | 單一社區 | 由社區 `Admin` 自行命名 |

## 核心規則

### 1. `Super admin`
- `PermissionID = 1`
- 可管理所有社區
- 可審核社區申請
- 可核可後建立新社區
- 可建立或重設任一社區的 `Admin`
- 可查看所有社區資料
- 可跨社區執行平台層操作

### 2. `Admin`
- `PermissionID = 2`
- 僅能管理自己的社區
- 可定義自己社區的 `PermissionID 3 ~ 7` 角色名稱與用途
- 可建立自己社區內的使用者帳號
- 可指派自己社區內的角色
- 不可管理其他社區
- 不可建立或修改 `Super admin`

### 3. 自訂角色
- `PermissionID = 3 ~ 7`
- 角色名稱可由社區 `Admin` 自訂
- 例如：
  - `3 = 保全`
  - `4 = 主委`
  - `5 = 委員`
  - `6 = 櫃台`
  - `7 = 財務`
- 雖然名稱可改，但授權比較只能看 ID：
  - `3` 權限大於 `4`
  - `4` 權限大於 `5`
  - `5` 權限大於 `6`
  - `6` 權限大於 `7`

## 授權判斷規則

### 1. 角色大小判斷
統一規則：

```text
PermissionID 數值越小，權限越高
1 > 2 > 3 > 4 > 5 > 6 > 7
```

建議實作方式：
- `operator.PermissionID <= target.PermissionID`
- 或依操作類型額外限制

### 2. 社區範圍判斷
- `PermissionID = 1` 可跨社區操作
- `PermissionID = 2 ~ 7` 僅能操作相同 `community_id` 的資料
- 若操作者與目標資料 `community_id` 不一致，應直接拒絕

### 3. 禁止只靠角色名稱判斷
不可使用以下方式判權：
- `role_name == "保全"`
- `role_name == "主委"`
- `role_name == "委員"`

原因：
- `3 ~ 7` 名稱可被各社區自行定義
- 同樣名稱在不同社區可能代表不同實務職責
- 系統授權必須保持可預測，因此只能根據 `PermissionID`

## 社區角色自訂規則

### 可自訂範圍
- 僅 `Admin` 可修改自己社區的 `PermissionID 3 ~ 7`
- 可自訂內容包含：
  - 顯示名稱
  - 描述
  - 後台標籤
  - 權限說明

### 不可自訂範圍
- 不可修改 `PermissionID = 1`
- 不可修改 `PermissionID = 2`
- 不可新增 `PermissionID > 7`
- 不可改變數字與權限排序

### 固定規則
- `1` 永遠是 `Super admin`
- `2` 永遠是 `Admin`
- 最大只到 `7`

## 與社區建立流程的關係
- `社區申請人` 先送出社區申請
- `Super admin` 審核核可後，系統才正式建立該社區的第一個 `Admin`
- 這個初始社區管理員的 `PermissionID` 必須固定為 `2`
- 建立完成後，該 `Admin` 才能去設定自己社區的 `3 ~ 7` 角色名稱

## 與 JWT / Middleware 的關係

### JWT 建議 claims
- `user_id`
- `permission_id`
- `community_id`
- `email`

### Middleware 建議
- `JWTAuthMiddleware`：解析使用者身份
- `CommunityContextMiddleware`：建立 `community_id` 上下文
- `PermissionMiddleware`：依 `permission_id` 判斷是否允許操作

## 建議資料設計

### 1. 系統固定權限表
用途：定義固定的權限 ID 規則。

建議欄位：
- `permission_id`
- `system_key`
- `default_name`
- `is_system_fixed`

### 2. 社區角色顯示設定表
用途：讓每個社區自訂 `3 ~ 7` 的顯示名稱，但不改變權限大小。

建議表名：
- `community_permission_profile`

建議欄位：
- `id`
- `community_id`
- `permission_id`
- `display_name`
- `description`
- `updated_by`
- `updated_at`

限制：
- `permission_id` 僅允許 `3`、`4`、`5`、`6`、`7`
- 同一社區同一 `permission_id` 僅一筆設定

## 使用範例

### 範例 1：Super admin 建立社區
- 操作者 `permission_id = 1`
- 允許審核社區申請並核可建立新社區
- 允許建立該社區第一個 `Admin`

### 範例 2：Admin 建立保全帳號
- 操作者 `permission_id = 2`
- 目標 `permission_id = 3`
- 同社區下允許

### 範例 3：角色 4 想修改角色 3
- 操作者 `permission_id = 4`
- 目標 `permission_id = 3`
- 因 `4 > 3`，不允許

### 範例 4：Admin 修改角色名稱
- 操作者 `permission_id = 2`
- 可修改自己社區的 `permission_id = 3 ~ 7` 顯示名稱
- 不可修改 `1` 與 `2`

## 實作備註
- 目前資料庫種子已存在 `1 ~ 6` 權限資料，但後續若採用本規格，建議重新整理成以 `1 ~ 7` 為核心。
- 若現有 `6 = 一般住戶` 已有歷史資料，需再評估是否沿用為社區自訂角色，或補 migration 將既有資料對齊新規格。
- 正式實作前，建議先決定：
  - `7` 是否就是最低權限的一般住戶
  - 還是 `3 ~ 7` 全部都是社區 admin 自訂的管理層角色

## 關聯文件
- [新增社區](/Users/kent/project/Community_Notification_System/docs/features/community_management/community_register/README.md)
- [社區基本設施預約](/Users/kent/project/Community_Notification_System/docs/features/facility_booking/facility_reservation/README.md)
- [專案分析文件](/Users/kent/project/Community_Notification_System/docs/analysis/project_analysis.md)
