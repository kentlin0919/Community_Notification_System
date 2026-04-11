# 使用者刪除

## 功能目標
刪除既有使用者資料。

## API
- 方法：`POST`
- 路徑：`/api/v1/deleteUser`
- 是否需要 JWT：是

## 回傳格式

### 成功回傳 `202 Accepted`
```json
{
  "Message": "Delete Sucessful"
}
```

### 錯誤回傳
`400 Bad Request`
```json
{
  "code": 400,
  "status": "Bad Request",
  "error": "Invalid input"
}
```

`401 Unauthorized`
```json
{
  "error": "缺少 Authorization header"
}
```

或

```json
{
  "error": "格式錯誤，應為 Bearer token"
}
```

或

```json
{
  "error": "無效的 token"
}
```

`404 Not Found`
```json
{
  "code": 404,
  "status": "Not Found",
  "error": "User Not Found"
}
```

`500 Internal Server Error`
```json
{
  "code": 500,
  "status": "Internal Server Error",
  "error": "Delete Error"
}
```

## Activity Diagram
```mermaid
flowchart TD
    A[接收刪除請求] --> B[JWT Middleware 驗證]
    B --> C{Token 是否有效}
    C -- 否 --> Z1[回傳 401]
    C -- 是 --> D{JSON 是否有效}
    D -- 否 --> Z2[回傳 400]
    D -- 是 --> E[查詢 user_info]
    E --> F{使用者是否存在}
    F -- 否 --> Z3[回傳 404]
    F -- 是 --> G[刪除 user_info]
    G --> H{刪除是否成功}
    H -- 否 --> Z4[回傳 500]
    H -- 是 --> I[回傳刪除成功]
```

## Sequence Diagram
```mermaid
sequenceDiagram
    participant Client as Client
    participant JWTmw as JWT Middleware
    participant Ctrl as UserController
    participant Repo as UserRepository
    participant DB as PostgreSQL

    Client->>JWTmw: POST /api/v1/deleteUser + Bearer Token
    JWTmw-->>Ctrl: 驗證通過
    Ctrl->>Repo: LoginRepository(email)
    Repo->>DB: SELECT * FROM user_info WHERE email = ?
    DB-->>Repo: UserInfo / not found
    Ctrl->>Repo: UserDeleteRepository(user)
    Repo->>DB: DELETE FROM user_info
    Repo-->>Ctrl: delete result
    Ctrl-->>Client: 202 Accepted
```

## 實作備註
- 雖註解提到需提供帳號密碼，但實作實際上只先查詢 email 對應使用者後執行刪除。
- 刪除前沒有再次驗證密碼。
