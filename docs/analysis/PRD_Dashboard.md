# Dashboard Design Specification (首頁儀表板設計規格)

## 1. 概述
本文件定義「社區通知系統」首頁儀表板的資訊架構與資料聚合邏輯。目標是根據使用者的權限等級 (Permission Level)，提供即時、相關且具備行動導向的數據摘要。

## 2. 角色與資訊權限映射
系統目前將使用者分為三種核心首頁視圖：

| 角色 | 權限等級 (Level) | 核心需求 | 呈現資訊 |
| :--- | :--- | :--- | :--- |
| **超級管理員** | 1-2 | 監控全平台運作狀況 | 全平台社區/住戶總數、待審核社區申請單 |
| **社區管理員** | 3-7 | 管理社區日常庶務 | 社區住戶總數、今日設施預約數、待處理改期單 |
| **住戶** | 8-9 | 個人行程與重要公告 | 個人即將到來的預約、社區最新公告 |

---

## 3. 資訊架構 (Information Architecture)

### 3.1 住戶 (Resident)
*   **即將到來的預約 (Upcoming Reservations)**: 顯示使用者未來三日內已核准的設施使用時段。
*   **最新公告 (Recent Messages)**: 顯示所屬社區發布的最新 5 則公告，確保重要資訊不遺漏。

## 角色資訊架構與 API 規範 (v1/home)

根據身份權限 (PermissionID) 返回對應的儀表板數據。

### 1. 住戶 (PermissionID: 8-9)
**核心價值**：自我服務、即時通知、設施預約。

#### JSON 欄位定義
| 欄位名 | 型態 | 範例值 | 說明 |
| :--- | :--- | :--- | :--- |
| `upcoming_reservations` | `Array` | `[...]` | 最近 3 筆未來的設施預約 |
| `recent_messages` | `Array` | `[...]` | 最近 5 筆社區公告或個人訊息 |
| `maintenance_status` | `Object` | `{...}` | 若有報修，顯示進度摘要 |

#### JSON 範例
```json
{
  "role": "resident",
  "data": {
    "upcoming_reservations": [
      {
        "id": 101,
        "facility_name": "健身房",
        "date": "2026-04-21",
        "start_time": "18:00",
        "end_time": "19:00",
        "status": "confirmed"
      }
    ],
    "recent_messages": [
      {
        "id": 2001,
        "title": "電梯定期保養通知",
        "content": "本週三 A 棟電梯將進行保養...",
        "create_time": "2026-04-20 10:00",
        "category": "announcement"
      }
    ]
  }
}
```

---

### 2. 社區管理員 (PermissionID: 3-7)
**核心價值**：營運效率、異常監控、申請審核。

#### JSON 欄位定義
| 欄位名 | 型態 | 範例值 | 說明 |
| :--- | :--- | :--- | :--- |
| `stats` | `Object` | `{...}` | 包含：今日預約數、待審核申請數、住戶總數 |
| `urgent_notices` | `Array` | `[...]` | 需要立即處理的事項（如：逾期報修） |
| `recent_activities` | `Array` | `[...]` | 最近發生的操作日誌摘要 |

#### JSON 範例
```json
{
  "role": "community_admin",
  "data": {
    "stats": {
      "resident_count": 150,
      "pending_applications": 5,
      "today_reservations": 12
    },
    "urgent_notices": [
      {
        "id": 501,
        "type": "application",
        "description": "林先生 申請加入社區",
        "time": "2 小時前"
      }
    ]
  }
}
```

---

### 3. 超級管理員 (PermissionID: 1-2)

**核心價值**：全平台監控、系統健康度、商務統計。

#### JSON 欄位定義

| 欄位名 | 型態 | 範例值 | 說明 |
| :--- | :--- | :--- | :--- |
| `platform_stats` | `Object` | `{...}` | 總社區數、總用戶數、今日推播總量 |
| `pending_communities` | `Array` | `[...]` | 待審核的社區註冊申請單 |

---

## 核心功能活動流程 (Activity Diagrams)

為了銜接前後端開發，以下將核心功能拆分為「UI 操作流」與「後端業務流」。

