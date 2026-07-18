# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Development
go run main.go          # Start server on :9080
air                     # Hot-reload development (preferred)
air -c .air-debug.toml  # Hot-reload with Delve debugger on port 40000

# Dependencies & Docs
go mod tidy             # Sync/clean dependencies
swag init -g main.go    # Regenerate Swagger docs after annotation changes

# Testing
go test ./...                                        # All tests
go test ./app/controller/v1/user/... -v -run TestX   # Single test

# Docker
docker compose up --build  # Full stack: PostgreSQL + API + Delve
```

## Architecture

This is a **multi-tenant community notification system** — each community is isolated by `community_id` carried in the JWT. Always filter queries by `community_id` to prevent cross-community data leakage.

### Layered structure

```
Controller (app/controller/v1/)  ← HTTP handling, request parsing
    ↓
Repository (app/repositories/)   ← Database access via GORM
    ↓
Database (database/)             ← Schema definitions, auto-migration on startup
```

All repositories return a generic wrapper:
```go
type RepositoryModel[T any] struct {
    Statue gorm.DB  // error state + rows affected
    Result T        // actual data
}
```

Controllers are registered via factory functions in `app/controller/v1/v1.go`.

### Middleware chain (applied in `main.go`)

1. **CORSMiddleware** — allows credentials from any origin
2. **JWTAuthMiddleware** — validates Bearer token; skips `/login`, `/register`, `/platform/getlist`, `/swagger/*`
3. **CookieMiddleware** — reads `session_id` cookie
4. **CommunityContextMiddleware** — extracts `community_id`, `user_id`, `permission_id` from JWT claims into `gin.Context`
5. **PermissionMiddleware** — gates admin routes by `permission_id` (1 = Super Admin, 2+ = Community Admin/custom)

### Permission levels

- `1` — Super Admin: can approve/reject community registrations globally
- `2+` — Community Admin or custom roles: facility management, user management within one community

### Key packages

| Path | Purpose |
|------|---------|
| `app/controller/v1/` | Route handlers, one sub-directory per feature |
| `app/repositories/` | GORM queries, one file per feature |
| `app/models/` | Request/response structs and GORM model definitions |
| `database/` | Table schemas and `AutoMigrate` calls |
| `middlewares/` | CORS, JWT, cookie, community-context middleware |
| `routers/` | Route registration for v1 and v2 |
| `utils/` | JWT generation helpers |
| `pkg/firebase/` | Firebase Cloud Messaging (FCM) initialization — optional at startup |
| `configs/` | `godotenv` loader |
| `docs/` | Swagger generated API docs only (`docs.go`, `swagger.json`, `swagger.yaml`) |
| `/Users/kent/Library/Mobile Documents/iCloud~md~obsidian/Documents/Community_Notification_System_docs` | PRD, architecture docs, feature specs, monthly commit summaries. **MUST BE REFERENCED BEFORE IMPLEMENTATION** |

### Environment variables (`.env`)

```
PORT=:9080
DB_HOST=127.0.0.1
DB_USER=postgres
DB_PASSWORD=...
DB_NAME=db_Community
DB_PORT=5432
DB_TIMEZONE=Asia/Taipei
JWTPASSWORD=...
```

The app auto-creates the database by connecting to the `postgres` admin database if `DB_NAME` does not exist.

Firebase FCM is optional — place `serviceAccountKey.json` in the project root. If absent the server still starts; FCM endpoints return 503.

### Testing pattern

Tests use SQLite in-memory + Gin test context. See `app/controller/v1/user/User_Login_test.go` for the setup pattern. Tests live alongside their source files (`_test.go` suffix, same package).

### Swagger

Swagger annotations live on controller handler functions. Run `swag init -g main.go` after any annotation change. View generated docs at `/swagger/index.html`.

### Commit convention

Follow conventional commits: `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, etc. Update `/Users/kent/Library/Mobile Documents/iCloud~md~obsidian/Documents/Community_Notification_System_docs/commit_summaries/commit_summary_YYYY_MM.md` when adding major features.


# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.