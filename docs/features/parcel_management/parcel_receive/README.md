# 管理員代收包裹流程 (Parcel Receive)

## 功能說明
當物業管理員（保全）收到外部快遞公司送達的包裹時，需登錄包裹資訊。系統將自動根據住戶資訊，向該戶所有成員發送 FCM 推播通知。

## 角色視角：管理員
1.  **輸入門牌**：確認包裹上的地址，在系統輸入門牌號碼。
2.  **載入成員**：系統自動根據門牌（家庭群組）載入該戶所有註冊成員。
3.  **選擇收件人**：透過下拉選單選擇包裹所屬的成員（或選擇「全家」）。
4.  **登錄**：輸入物流公司、單號與備註。
5.  **確認**：提交後，系統針對特定成員（或全家）發送推播通知。

## 流程圖

### 活動圖 (Activity Diagram)

```mermaid
stateDiagram-v2
    state "管理員操作" as AdminPart {
        A[收到實體包裹] --> B[輸入門牌號碼]
        B --> C: 自動載入家庭成員清單
        C --> D: 從下拉選單選擇特定收件人
        D --> E: 輸入物流資訊 (公司、單號)
        E --> F: 提交登錄
    }

    state "系統後台" as SystemPart {
        F --> G: 驗證並寫入資料庫 (ParcelInfo)
        G --> H: 狀態標記為 "待領取"
        H --> I: 根據選擇對象獲取 FcmToken
    }

    state "通知服務" as NotificationPart {
        I --> J: 發送精準通知 (FCM)
        J --> K: 住戶手機顯示個人化通知
    }
```

### 時序圖 (Sequence Diagram)

```mermaid
sequenceDiagram
    participant Admin as 管理員 (Web/App)
    participant API as API Server (Gin)
    participant DB as Database (PostgreSQL)
    participant FCM as Firebase (FCM)
    participant User as 住戶裝置

    Admin->>API: GET /api/v1/homes/:home_id/members (查詢住戶成員)
    API->>DB: SELECT id, name FROM user_infos WHERE home_id = ?
    DB-->>API: 返回成員清單
    API-->>Admin: 顯示下拉選單

    Admin->>API: POST /api/v1/parcels (Token, Courier, HomeID, TargetUserID)
    Note over API: 解析 CommunityID (Middleware)
    API->>DB: INSERT INTO parcel_infos (TargetUserID, Status: 1)
    DB-->>API: Success (ParcelID)
    
    alt 指定特定成員
        API->>DB: SELECT fcmtoken FROM user_infos WHERE id = TargetUserID
    else 指定全家 (預設)
        API->>DB: SELECT fcmtoken FROM user_infos WHERE home_id = HomeID
    end
    DB-->>API: 取得代送 Token 清單
    
    loop 對每個 Token
        API->>FCM: 發送推播 (標題: 包裹到府, 內容: [姓名] 您有包裹已由管理室代收)
        FCM-->>User: 顯示通知
    end
    
    API-->>Admin: 201 Created (已成功登錄並發送通知)
```

## API 規格概要

### 1. 登錄包裹
- **URL**: `POST /api/v1/parcels`
- **權限**: 管理員 (Role Level >= 2)
- **Request Body**:
    ```json
    {
      "home_id": 101,
      "courier_company": "黑貓宅急便",
      "tracking_number": "9876543210",
      "remark": "大型件"
    }
    ```
- **Response**: `201 Created`

## 異常處理
- **找不到住戶**：若輸入的 `home_id` 不存在，返回 `404 Not Found`。
- **重複單號**：若同一社區已有相同單號且處於「待領取」狀態，提示管理員確認。