### 1. 住戶端：快速預約流程 (Quick Booking)

#### [UI 操作] 活動圖

```mermaid
flowchart TD
    Start((開始)) --> List[首頁預約快捷區列出常用設施]
    List --> Click[住戶點擊快速預約圖示]
    Click --> Popup[彈出時段選取半屏視窗]
    Popup --> Cancel{使用者取消?}
    Cancel -- Yes --> Close[關閉視窗，返回首頁]
    Cancel -- No --> Time[選擇時間並點擊確認預約]
    Time --> Loading[按鈕進入 Loading 狀態並禁用重複點擊]
    Loading --> API{API 回傳成功?}
    API -- Success --> Anim[顯示成功動畫 Lottie]
    Anim --> Refresh[首頁清單更新數據]
    API -- Failure --> Error[顯示紅字錯誤訊息]
    Error --> Reset[恢復按鈕狀態]
    Refresh --> End((結束))
    Close --> End
    Reset --> End
```

#### [後端處理] 活動圖

```mermaid
flowchart TD
    Start((接收 POST /quick 請求)) --> Auth[JWT Middleware 驗證 Token 合法性]
    Auth --> Context[從 Context 提取 UserID 與 CommunityID]
    Context --> Priv{權限等級 >= 8?}
    Priv -- No --> F403[返回 403 Forbidden]
    Priv -- Yes --> DBCheck[查詢指定設施是否存在且啟用]
    DBCheck --> Trans[執行資料庫事務 Begin Transaction]
    Trans --> CapCheck{名額足夠?}
    CapCheck -- No --> Rollback[Rollback]
    Rollback --> F409[返回 409 Conflict]
    CapCheck -- Yes --> Insert[INSERT 預約記錄至 DB]
    Insert --> Notify[發送非同步通知給管理員]
    Notify --> Commit[Commit]
    Commit --> F201[返回 201 Created]
    F201 --> End((結束))
    F403 --> End
    F409 --> End
```

---

### 2. 管理員端：緊急公告與推播 (Instant Announcement)

#### [UI 操作] 活動圖
```mermaid
flowchart TD
    Start((開始)) --> Click[點擊緊急公告按鈕]
    Click --> Edit[開啟全螢幕編輯模式]
    Edit --> Input[輸入標題與簡短內文]
    Input --> Check[勾選同步發送推播通知]
    Check --> Preview[點擊預覽與發送]
    Preview --> Confirm{確認發送?}
    Confirm -- 否 --> Back[保留草稿並返回]
    Confirm -- 是 --> Progress[顯示發送中進度條]
    Progress --> Finish[顯示發送成功 Toast 提醒]
    Finish --> Home[導向回首頁]
    Home --> End((結束))
    Back --> End
```

#### [後端處理] 活動圖
```mermaid
flowchart TD
    Start((接收 POST /emergency 請求)) --> Auth[驗證管理員權限 Level 1-7]
    Auth --> Safe[過濾敏感字元與格式驗證]
    Safe --> Valid{驗證通過?}
    Valid -- No --> F400[返回 400 Bad Request]
    Valid -- Yes --> Save[寫入 Message_DB]
    Save --> Push{需推播?}
    Push -- No --> F200[返回 200 OK]
    Push -- Yes --> Job[啟動後台 Goroutine]
    Job --> Tokens[提取同社區之 FcmTokens]
    Tokens --> FCM[分批送往 Firebase Cloud Messaging]
    FCM --> Log[記錄發送數據]
    Log --> F200
    F200 --> End((結束))
    F400 --> End
```

---

## 核心功能互動流程 (Interactions)

#### 時序圖：後端聚合實作
```mermaid
sequenceDiagram
    participant App as Mobile Frontend
    participant Home as HomeController
    participant ReseRepo as ReservationRepository
    participant DB as PostgresDB

    App->>Home: POST /api/v1/reservation/quick (FacilityID, Time)
    Note over Home: 驗證 JWT Context 中的<br/>CommunityID & UserID
    Home->>ReseRepo: CreateQuickReservation(data)
    
    ReseRepo->>DB: SELECT count(*) FROM reservations WHERE facility_id=? AND time=?
    DB-->>ReseRepo: 尚有餘額
    
    ReseRepo->>DB: INSERT INTO reservation_info (data)
    DB-->>ReseRepo: Success
    
    ReseRepo-->>Home: Created Model
    Home-->>App: 201 Created (Updated Home Snapshot)
```

