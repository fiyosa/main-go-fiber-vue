# Todo: MCP Server (Read-Only)

Scope: read-only MCP server on `/api/mcp` (Streamable HTTP) with bearer auth, tools
`ping` / `whoami` / `search_users`, and one audit row per tool call.

Plan: `.opencode/planning/tasks/plan.md`

Out of scope: write tools, OAuth, PostgreSQL rate limiting, SSE tuning, testcontainers
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

- [ ] Task 2: Core files — auth, server, handler, audit + system tools
  - [ ] `internal/mcp/auth.go`: bearer token parse + `subtle.ConstantTimeCompare`
  - [ ] `internal/mcp/server.go`: `mcp.NewServer`, register `ping` + `whoami` via `mcp.AddTool`
  - [ ] `internal/mcp/handler.go`: `WithAuth` → Streamable HTTP handler; body-size limit;
        origin allowlist + OPTIONS/CORS
  - [ ] `internal/mcp/audit.go`: one `mcp_audit_logs` row per tool call (nil-DB no-op)
  - [ ] `internal/mcp/tools/system.go`: `ping` (version + time), `whoami` (actor, scopes)
  - [ ] `whoami` without actor → safe tool error, no panic
  - [ ] No tool accepts raw SQL, table name, file path, shell command, or URL
  - [ ] Verify: `go build ./...`, `go vet ./...`, `go test ./...` clean
  - [ ] Verify: unit tests for auth compare, origin check, body size, missing actor

### Checkpoint — Foundations

- [ ] Build/vet/test pass
- [ ] Auth rejects missing/invalid tokens; system tools respond via the handler
- [ ] Audit wrapper proven incl. nil-DB no-op
- [ ] Human reviews before `search_users` + audit migration

## Phase 2 — Domain Tools

- [ ] Task 3: `mcp_audit_logs` migration + `UserRepository.SearchUsers`
  - [ ] Add `000003_create_mcp_audit_logs.up.sql` / `.down.sql` following the existing pattern
  - [ ] Columns: id, tool_name, parameters, result_summary, error_msg, ip_address,
        user_agent, created_at
  - [ ] Add `SearchUsers(ctx, email, page, pageSize)` to the user repository interface
  - [ ] Implement it in the GORM repo: max 50 rows, nil-safe (empty slice, no panic)
  - [ ] Wire the migration into `internal/app/infrastructure/commands/migrate.go` if needed
  - [ ] Verify: `go build ./...`, `go vet ./...` clean
  - [ ] Verify: unit test for the nil-DB path

- [ ] Task 4: `search_users` tool
  - [ ] `internal/mcp/tools/users.go` with `jsonschema`-tagged input/output structs
  - [ ] Register in `server.go` with scope `mcp:read`
  - [ ] Paginated email search; page size capped at 50; 90-day range cap
  - [ ] Missing actor → safe tool error
  - [ ] Verify: `go build ./...`, `go vet ./...`, `go test ./...` clean
  - [ ] Verify: unit tests for normal call, missing actor, page-size cap

### Checkpoint — Domain Tools

- [ ] `search_users` returns users as JSON through the handler
- [ ] Migration follows the existing pattern
- [ ] Build/vet/test pass
- [ ] Human reviews before bootstrap/routing work

## Phase 3 — Integration

- [ ] Task 5: `internal/bootstrap/mcp.go`
  - [ ] Build the handler from `Dependencies` (config, DB, user repo, actor, version)
  - [ ] No tool logic; no `internal/routes` import from bootstrap
  - [ ] Return nil handler when `MCP_ENABLED=false`
  - [ ] Verify: `go build ./...`, `go vet ./...` clean
  - [ ] Verify: `MCP_ENABLED=false` registers no route

- [ ] Task 6: `internal/routes/mcp.go` + mount `/api/mcp`
  - [ ] Mount via `adaptor.HTTPHandler` at `MCP_PATH`, methods GET/POST/DELETE
  - [ ] Register before `/api/v1/*`, `/build/*`, and the `/*` SPA fallback
  - [ ] Verify: `curl -i localhost:8080/api/mcp` → 401 without token
  - [ ] Verify: correct bearer token → non-401 (MCP response)
  - [ ] Verify: `MCP_ENABLED=false` → `/api/mcp` reaches the normal `/api/*` JSON 404
  - [ ] Verify: `/healthz`, `/api/v1/*`, `/build/*`, and `/` unchanged
  - [ ] Verify: `go build ./...`, `go vet ./...`, `go test ./...` clean

### Checkpoint — Complete

- [ ] All Task 1–6 acceptance criteria met
- [ ] `ping`, `whoami`, `search_users` respond on `/api/mcp`
- [ ] Auth, size limit, and origin rules covered by unit tests
- [ ] Route order correct; `/api/*` 404s unchanged
- [ ] No changes under `docs/`, `tests/`, or `.github/`; `.env` untouched
- [ ] Human approves; ready for `/build`

## Deferred

- [ ] Write tools (`mcp:write`)
- [ ] PostgreSQL `rate_limits` integration
- [ ] SSE timeout/proxy-buffering tuning
- [ ] Testcontainers integration + contract/golden tests
- [ ] Per-tool scope enforcement beyond `mcp:read`

## Standing Rules

- Never log or commit `MCP_BEARER_TOKEN`; only `.env.example` gets a placeholder.
- Keep auth on the `net/http` side; Fiber locals are unavailable after the adaptor.
- Register `/api/mcp` before the SPA fallback.
- Tools call the same use cases as HTTP handlers; no duplicated business rules.
- Audit is written by the shared wrapper, never inside a tool body.
- English comments; do not edit `docs/`; do not create `tests/` or `.github/`.
- Verify with `go build ./...` / `go vet ./...` / `go test ./...`; do not run the binary.
