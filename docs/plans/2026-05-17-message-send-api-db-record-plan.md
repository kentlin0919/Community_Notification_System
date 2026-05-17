# Message Send API 與寄送紀錄資料庫欄位 Implementation Plan

> **For Hermes:** Use subagent-driven-development skill to implement this plan task-by-task.

**Goal:** 完善傳送 message API，讓所有寄出訊息都會建立可追蹤的資料庫紀錄，並以社區範圍查詢收件者與回傳批次發送結果。

**Architecture:** 延續既有 Gin controller / repository / model 分層。API 以 JWT context 的 `community_id` 做資料隔離，先建立 `message_info` 寄送紀錄，再依使用者 `Fcmtoken` 嘗試 FCM 推播，最後回傳 target/success/failure 統計。

**Tech Stack:** Go 1.23、Gin v1.10、GORM v1.30、Firebase Admin SDK、Swagger/swag。

---

## Task 1: 補測試保護寄送紀錄流程

**Objective:** 先建立 controller 行為測試，驗證新 API 會接受 selected_users request 並建立寄送紀錄。

**Files:**
- Create: `app/controller/v1/message/Message_SendMessage_test.go`
- Read: `app/controller/v1/user/User_Login_test.go`

**Steps:**
1. 使用 SQLite in-memory 初始化 `database.DB`。
2. AutoMigrate `user_db.UserInfo` 與 `message_db.MessageInfo`。
3. 建立兩個同社區使用者。
4. 呼叫 `POST /messages/send`，context 帶入 `user_id`、`permission_id`、`community_id`。
5. 預期 response 為 200，且 `message_info` 有兩筆同 batch 紀錄。
6. 先執行測試確認失敗：`go test ./app/controller/v1/message -run TestSendMessageCreatesRecords -v`。

## Task 2: 設計 message request/response 與 DB 欄位

**Objective:** 增加 API contract 與寄送紀錄必要欄位。

**Files:**
- Modify: `app/models/message/Message_Model.go`
- Modify: `database/Message_DB/Message_Schema.go`

**Schema 欄位:**
- `SenderID string`
- `CommunityID uint64`
- `BatchID string`
- `TargetType string`
- `Category string`
- `EntityType string`
- `EntityID string`
- `FcmStatus string`
- `FcmMessageID string`
- `FcmError string`

**API models:**
- `SendMessageRequest`
- `SendMessageResponse`
- `SendMessageResponseData`
- `SendMessageResult`

## Task 3: 實作 repository 查詢與建立紀錄

**Objective:** 將收件者查詢與 message_info 建立集中在 message repository。

**Files:**
- Modify: `app/repositories/message/Message_repository.go`

**Functions:**
- `FindMessageRecipientsRepository(communityID uint64, req *message_model.SendMessageRequest)`
- `CreateMessageRecordsRepository(senderID string, communityID uint64, req *message_model.SendMessageRequest, recipients []*user_db.UserInfo, batchID string)`
- `UpdateMessageFCMResultRepository(messageID, status, fcmMessageID, fcmError string)`

**Rules:**
- `selected_users` 支援 `recipient_user_ids` 與 `recipient_emails`。
- `community` 查詢同社區全部使用者。
- 所有查詢都必須帶 `community_id = ?`。
- 沒有收件者時由 controller 回 404。

## Task 4: 改寫 controller 與 route

**Objective:** 實作新 API 並保留舊路徑相容。

**Files:**
- Modify: `app/controller/v1/message/Message_SendMessage.go`
- Modify: `routers/api/v1/v1.go`
- Modify: `routers/api/v2/v2.go`

**Routes:**
- Keep: `POST /api/v1/sendmessage`
- Add: `POST /api/v1/messages/send`
- Keep v2 existing behavior or add matching route if currently maps to v1 controller。

**Behavior:**
- 驗證 JWT context：`user_id`、`community_id`。
- 驗證 request：`title`、`body`、`target_type`。
- 先建立 DB 記錄。
- Firebase 未初始化時不阻止落庫；回傳 `fcm_status=skipped` 與 failure_count。
- FCM 部分失敗仍回 200 並在 body 顯示 failure_count。

## Task 5: 更新文件與 Swagger

**Objective:** 讓 README、功能文件與 Swagger 與新 API 同步。

**Files:**
- Modify: `README.md`
- Modify: `docs/features/message_send/README.md`
- Generated if available: `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`

**Steps:**
1. 更新 README 傳送訊息時序圖與資料表欄位說明。
2. 更新 message_send 文件：request/response、資料庫欄位、流程圖、時序圖。
3. 執行 `swag init -g main.go`，若 swag 不存在則記錄未執行原因。

## Task 6: 驗證

**Objective:** 確認編譯與測試通過。

**Commands:**
- `gofmt -w app/models/message/Message_Model.go database/Message_DB/Message_Schema.go app/repositories/message/Message_repository.go app/controller/v1/message/Message_SendMessage.go app/controller/v1/message/Message_SendMessage_test.go routers/api/v1/v1.go routers/api/v2/v2.go`
- `go test ./app/controller/v1/message -run TestSendMessageCreatesRecords -v`
- `go test ./... -v`
- `swag init -g main.go` if available

## Acceptance Criteria

- `POST /api/v1/messages/send` 可傳送 selected_users/community 訊息。
- `POST /api/v1/sendmessage` 保留相同行為。
- 每位收件者都會建立 `message_info` 紀錄，即使 Firebase 未初始化或 FCM 失敗。
- 寄送紀錄包含 sender、community、batch、target、category/entity、FCM 狀態欄位。
- 收件者查詢不得跨社區。
- 測試與 gofmt 通過。
