# Repository Guidelines

## Project Structure & Module Organization
The entry point lives in `main.go`, which wires middleware, routes, and Swagger. Domain logic sits under `app/`, dividing controllers, models, and repositories by feature and version (for example, `app/controller/v1/user`). Database schema helpers are in `database/`, while middleware such as JWT and CORS lives in `middlewares/`. Generated API docs stay in `docs/`; do not edit them manually. Temporary build artifacts from `air` land in `tmp/` and should remain untracked. No dedicated test directory exists yet—place Go tests alongside the packages they cover.

## Build, Test, and Development Commands
Run `go mod tidy` after adding dependencies to keep `go.mod` clean. Use `go run main.go` for a one-off local server, or `air` for hot-reload via `.air.toml`. Regenerate Swagger after updating annotations with `swag init -g main.go`. Execute `go test ./...` to run the full suite once tests exist.

## Coding Style & Naming Conventions
Format every Go change with `gofmt` (or `go fmt ./...`) before committing; the repo follows standard Go tab indentation and import sorting. Package names stay lowercase with no underscores. Exported types and functions use PascalCase (`UserController`), while locals prefer camelCase. Keep request/response structs in `app/models`, and match JSON field tags to the casing exposed in the API.

## Testing Guidelines
Add `_test.go` files beside the implementation files (`app/controller/v1/user/User_Login_test.go`). Favor table-driven tests and focus on controller behavior through Gin test contexts or repository logic with a mocked DB. Run `go test ./... -v` before opening a pull request and ensure new features include coverage.

## Commit & Pull Request Guidelines
Follow the existing conventional commits style—prefix messages with scopes like `feat:`, `fix:`, or `ref:` plus a short summary (e.g., `feat: add message sender`). For pull requests, include a concise description, testing notes, and screenshots or cURL examples for API changes. Link related tickets or issues in the description and update Swagger output when contracts change.

## Security & Configuration Tips
Store secrets only in `.env`; never commit that file. Verify `JWTPASSWORD` and database credentials before running migrations. Middleware assumes HTTPS when setting cookies—use secure transport in staging and production. Rotate JWT secrets if credentials leak and revoke affected tokens.

## 回覆語言
- 繁體中文

## README 書寫條件
- 要有詳細的版本資訊
- 詳細的檔案結構
- 針對每個功能畫出時序圖
- 針對每個系統的環境安裝的詳細流程以及指令
- 每次執行時都對 README.md 進行更新確保符合現在專案


## 提交與 Pull Request 規範
- Commit：簡短、命令式、標明範疇（如：`fix(chart): clamp pan at edges`）
- 語言 ： 請以繁體中文書寫
- 請將相關變更分組，避免不相關的重構
- PR 必須包含：
  - 變更摘要與原因
  - UI 變更請附截圖/GIF
  - 手動測試步驟與影響模組/路徑
  - 若適用請附上相關 issue 或任務 ID
- 分支命名：`feature/<name>`、`fix/<name>`、`chore/<name>`
- 每個 PR 需至少一位審核者
- 請保持 PR 規模小且聚焦（盡量少於 300 行）
- 請將所有commit 經過逐行分析後放進/docs中的commit_summary_2025_{當月月份}.md 的文件中
- commit_summary_2025_{當月月份}.md 請以時間新到舊排序

## 使用套件的版本
- 確保環境安裝腳本有確認是否安裝必要套件