---

### 2. 管理員度：緊急公告與 FCM 推播 (Instant Announcement)

此功能確保管理員能快速發布資訊並強制推播至住戶手機。

#### 活動圖：發送流程
```mermaid
activityDiagram
    start
    :管理員首頁點擊「緊急公告」按鈕;
    :輸入標題與簡短內文;
    :確認發布;
    fork
        :後端持久化至 Message_DB;
    fork again
        :異步調用 Firebase 核心組件;
        :依照 CommunityID 查詢登錄的 FcmTokens;
        :下發 FCM Push Notification;
    end fork
    :顯示「發送成功」並更新首頁公告看板;
    stop
```

#### 時序圖：跨服務調度
```mermaid
sequenceDiagram
    participant Admin as Admin Web
    participant Home as HomeController
    participant MsgRepo as MessageRepository
    participant FB as pkg/firebase (Internal Service)
    participant FCM as Firebase Cloud Messaging (External)

    Admin->>Home: POST /api/v1/message/emergency (Title, Content)
    Home->>MsgRepo: CreateEmergencyMessage(data)
    
    MsgRepo->>database: INSERT INTO message_info
    database-->>MsgRepo: OK
    
    par 異步推播
        Home->>FB: SendNotificationToCommunity(communityID, payload)
        FB->>database: SELECT fcm_token FROM user_info WHERE community_id=?
        database-->>FB: List of Tokens
        FB->>FCM: SendMulticast(payload)
    end

    Home-->>Admin: 200 OK (Message Sent & Pushed)
```

---

### 3. 超級管理員：系統監控彙整 (System Health)

#### 時序圖：數據聚合
```mermaid
sequenceDiagram
    participant Dash as Super Admin Dashboard
    participant Home as HomeController
    participant CommRepo as CommunityRepository
    participant UserRepo as UserRepository

    Dash->>Home: GET /api/v1/home (Admin Identity)
    
    par 並行查詢
        Home->>CommRepo: GetPlatformGlobalStats()
        Home->>UserRepo: GetTotalUserCount()
        Home->>CommRepo: GetPendingCommunityApplications()
    end

    Home-->>Dash: SuperAdminHomeResponse (JSON)
```

---

## 核心功能模組 (Feature Modules)

首頁儀表板不僅是數據展示，更是關鍵功能的「快速入口」。

### 1. 住戶端功能模組
*   **快速預約 (Quick Booking)**：在即將到來的預約下方，提供「新增預約」按鈕，直接跳轉至最常使用的設施。
*   **包裹管理 (Parcel Tracking)**：若有待領取包裹，在儀表板中央顯示「待領包裹：2 份」紅色提醒，並附上領取 QRCode 按鈕。
*   **繳費中心 (Billing Hub)**：顯示當月管理費狀態。若未繳費，顯示「立即繳費」快捷導向。
*   **報修進度 (Maintenance Tracking)**：提供「我要報修」快捷鍵，並顯示最近一筆報修單的即時狀態（如：派工中、已完工）。

---

### 2. 社區管理員功能模組
*   **審核工作台 (Approval Workbench)**：集中顯示「待審核住戶」、「待審核預約」之數字標籤，點擊直接進入對應清單。
*   **群發通知 (Instant Announcement)**：提供「發送緊急公告」快捷按鈕，整合標題與簡短內文輸入，一鍵推播至全社區。
*   **設施看板 (Facility Dashboard)**：提供各設施當前使用率與清消進度，並具備「一鍵禁用」功能（應用於維修突發狀況）。

---

### 3. 超級管理員功能模組
*   **平台警報 (System Alerts)**：若有 API 回應延遲或 Firebase 推播失敗率過高，在首頁置頂顯示系統警告。
*   **入駐進度 (Onboarding Status)**：地圖式或清單式顯示各社區的活躍度，並提供「待審核社區」名單。

