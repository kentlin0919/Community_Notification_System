---
name: db-migration-helper
description: 協助管理資料庫變動，確保 GORM 模型正確註冊到 AutoMigrate 流程中。
---

# DB Migration Helper Skill

## When to use this skill
- 當新增或修改了 `database/` 目錄下的 Schema 檔案時。
- 當需要確保 `database/db.go` 的 `InitDB` 函式包含最新的 AutoMigrate 模型時。

## How to use it
1. **分析模型變動**：檢查 `database/` 下各模組的 Schema 定義。
2. **更新註冊**：打開 `database/db.go`，在 `db.AutoMigrate(...)` 中加入新的結構體。
3. **驗證結構**：確保模型標籤（如 `gorm:"primaryKey"` 或 `json:"id"`）符合專案規範。
4. **檢查連線**：確保 `.env` 中的資料庫資訊正確，以便本機測試。

## Rules and Patterns
- Schema 應依功能分類放在 `database/` 的子目錄中。
- 欄位名稱應使用 PascalCase，JSON 標籤應為 snake_case。
- 複雜的關聯必須在 Schema 註解中說明。
