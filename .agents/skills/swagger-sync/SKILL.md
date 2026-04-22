---
name: swagger-sync
description: 自動化 Swagger API 文件的更新流程，包含語法檢查與 swag init 執行。
---

# Swagger Sync Skill

## When to use this skill
- 當修改了 Controller 函式的註解時。
- 當新增了 API 路由，需要更新 `/swagger/index.html` 的內容時。

## How to use it
1. **語法掃描**：檢查專案中 `// @` 開頭的註解是否符合 Swaggo 規範。
2. **更新主文件**：確保 `main.go` 中的全局資訊（如 `@host`, `@version`）是最新的。
3. **執行生成指令**：
    - 運行 `swag init -g main.go` 重新產生 `docs/`。
4. **驗證產出**：確認 `docs/docs.go` 與 `docs/swagger.json` 已更新。

## Rules and Patterns
- 嚴禁手動修改 `docs/` 下由 Swaggo 產生的檔案。
- 每個 Controller 動作必須包含 `@Summary`, `@Tags`, `@Success` 等基本標註。
- 確保 API 分類（Tags）與專案模組對應。
