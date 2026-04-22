# Antigravity Skills 總覽

本目錄定義了 `Community_Notification_System` 專案中的 **Skills**。根據 Antigravity 框架，Skill 是一套標準化的指令集與工作流，用於指導 AI 代理人（或系統自動化程序）執行特定的複雜任務。

## 為什麼需要 Skills？

在大型社區通知系統中，許多操作涉及跨模組的邏輯（例如：審核申請單後觸發推播，並寫入日誌）。透過 Skills，我們可以：
1. **標準化流程**：確保每次執行相同任務時，邏輯路徑一致。
2. **提升 AI 效率**：讓 AI 助理能快速識別「工具」的使用時機。
3. **降低耦合**：將業務流程與底層 API 實作分離。

## Skills 分類

| 分類 | 說明 | 對應文件 |
| --- | --- | --- |
| **自動化邏輯 (Automation)** | 處理需要決策與邏輯判斷的自動化任務。 | [自動化審核](skill_auto_approval.md) |
| **系統維護 (Maintenance)** | 處理定期執行的資料清理與狀態同步。 | [系統維護](skill_maintenance_cron.md) |
| **通訊分發 (Communication)** | 處理訊息推播、Email 或簡訊的發送流。 | [訊息推播](skill_notification_dispatch.md) |

## Skill 組成結構

每個 Skill 文件均包含以下章節：
- **Definition (定義)**：該 Skill 的核心目標。
- **Context (上下文)**：執行時所需的必要資訊（如 `community_id`）。
- **Triggers (觸發條件)**：什麼情況下會啟動此 Skill。
- **Workflow (工作流)**：詳細的執行步驟，包含時序圖。
- **Constraints (限制)**：安全規範與權限要求。

## 執行環境

Skills 透過專案的 `v1/v2` 控制器實作，並受到以下中介層保護：
- `JWTAuthMiddleware`：身份驗證。
- `CommunityContextMiddleware`：租戶資料隔離。
- `PermissionMiddleware`：權限等級校驗。
