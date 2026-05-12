# 智慧訪客系統流程圖 (Activity Diagram: Visitor Pass)

本文件描述智慧訪客系統 (Visitor Management 2.0) 的活動流程。

## 流程圖

```mermaid
sequenceDiagram
    participant Resident as 住戶 (APP)
    participant Backend as 後端伺服器 (API)
    participant DB as 資料庫 (Visitor_Passes / Access_Logs)
    participant IoT as 門口機 (IoT Device)
    participant Visitor as 訪客

    Resident->>Backend: 1. 輸入訪客資訊與預期時間
    Backend->>DB: 2. 建立 Visitor_Pass 紀錄
    DB-->>Backend: 3. 回傳 Pass ID 與加密金鑰
    Backend-->>Resident: 4. 回傳 QRCode 內容 (加密金鑰)
    Resident->>Visitor: 5. 透過社群軟體分享 QRCode
    
    Visitor->>IoT: 6. 在門口機掃描 QRCode
    IoT->>Backend: 7. 發送驗證請求 (Token/金鑰)
    Backend->>DB: 8. 查詢 Visitor_Pass 是否有效
    
    alt 金鑰有效且在時間內
        DB-->>Backend: 9a. 回傳有效
        Backend->>DB: 10a. 寫入 Access_Logs (通行紀錄)
        Backend-->>IoT: 11a. 回傳開門指令
        IoT-->>Visitor: 12a. 門鎖開啟
        Backend->>Resident: 13a. FCM 推播「訪客已進入」
    else 金鑰無效或逾期
        DB-->>Backend: 9b. 回傳無效
        Backend->>DB: 10b. 寫入 Access_Logs (拒絕紀錄)
        Backend-->>IoT: 11b. 回傳拒絕指令
        IoT-->>Visitor: 12b. 顯示錯誤或警示
    end
```

## 關鍵技術節點說明
1. **QRCode 安全性**：金鑰中應包含 `PassID` 與 `Signature`，避免被輕易竄改。
2. **時效性檢查**：後端驗證時，必須嚴格檢查 `StartTime` 與 `EndTime`。
3. **單次/多次使用**：若為單次使用，後端在回傳開門指令前，需同步將 `Visitor_Pass` 標記為 `IsUsed = true`。
4. **推播通知**：當訪客成功刷入，系統調用現有的 FCM 服務發送通知給發出邀請的 `OwnerID`。
