# 使用者註冊

## 功能目標
建立新帳號，檢查帳號重複、加密密碼、產生 JWT 與 Session Cookie，並寫入 `user_info`。

## 社區上下文規則
- 註冊資料中的 `community_id` 應視為後續社區上下文的來源之一。
- 註冊成功後若立即簽發 JWT，建議將 `community_id` 一併寫入 claims。
- 若註冊時未綁定社區，後續受保護 API 應限制其只能操作非社區型功能。

## API
- 方法：`POST`
- 路徑：`/api/v1/register`
- 是否需要 JWT：否

## 回傳格式

### 成功回傳 `200 OK`
```json
{
  "message": "Register successful",
  "token": "<jwt-token>",
  "user_info": {
    "PermissionId": 0,
    "name": "",
    "email": "",
    "Home_id": "",
    "Birthdaytime": "0001-01-01T00:00:00Z",
    "Platform": 0,
    "Session_id": "",
    "community_id": 0
  }
}
```

註：目前程式只設定 `message` 與 `token`，`user_info` 未填值，序列化後會出現零值欄位。

### 錯誤回傳
`400 Bad Request`
```json
{
  "code": 400,
  "status": "Bad Request",
  "error": "密碼長度至少需 8 碼"
}
```

其他可能的 `400` 訊息：
- `Invalid input`
- `已經註冊過了`
- `Register error`

`500 Internal Server Error`
```json
{
  "code": 500,
  "status": "Internal Server Error",
  "error": "JWT 簽發失敗"
}
```

## Activity Diagram
```mermaid
flowchart TD
    A[接收註冊請求] --> B{JSON 是否有效}
    B -- 否 --> Z1[回傳 400 Invalid input]
    B -- 是 --> C{密碼長度 >= 8}
    C -- 否 --> Z2[回傳 400 密碼長度不足]
    C -- 是 --> D[查詢帳號是否已存在]
    D --> E{是否已註冊}
    E -- 是 --> Z3[回傳 400 已經註冊過了]
    E -- 否 --> F[bcrypt 產生雜湊密碼]
    F --> G[建立 UserInfo 與 UUID]
    G --> H[產生 JWT]
    H --> I[寫入 user_info]
    I --> J[設定 session_id Cookie]
    J --> K[回傳註冊成功]
```

## Sequence Diagram
```mermaid
sequenceDiagram
    participant Client as Client
    participant Ctrl as UserController
    participant Repo as UserRepository
    participant DB as PostgreSQL
    participant JWT as JWT Utility

    Client->>Ctrl: POST /api/v1/register
    Ctrl->>Repo: LoginRepository(email)
    Repo->>DB: SELECT * FROM user_info WHERE email = ?
    DB-->>Repo: Existing user / not found
    Repo-->>Ctrl: 檢查結果
    Ctrl->>Ctrl: bcrypt.GenerateFromPassword()
    Ctrl->>JWT: GenerateJWT(email)
    JWT-->>Ctrl: token
    Ctrl->>Repo: RegisterRepository(user_info)
    Repo->>DB: INSERT INTO user_info
    DB-->>Repo: create result
    Repo-->>Ctrl: success
    Ctrl-->>Client: 200 OK + token + cookie
```

## 實作備註
- `Register` 請求中的 `Email` JSON key 為大寫開頭，與登入的 `email` 命名不同。
- 註冊成功回應訊息目前為英文 `Register successful`。
- 目前沒有把註冊行為寫入 `user_log`。
- 後續若全面啟用社區隔離，註冊流程應驗證 `community_id` 是否存在，避免產生無效社區關聯。
