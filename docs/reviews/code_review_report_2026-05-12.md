# Community Notification System Backend Code Review Report

- 日期：2026-05-12
- Reviewer：Hermes CLI Agent
- 範圍：main 啟動流程、routing、middleware、user/message controller、message repository、db 初始化
- 審查方式：靜態程式碼審查（未啟動整合環境）

## 一、整體結論

目前程式可看出已逐步往「分層（controller/repository/model）」靠攏，但離 clean architecture 仍有明顯落差：

1. 邊界責任尚未清楚（controller 直接操作 DB 與 session/cookie，repository 只承載部分資料邏輯）。
2. 驗證、授權、錯誤回應格式未統一，導致 API 行為不一致。
3. middleware 鏈存在可維運風險（JWT 例外路徑硬編碼、token 驗證後未 `c.Next()`）。
4. 架構上強耦合全域狀態（`database.DB`, `firebase.FcmClient`），不利測試與擴充。

---

## 二、關鍵問題（Critical）

### C1. JWT middleware 在成功驗證後未明確呼叫 `c.Next()`
- 檔案：`middlewares/jwt_middleware.go:13-83`
- 現象：skip path 有 `c.Next()`；失敗時 `Abort()`；但成功分支最後沒有 `c.Next()`。
- 風險：目前 Gin 可能因 handler return 後繼續 chain，但此寫法語意不明確，容易在後續維護時產生誤判或 side effect（尤其 middleware 增加後）。
- 建議：成功驗證後明確 `c.Next(); return`，避免流程語意不完整。

### C2. 全域套用 JWT，再用 skipPaths 白名單「反向排除」
- 檔案：`main.go:40-43`、`middlewares/jwt_middleware.go:16-24`
- 現象：JWT middleware 掛在全域路由，靠路徑字串排除 `/api/v1/login` 等。
- 風險：
  - 新增公開 API 時容易忘記加白名單造成誤封鎖。
  - 路由重構（path 改名）時容易漏改而出現認證漏洞或誤攔截。
- 建議：改為 route group 分層（public/private），而非字串白名單判斷。

### C3. Controller 與資料層耦合過深，破壞清晰分層
- 檔案：`app/controller/v1/user/User_Login.go:92-97`
- 現象：`UserLogin` 在 controller 直接 `database.DB.Model(...).Updates(...)`，繞過 repository。
- 風險：
  - 商業規則散落 controller，難測試、難追蹤。
  - 後續 transaction 或審計需求無集中入口。
- 建議：抽到 application/usecase/service，再由 repository 實作 DB 行為。

---

## 三、高風險設計缺陷（Warnings）

### W1. 錯誤回應模型混用，不一致
- 檔案：`app/models/model/model.go` + 多個 controller
- 現象：同專案同路徑中，同時使用 `model.NewErrorRequest`、`NewGlobalErrorRequestWithMsg`、以及 `gin.H`。
- 影響：前端難以穩定解析；監控與錯誤碼對帳困難。
- 建議：定義單一錯誤回應 contract（含 `request_id`、`code`、`internal_code`、`message`）。

### W2. Message controller 引用錯誤的 repository package 與殘留 debug
- 檔案：`app/controller/v1/message/Message_SendMessage.go:6,44`
- 現象：import `app/repositories/user` 來拿 `UserInfoListRepository`，且 `fmt.Print(UserInfoList)`。
- 影響：邊界語意混亂（message domain 依賴 user repo 細節），debug 輸出污染 log。
- 建議：搬到 message usecase 層；移除 `fmt.Print`；統一 logger。

### W3. 請求與回應語言/訊息格式不一致
- 例子：`UserRegister` 同時有中英混用（`Invalid input`, `Register successful`）。
- 影響：產品體驗與維運（客服/稽核）成本增加。
- 建議：建立 i18n 或至少統一 API 訊息語系。

### W4. DB 初始化與 migration 綁在啟動流程
- 檔案：`database/db.go:115-142`
- 現象：啟動即執行建表。
- 風險：
  - 多實例/高併發部署時，schema 變更不可控。
  - 不利灰度與回滾。
- 建議：將 migration 與 app runtime 拆離（CI/CD migration job）。

---

## 四、一般改善建議（Suggestions）

1. `router.Run(port)` 應處理錯誤，且 `PORT` 應提供預設值或啟動前驗證。  
2. JWT claims 型別轉換建議集中封裝（目前 middleware 用 float64/int 混轉）。  
3. `checkAccount` 命名與語意可更明確（如 `isAccountAvailable`）。  
4. `Statue` 欄位拼字（疑似 `Status`）建議統一，避免誤用。  
5. 新增 request_id 後，錯誤回應建議帶上 `X-Request-ID` 以便追蹤。

---

## 五、對目前架構的質疑（Architecture Challenges）

以下是我建議你們團隊「必須回答」的架構問題：

### Q1. 這個專案目前到底是「分層 MVC」還是「Clean Architecture」？
你們文件要求 clean architecture，但目前 controller 仍直接 touching DB 與 infra（cookie、firebase、gorm）。

若要落實 clean architecture，至少要回答：
- Usecase/Application Service 層在哪裡？
- Domain 規則（例如註冊、登入、權限判斷）是否可脫離 Gin/Gorm 執行？

### Q2. 為什麼認證策略是「全域攔截 + 白名單例外」？
這是最容易在功能成長後出洞的模式。是否評估過：
- Public route group + Auth route group 的明確界線？
- 權限策略（RBAC）是否應該從 path 字串遷移到 policy 層？

### Q3. 為何 message domain 需要知道 user repository 細節？
目前看起來 message flow 為了找 user list 直接依賴 user repository。這會造成跨 domain 耦合。

應該討論：
- 是否需要 NotificationRecipientProvider 介面化？
- message 只是「傳送」還是也要負責「收件者決策」？

### Q4. 全域 singleton（`database.DB`, `firebase.FcmClient`）是否能支撐測試與多環境？
目前測試會很難做 dependency injection / mock。

應評估：
- 是否以 constructor 注入取代全域變數？
- 是否建立 `AppContainer` 或 wire/fx 進行組裝？

### Q5. API 錯誤模型要不要當成「穩定契約」管理？
現在錯誤回應有多種 shape。若前端與第三方整合擴大，這會直接變成維運風險。

---

## 六、建議的重構優先順序（低風險、可漸進）

1. 先修 middleware 流程語意與 route group（public/private）。
2. 建立統一 ErrorResponse（含 internal code/request id）。
3. 將 controller 內 DB 直連搬到 usecase + repository。
4. 把 message 與 user 的跨 domain 依賴改成介面。
5. 最後才處理全域 singleton 到 DI。

---

## 七、這次審查的限制

- 本次為靜態 code review，未實際串接資料庫與 Firebase 執行 E2E。
- 尚未覆蓋每一個 controller（以高風險路徑抽樣）。
- 建議下一步補一輪：`go test ./...` + middleware/route 整合測試 + auth matrix 測試。
