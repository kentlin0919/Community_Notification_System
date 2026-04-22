---
name: permission-checker
description: 靜態檢查 API 路由設定，確保敏感操作 (如 /admin) 有正確套用權限驗證中介層。
---

# Permission Checker Skill

## When to use this skill
- 當在 `routers/` 中新增了管理端（Admin）API 時。
- 在發布 PR 前，進行安全性自我稽核時。

## How to use it
1. **路由掃描**：掃描 `routers/api/v1/v1.go`。
2. **路徑分析**：找出所有包含 `/admin/`, `/delete`, `/update` 等敏感關鍵字的路徑。
3. **檢查 Middleware**：確認這些路由是否被包在套用了 `middlewares.PermissionMiddleware(level)` 的路由組中。
4. **回報漏洞**：若發現公開（未受限）的管理員介面，應立即標註並修正。

## Rules and Patterns
- 社區級別的管理 API 必須檢查 `community_id` 上下文。
- Super Admin API 必須要求最高權限等級。
- 檢查結果應記錄在 `docs/analysis/security_audit.md`（若存在）。
