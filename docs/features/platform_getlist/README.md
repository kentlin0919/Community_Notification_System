# 平台列表查詢

## 功能目標
提供前端取得系統支援的平台名稱清單。

## API
- 方法：`GET`
- 路徑：`/api/v1/platform/getlist`
- 是否需要 JWT：否

## 回傳格式

### 成功回傳 `200 OK`
```json
{
  "total": 3,
  "platforms": [
    { "name": "web" },
    { "name": "App" },
    { "name": "Desktop" }
  ]
}
```

### 錯誤回傳
`500 Internal Server Error`
```json
{
  "code": 500,
  "status": "Internal Server Error",
  "error": "取得平台資料失敗"
}
```

## Activity Diagram
```mermaid
flowchart TD
    A[接收平台查詢請求] --> B[查詢 platform_info]
    B --> C{查詢是否成功}
    C -- 否 --> Z1[回傳 500]
    C -- 是 --> D[轉成 PlatformListResponse]
    D --> E[回傳平台清單]
```

## Sequence Diagram
```mermaid
sequenceDiagram
    participant Client as Client
    participant Ctrl as PlatformController
    participant Repo as PlatformRepository
    participant DB as PostgreSQL

    Client->>Ctrl: GET /api/v1/platform/getlist
    Ctrl->>Repo: PlatformRepository()
    Repo->>DB: SELECT * FROM platform_info
    DB-->>Repo: platform rows
    Repo-->>Ctrl: total + items
    Ctrl-->>Client: 200 OK + platforms[]
```

## 實作備註
- 預設資料為 `web`、`App`、`Desktop`。
