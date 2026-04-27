# 管理員代收包裹流程 (Parcel Receive)

## 功能說明
當物業管理員（保全）收到外部快遞公司送達的包裹時，需登錄包裹資訊。系統將自動根據住戶資訊，向該戶所有成員發送 FCM 推播通知。

## 角色視角：管理員
1.  **收貨**：確認包裹上的收件人地址（樓層、房號）。
2.  **登錄**：開啟管理後台，輸入物流公司、單號與選擇住戶。
3.  **確認**：系統顯示「登錄成功並已發送通知」。

## 流程圖

### 活動圖 (Activity Diagram)

```mermaid
stateDiagram-v2
    state "管理員操作" as AdminPart {
        A[收到實體包裹] --> B[於系統輸入物流資訊 (公司、單號)]
        B --> C: 選擇所屬住戶
        C --> D: 提交登錄
    }

    state "系統後台" as SystemPart {
        D --> E: 驗證資訊並寫入資料庫 (ParcelInfo)
        E --> F: 狀態標記為 "待領取"
        F --> G: 查詢該住戶所有裝置的 FcmToken
    }

    state "通知服務" as NotificationPart {
        G --> H: 發送包裹通知 (FCM)
        H --> I: 住戶手機顯示通知
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

    Admin->>API: POST /api/v1/parcels (Token, Courier, TrackingNo, HomeID)
    Note over API: 解析 CommunityID (Middleware)
    API->>DB: INSERT INTO parcel_infos (Status: 1)
    DB-->>API: Success (ParcelID)
    
    API->>DB: SELECT fcm_tokens FROM user_infos WHERE home_id = ?
    DB-->>API: List of Tokens
    
    loop 對每個 Token
        API->>FCM: 發送推播 (標題: 包裹到府, 內容: 您的包裹已由管理室代收)
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
