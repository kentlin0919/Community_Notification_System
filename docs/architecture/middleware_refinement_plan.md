# Middleware 優化與落實計畫 (Refinement Plan)

## 1. 核心目標
將「共用規則」從業務邏輯（Controller）中徹底拆分，確保 Controller 只負責處理具體的商業流程，而驗證、權限控制與資料隔離則交由 Middleware 處理。

## 2. 尚未落實的現有組件 (Existing Gaps)

### A. 權限驗證落實 (Permission Authorization)
*   **現況**：已實作 `MinimalPermissionMiddleware`，但未在 `routers/api/v1/v1.go` 中套用。
*   **影響**：目前的 API 只要有 JWT 就能進入，缺乏細粒度的權限控制（例如住戶可以呼叫管理員 API）。
*   **待辦事項**：
    - [ ] 在路由層使用 `Group` 並掛載 `MinimalPermissionMiddleware(level)`。
    - [ ] 移除 Controller 中手動判斷 `permission_id` 的邏輯。

### B. 社區資料隔離 (Data Isolation / Context)
*   **現況**：已實作 `CommunityContextMiddleware`，但僅在註解中提及，未實際綁定路由。
*   **影響**：Controller 內充斥著 `ctx.Get("community_id")` 的存在檢查與跨社區權限手動判斷。
*   **待辦事項**：
    - [ ] 將需要社區上下文的路由（如設施預約、社區管理）歸類到 `communityGroup`。
    - [ ] 讓 Middleware 統一處理 `community_id` 的合法性檢查，若不存在則直接 `AbortWithStatusJSON(401)`。

## 3. 建議新增的共用處理層 (Proposed Middlewares)

根據 Middleware 的價值，建議後續新增以下層級：

| 功能層級 | 價值 | 實作建議 |
| :--- | :--- | :--- |
| **全域錯誤處理 (Error Handling)** | 統一回傳格式 | 捕捉 `panic` 或處理 Controller 回傳的標準錯誤，避免程式碼中到處都是 `if err != nil { ctx.JSON(...) }`。 |
| **結構化日誌 (Logging)** | 追蹤請求脈絡 | 紀錄每個請求的 Request ID、執行時間、路徑與狀態碼，方便線上問題排查。 |
| **速率限制 (Rate Limiting)** | 防止惡意攻擊 | 針對 `Login` 或 `Register` 等 API 進行限制，防止暴力破解。 |
| **請求驗證 (Validation)** | 預先擋掉無效請求 | 針對重複性的 Body 格式進行預檢查，確保進到業務層的資料是乾淨的。 |

## 4. 控制器 (Controller) 重構範例

### 優化前 (目前的做法)
```go
func (c *Controller) Handle(ctx *gin.Context) {
    // ❌ 共用規則：驗證登入
    userID, _ := ctx.Get("user_id")
    
    // ❌ 共用規則：檢查權限
    permID, _ := ctx.Get("permission_id")
    if permID > 2 {
        ctx.JSON(403, "權限不足")
        return
    }

    // ✅ 真正的工作：處理業務
    logic.DoSomething(userID)
}
```

### 優化後 (理想做法)
```go
// 路由層負責規則
adminGroup := r.Group("/admin").Use(Auth(), Permission(2))

// Controller 只負責業務
func (c *Controller) Handle(ctx *gin.Context) {
    // 假設資料已經在 Middleware 驗證並準備好了
    userID := ctx.MustGet("user_id").(string)
    logic.DoSomething(userID)
}
```

---
*文件更新日期：2026-04-17*
