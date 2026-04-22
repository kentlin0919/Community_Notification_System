---
name: test-generator
description: 協助為 Go 專案組件產生 Table-driven 測試模板，加速開發與驗證流程。
---

# Test Generator Skill

## When to use this skill
- 當新增了 Controller 或 Repository 邏輯，需要建立單元測試時。
- 當需要針對現有功能補充測試案例以提升覆蓋率時。

## How to use it
1. **定位目標檔案**：確定要測試的 Go 檔案路徑。
2. **分析函式簽名**：識別需要測試的函式與依賴項目。
3. **建立測試檔案**：在同目錄下建立 `filename_test.go`。
4. **撰寫模板**：
    - 建立測試結構體 `tests := []struct{...}`。
    - 使用 `gin.CreateTestContext` 模擬 HTTP 請求。
    - 使用 `assert` 庫（若專案有引入）或標準 `t.Errorf` 驗證結果。
5. **執行測試**：運行 `go test -v path/to/pkg`。

## Rules and Patterns
- 測試檔案必須放在實作檔案旁邊。
- 優先進行「邊界案例」測試（如：空輸入、非法權限）。
- 測試名稱應清楚描述測試意圖（如 `TestUserLogin_Success`）。
