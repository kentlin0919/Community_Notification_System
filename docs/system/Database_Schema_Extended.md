# 智慧社區擴展資料庫設計 (Database Schema Extended)

> [!NOTE]
> **重要說明**：本文件中定義的資料表目前僅供未來 2.0 版本設計參考，尚未正式實裝至系統中。

本文件定義了未來達成「智慧社區 2.0」所需的擴展資料表。這些設計遵循 GORM 規範，並考慮了與現有 `User_DB` 與 `Community_DB` 的關聯。

## 1. 智慧通行模組 (Smart Access)

### Visitor_Passes (訪客邀請碼表)
用於記錄住戶發送給訪客的臨時通行金鑰。

| 欄位名 | 類型 | 描述 |
| :--- | :--- | :--- |
| **ID** | uint (PK) | 邀請碼唯一識別碼 |
| **CommunityID** | uint (FK) | 所屬社區 |
| **OwnerID** | uint (FK) | 發起住戶 (User.ID) |
| **VisitorInfo** | string | 訪客姓名/電話備註 |
| **AccessKey** | string (Index) | 加密後的 QRCode 原始內容 |
| **StartTime** | datetime | 有效開始時間 |
| **EndTime** | datetime | 有效結束時間 |
| **IsUsed** | bool | 是否已使用 |
| **MaxUsage** | int | 允許進入次數 (1 代表一次性，>1 代表多次可用) |

### Access_Logs (通行紀錄表)
追蹤住戶與訪客的開門紀錄。

| 欄位名 | 類型 | 描述 |
| :--- | :--- | :--- |
| **ID** | uint (PK) | 紀錄編號 |
| **TargetID** | uint | 開門者 ID (可能是 User 或 Visitor_Pass) |
| **TargetType** | string | 類型: `RESIDENT`, `VISITOR`, `STAFF` |
| **GateID** | string | 設備 ID (如: 大門口機 01, 地下室電梯 02) |
| **Status** | int | 0: 成功, 1: 拒絕, 2: 逾期 |
| **Timestamp** | datetime | 通行時間 |

---

## 2. 物業維修模組 (Maintenance & Asset)

### Maintenance_Tickets (報修工作單)
擴展現有的報修流程，加入派工與感應標籤。

| 欄位名 | 類型 | 描述 |
| :--- | :--- | :--- |
| **ID** | uint (PK) | 單號 |
| **ReporterID** | uint | 報修人 ID |
| **AssigneeID** | uint | 維修負責人員 ID (員工) |
| **AssetTag** | string | 損壞設備的 NFC/QRCode 標籤 ID |
| **Status** | int | 0-待受理, 1-處理中, 2-已完工, 3-驗收失敗 |
| **RepairPhoto** | string | 修復後的照片 URL |
| **MaterialCost** | decimal | 耗材費用紀錄 |

### Asset_Maintenance_Logs (設備維修履歷)
用於追蹤單一設備的完整生命週期。

| 欄位名 | 類型 | 描述 |
| :--- | :--- | :--- |
| **ID** | uint (PK) | 紀錄編號 |
| **AssetID** | uint (FK) | 關聯資產 |
| **Action** | string | 維養動作 (e.g., 更換電池, 定期保養) |
| **Technician** | string | 技術人員姓名 |
| **Result** | string | 檢測結果 |

---

## 3. 智慧能耗管理 (Smart Utility)

### Utility_Readings (度數申報表)
用於住戶每月自行申報或智慧電表自動上傳。

| 欄位名 | 類型 | 描述 |
| :--- | :--- | :--- |
| **ID** | uint (PK) | 紀錄編號 |
| **UserID** | uint | 住戶 ID |
| **Type** | int | 1-瓦斯, 2-水, 3-電 |
| **PreviousValue** | decimal | 上期度數 |
| **CurrentValue** | decimal | 本期度數 |
| **PhotoURL** | string | 度數表拍照證明 (非自動化設備時使用) |
| **BillStatus** | int | 0-未生成帳單, 1-已出帳, 2-已繳清 |

### EV_Charging_Logs (EV 充電樁負載管理)
用於 ESG 能耗追蹤與電力負載平衡。

| 欄位名 | 類型 | 描述 |
| :--- | :--- | :--- |
| **ID** | uint (PK) | 紀錄編號 |
| **StationID** | string | 充電樁編號 |
| **ResidentID** | uint | 使用住戶 |
| **PowerConsumption** | decimal | 消耗電量 (kWh) |
| **Duration** | int | 充電時數 (分鐘) |
| **LoadLimitApplied** | bool | 是否曾因負載平衡被限流 |

## 4. 設計規範
1. **Soft Delete**: 所有擴展表必須符合 `gorm.Model`，具備 `deleted_at` 以支援軟刪除。
2. **Audit Logging**: 關鍵變更（如訪客授權）必須同步寫入 `UserLog_DB`。
3. **Data Localization**: 對於高頻率存取的 `AccessKey`，應同時考慮 Redis 暫存設計。
