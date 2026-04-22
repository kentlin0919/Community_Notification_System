# 智慧社區功能矩陣與發展藍圖 (Roadmap)

本文件基於「智慧社區」核心價值：**安防、便捷、節能、互助**，對系統功能進行深度分析。

## 1. 功能成熟度矩陣 (Functional Maturity Matrix)

| 模組 | 功能名稱 | 狀態 | 描述 |
| :--- | :--- | :--- | :--- |
| **基礎管理** | 使用者/社區/權限 | `DONE` | 支援多租戶、JWT 權限分流 (RBAC)。 |
| **基礎管理** | 設施預約/包裹 | `DONE` | 基礎預約邏輯與包裹狀態追蹤。 |
| **智慧通行** | 手機 QRCode 開門 | `DOC` | 利用時間戳記與住戶 ID 生成動態碼，整合門口機。 |
| **智慧通行** | 訪客邀請碼 | `TODO` | 下發臨時權限金鑰給訪客，含區域與時間限制。 |
| **智慧生活** | 瓦斯/度數表登錄 | `TODO` | 每月自動提醒登錄，與後端帳單系統整合。 |
| **智慧生活** | 聚合式首頁儀表板 | `DONE` | UI 已設計，API 結構已定義 (Aggregator API)。 |
| **安防通知** | 緊急公告 (FCM) | `DONE` | 整合 Firebase 下發全社區推播。 |
| **安防通知** | AI 跌倒/入侵偵測 | `FUTURE` | 整合第三方影像辨識串流與警報下發。 |
| **綠能節能** | EV 充電樁管理 | `FUTURE` | 充電位預約、計費管理與電力監控。 |

## 2. 發展階段規劃 (Development Phases)

> [!NOTE]
> 為了確保系統核心穩定，IOT 與自動化連動功能目前不在第一階段開發範圍，將作為 Phase 3 的長期規劃。

### Phase 1: 基礎社區平台 (當前重點)
*   使用者帳號、社區管理。
*   基礎設施預約 (無硬體連動)。
*   文字/圖片通知推播。

### Phase 2: 管理流程數位化 (進行中)
*   訪客登記流程 (純軟體紀錄)。
*   報修系統數位化。
*   首頁儀表板資訊聚合。

### Phase 3: 未來智慧擴展 (Future Planning)
*   **IOT 門禁實裝**：手機 QRCode 開門、電梯聯動。
*   **能源監控**：智慧度數表、EV 充電樁負載管理。
*   **AI 安防**：影向偵測與自動警報。

## 3. 核心增益功能詳細定義 (Perfecting missing parts)

### A. 智慧訪客系統 (Visitor Management 2.0)
*   **痛點**：訪客到達需對講機溝通，容易造成住戶隱私曝露。
*   **智慧化設計**：
    *   住戶預先於 APP 輸入訪客資訊 -> 生成 QRCode -> 訪客感應進門 -> 推播「訪客已進入」至住戶手機。
    *   **技術缺口**：需新增 `Visitor_Access_Logs` 表與暫存在 Redis 的 `SessionToken`。

### B. 公共資產與報修閉環 (Maintenance Lifecycle)
*   **痛點**：目前的報修僅有狀態顯示，缺乏維修軌跡與備品管理。
*   **智慧化設計**：
    *   管理員可指派物業服務人員 -> 工作人員掃描設備標籤 (Tag) -> 上傳修復照片 -> 系統自動記錄維修歷史。
    *   **技術缺口**：在 `Facility_DB` 中擴展 `Maintenance_Record` 結構。

### C. 雲端社區服務聚合 (Service Aggregator)
*   **智慧化設計**：
    *   整合特約商店優惠、家政服務預約、洗衣代收。
    *   **技術缺口**：建立標準 API Adapter，用於對接外部廠商。

## 3. 系統架構優化建議 (Codex Compliance)

為了達到「規整」與「完美」的要求，現有架構應進行以下技術債修正：
1.  **統一錯誤代碼 (Global Error Codes)**：建立 `utils/errors/codes.go`，避免字串直接寫死。
2.  **DTO (Data Transfer Object) 隔離**：確保 `app/models` 中 Database Schema 與 API Request/Response 結構完全隔離。
3.  **Traceability**：在所有 `middlewares` 中加入 `X-Request-ID`，方便追蹤跨模組請求（尤其是 Aggregator API）。

## 4. 下一步計畫

1.  **文件規整化**：將檔案按照 `/docs` 的新架構重新編排。
2.  **完善設計文件**：產出 `Activity_Diagram_Visitor_Pass.md`。
3.  **更新資料庫設計**：在 `docs/system/Database_Design.md` 補齊上述 `Visitor_Access_Logs` 與 `Maintenance_Record`。
