# Todo: MCP Server (Database Access)

Scope: MCP server on `/api/mcp` (Streamable HTTP) with bearer auth and tools
`ping` / `whoami` / `run_sql`. `run_sql` executes raw SQL including DDL/DML. No audit table,
no query logging.

Plan: `.opencode/planning/tasks/plan.md`

Out of scope: OAuth, PostgreSQL rate limiting, SSE tuning, testcontainers
integration tests (no local Postgres).

## Phase 1 — Foundations

- [x] Task 1: MCP config fields + `.env.example`
  - [x] Add `MCP_ENABLED`, `MCP_PATH`, `MCP_BEARER_TOKEN`, `MCP_ACTOR_USER_ID`,
        `MCP_SCOPES`, `MCP_MAX_BODY_BYTES`, `MCP_ALLOWED_ORIGINS`, `MCP_SESSION_TIMEOUT`,
        `MCP_STATELESS` to the envconfig struct with documented defaults
  - [x] Fail startup when `MCP_ENABLED=true` and `MCP_BEARER_TOKEN` is empty
  - [x] Compute `bearerHash` (SHA-256) once at load time
  - [x] Add the `MCP_*` block to `.env.example` (never edit `.env`)
  - [x] Verify: `go build ./...`, `go vet ./...` clean
  - [x] Verify: unit test for defaults + empty-token failure

- [x] Task 2: Core files — auth, server, handler + system tools
  - [x] `internal/mcp/auth.go`: bearer token parse + `subtle.ConstantTimeCompare`
  - [x] `internal/mcp/server.go`: `mcp.NewServer`, register `ping` + `whoami` via `mcp.AddTool`
  - [x] `internal/mcp/handler.go`: `WithAuth` → Streamable HTTP handler; body-size limit;
        origin allowlist + OPTIONS/CORS
  - [x] `internal/mcp/tools/system.go`: `ping` (version + time), `whoami` (actor, scopes)
  - [x] `whoami` without actor → safe tool error, no panic
  - [x] No audit or database write code in the core package
  - [x] Verify: `go build ./...`, `go vet ./...`, `go test ./...` clean
  - [x] Verify: unit tests for auth compare, origin check, body size, missing actor

### Checkpoint — Foundations

- [x] Build/vet/test pass
- [x] Auth rejects missing/invalid tokens; system tools respond via the handler
- [x] No audit/DB code in the core package
- [x] Human reviews before the SQL tool

## Phase 2 — Database Tools

- [x] Task 3: reverted — no user repository change (raw SQL instead of `search_users`)
  - [x] `UserRepository` keeps only `FindByEmail`
  - [x] No migration and no `mcp_audit_logs` table — MCP must not create tables
  - [x] Verify: `go build ./...`, `go vet ./...` clean

- [x] Task 4: `run_sql` tool
  - [x] `internal/mcp/tools/sql.go` with `jsonschema`-tagged input/output structs
  - [x] `internal/mcp/sql_executor.go`: `*sql.DB` executor, no allowlist
  - [x] Register in `server.go`; wire config `MCP_SQL_MAX_ROWS` / `MCP_SQL_TIMEOUT`
  - [x] `RETURNING` statements read as queries; others report `rowsAffected`
  - [x] Rows capped at `MCP_SQL_MAX_ROWS`; timeout per statement
  - [x] Missing actor, empty SQL, nil DB → safe tool errors
  - [x] Verify: `go build ./...`, `go vet ./...`, `go test ./...` clean
  - [x] Verify: unit tests for classification, normalization, tool call via SDK

### Checkpoint — Database Tools

- [x] `run_sql` works through the handler with a fake executor
- [x] No new database table or migration from MCP
- [x] Build/vet/test pass
- [ ] Human reviews before bootstrap/routing work

## Phase 3 — Integration

- [x] Task 5: `internal/bootstrap/mcp.go`
  - [x] `newMCPHandler(cfg, sqlDB)` returns nil when `MCP_ENABLED=false` or config nil
  - [x] Builds `mcp.HandlerConfig` from config (auth, scopes, origins, limits, SQL settings)
  - [x] `Application.MCPHandler` field wired in `New()`
  - [x] No tool logic; no `internal/routes` import from bootstrap
  - [x] Verify: `go build ./...`, `go vet ./...` clean
  - [x] Verify: `MCP_ENABLED=false` → handler nil; configured token accepted; missing token 401

- [x] Task 6: `internal/routes/mcp.go` + mount `/api/mcp`
  - [x] Mount via `adaptor.HTTPHandler` at `MCP_PATH` for all methods (`app.All`)
  - [x] Registered before `RegisterAPI` / `RegisterWeb` in both `Router()` and `New()`
  - [x] No-op when handler nil (`MCP_ENABLED=false`) → `/api/*` JSON 404 unchanged
  - [x] Verify (Fiber `app.Test`): GET/POST/DELETE without token → 401
  - [x] Verify: configured token reaches the handler
  - [x] Verify: SPA fallback and `/api/v1/*` 404 unchanged
  - [x] Verify: `go build ./...`, `go vet ./...`, `go test ./...` clean

### Checkpoint — Complete

- [x] All Task 1–6 acceptance criteria met (Task 3 reverted for raw SQL)
- [x] `ping`, `whoami`, `run_sql` registered on `/api/mcp`
- [x] Auth, size limit, and origin rules covered by unit tests
- [x] Route order correct; `/api/*` 404s unchanged
- [x] No changes under `docs/`, `tests/`, or `.github/`; `.env` untouched
- [ ] Human reviews raw-SQL exposure before deployment

## Deferred

- [ ] Write tools beyond raw SQL (`mcp:write`)
- [ ] PostgreSQL `rate_limits` integration
- [ ] SSE timeout/proxy-buffering tuning
- [ ] Testcontainers integration + contract/golden tests
- [ ] Per-tool scope enforcement beyond `mcp:read`

## Standing Rules

- Never log or commit `MCP_BEARER_TOKEN`; only `.env.example` gets a placeholder.
- Keep auth on the `net/http` side; Fiber locals are unavailable after the adaptor.
- Register `/api/mcp` before the SPA fallback.
- Tools call the same use cases as HTTP handlers when they exist; raw SQL is the exception.
- MCP has no migration and no audit table; query text is never logged.
- English comments; do not edit `docs/`; do not create `tests/` or `.github/`.
- Verify with `go build ./...` / `go vet ./...` / `go test ./...`; do not run the binary.
