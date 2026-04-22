# 智慧通行 (Smart Access) 需求規格書

> [!NOTE]
> **Status**: Future Planning / Concept Only (非當前開發範圍)

本文件定義「智慧社區」中最重要的基礎設施：智慧門禁與訪客系統。

## 1. 核心業務場景 (Business Scenarios)

### 場景 A：住戶藍牙/手機感應
住戶到達大門或電梯口，無需尋找感應卡，開啟 APP 透過藍牙 (BLE) 或點擊「手機開門」即可通行。
*   **技術實作**：手機 APP 調用 `/api/v1/access/open` -> 後端校驗 JWT 與時間權限 -> 下發指令至 IoT Gateway -> 繼電器觸發開門。

### 場景 B：訪客動態邀請碼 (Visitor Pass)
住戶預計有訪客到訪，預先生成 QRCode 發送給訪客。

*   **動態加密設計 (Perfected TOTP)**：
    邀請碼採用 `TOTP (Time-based One-time Password)` 演算法，將 `CommunityID + UserID + Timestamp` 進行雜湊，每 30 秒更換一次。即便 QRCode 被截圖發送給他人，若超過有效時距將自動失效。

*   **電梯聯動邏輯 (Lift Control Logic)**：
    ```mermaid
    graph TD
        A[訪客感應門口機 QRCode] --> B{後端校驗 AccessKey}
        B -->|有效| C[觸發大門電子鎖]
        B -->|有效| D[發送指令至電梯控制器]
        D --> E[電梯召喚至 1 樓]
        E --> F[授權住戶指定樓層]
        C --> G[推播通知至住戶 APP]
        F --> G
    ```

*   **流程圖**：
    ```mermaid
    sequenceDiagram
        participant Resident as 住戶 APP
        participant Server as 後端伺服器 (Go)
        participant Visitor as 訪客手機
        participant Hardware as 門口辨識器 (IoT)

        Resident->>Server: 請求生成邀請碼 (VisitorID, 有效時間)
        Server-->>Resident: 回傳加密 AccessKey
        Resident->>Visitor: 透過 LINE/簡訊 分享 QRCode
        Visitor->>Hardware: 感應 QRCode
        Hardware->>Server: 校驗 AccessKey
        Server-->>Hardware: 通行許可 (Open Door)
        Server->>Resident: 推播通知「訪客已進入」
    ```

## 2. 功能清單 (Feature List)

| ID | 功能名稱 | 優先級 | 描述 |
| :--- | :--- | :--- | :--- |
| SA-01 | 動態 QRCode 生成 | `P0` | 基於 TOTP 演算法，每 30 秒更新一次，防止截圖盜用。 |
| SA-02 | 邀請碼有效期管理 | `P0` | 支援「一次性」或「時段重複性（如每週五 14:00）」授權。 |
| SA-03 | 門鎖狀態即時監控 | `P1` | 監測大門是否遭非法長時間開啟或強行進入。 |
| SA-04 | 訪客通行權限限制 | `P0` | 限制訪客僅能到達指定樓層與區域。 |

## 3. 安全性規範 (Security Requirements)

1.  **加密傳輸**：所有的通行權限交換必須經過雙層加密（AES-128）。
2.  **Rate Limiting**：限制單一住戶每天生成的邀請碼數量，防止資源濫用。
3.  **Audit Trail**：所有的開門行為必須包含 Timestamp, UserID, DeviceID 與結果，且記錄不可修改。

## 4. 硬體整合建議 (Integration)

*   **通訊協定**：建議使用 MQTT 或 WebSocket 與 IoT Gateway 連線，以降低通訊延遲。
*   **離線支援**：Edge 運算模組應儲存 24 小時內的住戶白名單，避免網路中斷導致住戶無法進門。