---

## Swagger 產生規範

為了讓 Swagger UI 正確顯示這些聚合結構，必須在 Controller 與 Model 中定義特定的註解。

### 1. Go Struct 定義範例
```go
// HomeResponse 首頁聚合回應結構
type HomeResponse struct {
    Role string      `json:"role" example:"resident"`
    Data interface{} `json:"data"` // 根據角色動態裝載
}

// ResidentDashboard 住戶視圖
type ResidentDashboard struct {
    UpcomingReservations []Reservation `json:"upcoming_reservations"`
}
```

### 2. Controller 註解範例 (進階：使用 Model Composition)
在 `app/controller/v1/home/Home_Controller.go` 的方法上方加入：

```go
// GetHomeDashboard 取得首頁儀表板
// @Summary 取得首頁儀表板資訊
// @Description 根據目前登入者的權限，聚合回傳對應的統計數據、預約清單與訊息。
// @Tags Home
// @Produce json
// @Success 200 {object} HomeResponse{data=ResidentDashboard} "住戶身份：成功返回儀表板數據"
// @Success 200 {object} HomeResponse{data=AdminDashboard} "管理員身份：成功返回統計與待辦"
// @Failure 401 {object} model.ErrorRequest "權限不足"
// @Router /api/v1/home [get]
// @Security BearerAuth
func (h *HomeController) GetHomeDashboard(ctx *gin.Context) { ... }
```
> [!TIP]
> 使用 `HomeResponse{data=SpecificStruct}` 語法能動態替換結構中的 `interface{}` 欄位型態，使 Swagger UI 能渲染出正確的 JSON schema 而非單純的 `object`。

### 3. 如何產生 Swagger 文件
1.  **安裝工具**：`go install github.com/swaggo/swag/cmd/swag@latest`
2.  **執行產生指令**：
    ```bash
    swag init -g main.go
    ```
    這會掃描 `main.go` 與所有引用的 package 註解，並更新 `docs/swagger.json`。
3.  **預覽網址**：`http://localhost:9080/swagger/index.html`

---

## 4. 資料聚合邏輯 (Sequence Diagram)

為了減少前端請求次數，首頁採「聚合式 API」設計。以下為資料聚合流向圖：

```mermaid
sequenceDiagram
    participant User as Frontend Client
    participant Controller as HomeController
    participant Auth as JWT Context
    participant Repo as Backend Repositories
    participant DB as PostgreSQL

    User->>Controller: GET /api/v1/home
    Controller->>Auth: 提取 PermissionID & CommunityID
    
    alt Super Admin (Level 1-2)
        Controller->>Repo: GetGlobalPlatformStats()
        Controller->>Repo: ListPendingApplications()
    else Community Admin (Level 3-7)
        Controller->>Repo: GetCommunityUserCount()
        Controller->>Repo: GetDailyReservationCount()
        Controller->>Repo: CountPendingRescheduleRequests()
    else Resident (Level 8-9)
        Controller->>Repo: GetReservationsByUserID()
        Controller->>Repo: GetRecentMessagesByCommunity()
    end

    Repo->>DB: 執行 SQL 查詢
    DB-->>Repo: 回傳數據結果
    Repo-->>Controller: 包裝 RepositoryModel
    Controller->>Controller: 聚合各項數據至 Response Struct
    Controller-->>User: 回傳 JSON Dashboard Data
```

---

## 5. API 規格定義 (規劃中)

### 5.1 Endpoint
`GET /api/v1/home`

### 5.2 預期回傳結構節錄
*   **住戶模式**: `{"upcoming_reservations": [...], "recent_messages": [...]}`
*   **管理員模式**: `{"total_residents": 120, "today_reservation_count": 15, "pending_reschedule_count": 3}`

---

## 6. 後續開發指南
1. **Repository 層**: 需擴充各模組的查詢方法，建議使用 GORM 的 `Count` 與 `Limit` 進行優化。
2. **安全性**: 嚴格檢查 `JWT` 內的 `community_id` 與 `permission_id`，禁止越權查詢他區或其他權限之數據。
