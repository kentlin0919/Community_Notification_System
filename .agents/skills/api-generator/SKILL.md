---
name: api-generator
description: 自動產生專案模組的基礎架構 (Controller, Repository, Model, Route)，確保符合 Clean Architecture 規範。
---

# API Generator Skill

## When to use this skill
- 當使用者需要新增一個新的 API 業務模組時。
- 當需要確保新模組的 Package 名稱與目錄結構符合 `app/` 下的慣例時。

## How to use it
1. **確認模組名稱**：確認要建立的模組名稱（例如 `Announcement`）。
2. **分析現有結構**：參考 `app/controller/v1/user` 或 `app/repositories/user` 的實作方式。
3. **產生檔案**：
    - 在 `app/models/` 建立對應的模型目錄。
    - 在 `app/repositories/` 建立 Repository 實作。
    - 在 `app/controller/v1/` 建立 Controller 檔案。
    - 在 `routers/api/v1/v1.go` 註冊路由。
4. **檢查 Package**：確保所有產生的 Go 檔案 Package 名稱正確且 Import 路徑完整。

## Rules and Patterns
- 遵循 `controller -> repository -> model` 的依賴關係。
- 路由必須註冊在 `v1` 組下，且標註適當的 Swagger 註解。
- 檔名格式應為 `Module_Action.go` 或 `Module_Controller.go`。
