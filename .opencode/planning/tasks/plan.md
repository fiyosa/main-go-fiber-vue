# Implementation Plan: MCP Server (Database Access)

## Overview

Add an **MCP server for direct PostgreSQL access** to the existing Go/Fiber app. The
server speaks the MCP **Streamable HTTP** transport on `/api/mcp` (same port as the HTTP API),
authorizes every request with a static **bearer token** (`MCP_BEARER_TOKEN`) compared via
SHA-256 `subtle.ConstantTimeCompare`, and exposes `ping`, `whoami`, and `run_sql`. The
`run_sql` tool executes raw SQL, including DDL and DML (user decision); nothing about tool
calls is persisted.

MCP is a delivery adapter like the HTTP routes: no separate process, port, or queue.

Out of scope: OAuth, per-token rate limiting on PostgreSQL, SSE optimization, and
testcontainers integration tests (no local Postgres).

## Architecture Decisions

- **SDK**: official `github.com/modelcontextprotocol/go-sdk/mcp`; server built with
  `mcp.NewServer`, tools registered with `mcp.AddTool` (input/output structs with `jsonschema`
  tags → auto JSON-schema). Verify exact constructor names against the pinned version at
  coding time (SDK still evolving).
- **Transport**: Streamable HTTP handler via `mcp.NewStreamableHTTPHandler(...)` → `http.Handler`,
  mounted on Fiber with `adaptor.HTTPHandler`.
- **Route order** (in `internal/routes`): `/healthz`,`/readyz` → `/api/mcp` → `/api/v1/*` →
  `/build/*` → `/*` SPA fallback. MCP must mount before the fallback.
- **Auth**: net/http-side wrapper `WithAuth` runs before the SDK handler. Bearer token →
  SHA-256 → `subtle.ConstantTimeCompare` against configured hash. Actor bound from
  `MCP_ACTOR_USER_ID` into context; `whoami`/`run_sql` fail when the actor is unset.
  Never log the token.
- **Stateless**: `MCP_STATELESS=true`, JSON response mode preferred; no server-side session
  state, no sticky-session dependency.
- **File layout** follows docs 19.6 but lands under `internal/mcp/` (existing `internal/mcp/`
  package with `tools/`); bootstrap+route in `internal/bootstrap/mcp.go` and
  `internal/routes/mcp.go` as documented.
- **Config**: add `MCP_ENABLED`, `MCP_PATH`, `MCP_BEARER_TOKEN`, `MCP_ACTOR_USER_ID`,
  `MCP_SCOPES`, `MCP_MAX_BODY_BYTES`, `MCP_ALLOWED_ORIGINS`, `MCP_SESSION_TIMEOUT`,
  `MCP_STATELESS`, `MCP_SQL_MAX_ROWS`, `MCP_SQL_TIMEOUT` to the envconfig config struct;
  `MCP_BEARER_TOKEN` required when `MCP_ENABLED=true` (startup fails if empty). Update
  `.env.example` only, never `.env`.
- **No audit persistence** (user decision): tool calls and query text are neither stored nor
  logged. No `mcp_audit_logs` table or migration is created.
- **`run_sql`**: raw SQL via `database/sql`, no allowlist. `RETURNING` statements are read as
  queries; other statements report `rowsAffected`. Result rows are capped by
  `MCP_SQL_MAX_ROWS` and each statement is bounded by `MCP_SQL_TIMEOUT`.
- **Scope**: `mcp:read` granted by default via `MCP_SCOPES`.

## Dependency Graph

```text
config MCP_* fields (internal/config/config.go) + .env.example
        │
        ▼
internal/mcp/{auth,server,handler}.go             tools/system.go (ping, whoami)
        │                                               │
        │        ┌──────────────────────────────────────┤
        ▼        ▼
internal/mcp/handler.go  (WithAuth → StreamableHTTP handler)
        │
        ▼
internal/bootstrap/mcp.go  (wires Dependencies → http.Handler; NO tool logic)
        │
        ▼
internal/routes/mcp.go  (adaptor.HTTPHandler → GET/POST /api/mcp)
        │
        ▼
internal/bootstrap/app.go  (RegisterMCP before RegisterWeb /**/ fallback)
        │
        ▼
tools/sql.go (run_sql) ──> DBExecutor ──> *sql.DB (raw statements)
```

## Task List

### Phase 1: Foundations — config, auth, system tools, server skeleton

- [x] Task 1: MCP config fields + `.env.example`
  - Description: Add `MCP_ENABLED` (default `true`), `MCP_PATH` (`/api/mcp`),
    `MCP_BEARER_TOKEN` (required when enabled), `MCP_ACTOR_USER_ID`, `MCP_SCOPES`
    (`mcp:read`), `MCP_MAX_BODY_BYTES` (`1048576`), `MCP_ALLOWED_ORIGINS`,
    `MCP_SESSION_TIMEOUT` (`60s`), `MCP_STATELESS` (`true`) to the envconfig config struct.
    Startup fails when `MCP_ENABLED=true` and the token is empty. Add an `MCP_*` block to
    `.env.example`. Extract `bearerHash` (SHA-256 of token) at load time.
  - Acceptance criteria:
    - [x] Config fields load from env with documented defaults.
    - [x] Constructing the MCP config with `MCP_ENABLED=true` and an empty token fails.
    - [x] `.env.example` documents every `MCP_*` variable; `.env` untouched.
    - [x] No MCP package yet imports the DB.
  - Verification:
    - [x] `go build ./...` and `go vet ./...` clean.
    - [x] Small unit test: defaults + empty-token failure (config package).
  - Dependencies: None.
  - Files likely touched:
    - `internal/config/config.go`, `internal/config/config_test.go`
    - `.env.example`
  - Estimated scope: Small (1–2 files + env block).

- [x] Task 2: MCP core files — auth, server, handler + system tools
  - Description: Create `internal/mcp/{auth,server,handler}.go`. `auth.go` parses the
    bearer token, stores `bearerHash`, and compares with `subtle.ConstantTimeCompare`.
    `server.go` builds `mcp.NewServer`, registers `ping`/`whoami` via `mcp.AddTool` with
    `jsonschema`-tagged input/output structs, holds a server version. `handler.go` wires
    `WithAuth` around the `mcp.NewStreamableHTTPHandler` result and enforces
    `MCP_MAX_BODY_BYTES` + `MCP_ALLOWED_ORIGINS` (OPTIONS/CORS).
    `tools/system.go` implements `ping` (version + time) and `whoami` (actor, scopes,
    version; fails when the actor is unset).
  - Acceptance criteria:
    - [x] Correct bearer token passes the wrapper; missing/malformed/wrong token → 401.
    - [x] `ping`/`whoami` are registered; `whoami` returns actor + scopes.
    - [x] Actor missing for `whoami` → tool error (safe message), not a panic.
    - [x] Oversized body rejected; disallowed origin rejected; correct CORS on OPTIONS.
    - [x] No raw SQL, file path, shell command, or URL accepted by any tool.
  - Verification:
    - [x] `go build ./...`, `go vet ./...`, `go test ./...` clean.
    - [x] Unit tests: auth compare, origin check, body size, whoami-without-actor.
  - Dependencies: Task 1.
  - Files likely touched:
    - `internal/mcp/auth.go`, `internal/mcp/server.go`, `internal/mcp/handler.go`
    - `internal/mcp/tools/system.go`
    - `internal/mcp/*_test.go`
    - `go.mod` (+ indirect `github.com/google/jsonschema-go/jsonschema`)
  - Estimated scope: Medium (5–7 files).

### Checkpoint: Foundations

- [x] `go build ./...`, `go vet ./...`, `go test ./...` pass.
- [x] Auth wrapper rejects missing/invalid tokens; system tools respond through the handler.
### Phase 2: Database Access — run_sql

- [x] Task 3: reverted — MCP uses raw SQL, so no user repository change
  - Description: `search_users` was dropped in favor of the raw `run_sql` tool. The
    `SearchUsers` addition and its migration/audit work were removed; the user repository
    keeps only `FindByEmail`.
  - Acceptance criteria:
    - [x] `UserRepository` unchanged apart from `FindByEmail`.
    - [x] No new migration and no `mcp_audit_logs` table.
  - Verification:
    - [x] `go build ./...`, `go vet ./...` clean.
    - [x] `go test ./...` passes.
  - Files touched:
    - `internal/app/application/core/repository/user/repository.go`
    - `internal/app/application/core/repository/user/repository_test.go`

- [x] Task 4: `run_sql` tool
  - Description: Create `internal/mcp/tools/sql.go` and `internal/mcp/sql_executor.go`.
    Input `{ sql: string }`; output `{ columns, rows, rowCount, rowsAffected, truncated,
    durationMs }`. Any statement is allowed, including DDL and DML (user decision). The
    executor uses `database/sql` directly, treats `RETURNING` statements as queries, caps
    rows at `MCP_SQL_MAX_ROWS`, and enforces `MCP_SQL_TIMEOUT`. Actor must be set. No audit
    or logging of query text.
  - Acceptance criteria:
    - [x] Tool registered in `server.go`; visible as `run_sql`.
    - [x] Row-returning and non-returning statements both return a structured result.
    - [x] `Rows` capped at `MCP_SQL_MAX_ROWS`; `truncated` set when cut.
    - [x] Missing actor, empty SQL, and nil database → safe tool errors.
    - [x] No query text is persisted or logged.
  - Verification:
    - [x] `go build ./...`, `go vet ./...`, `go test ./...` clean.
    - [x] Unit tests: statement classification, value normalization, tool call via SDK.
  - Dependencies: Task 2.
  - Files:
    - `internal/mcp/tools/sql.go`
    - `internal/mcp/sql_executor.go`
    - `internal/mcp/server.go`, `internal/mcp/handler.go`
    - `internal/config/config.go`, `.env.example`

### Checkpoint: Database Tools

- [x] `run_sql` works through the handler with a fake executor.
- [x] No new database table or migration from MCP.
- [x] `go build ./...`, `go vet ./...`, `go test ./...` pass.

### Phase 3: Integration — bootstrap wiring + route registration

- [ ] Task 5: `internal/bootstrap/mcp.go` (dependency wiring only)
  - Description: Build the MCP handler from `Dependencies` — config, `*gorm.DB`,
    user repository, actor, version. Returns an `http.Handler` (already wrapped with auth +
    limits). Contains **no tool logic**. Called from `internal/bootstrap/app.go`
    before the SPA fallback is registered.
  - Acceptance criteria:
    - [ ] Returns `nil, nil` (or skips registration) when `MCP_ENABLED=false`.
    - [ ] Startup errors when enabled and token empty (delegated to config).
    - [ ] No import of `internal/routes` from bootstrap (dependency direction preserved).
  - Verification:
    - [ ] `go build ./...`, `go vet ./...` clean.
    - [ ] `MCP_ENABLED=false` → handler nil and no route registered.
  - Dependencies: Tasks 2, 4.
  - Files likely touched:
    - `internal/bootstrap/mcp.go`
    - `internal/bootstrap/app.go`, `internal/bootstrap/wire_api.go`
  - Estimated scope: Small (1–3 files).

- [ ] Task 6: `internal/routes/mcp.go` + mount `/api/mcp`
  - Description: Register the MCP handler through `adaptor.HTTPHandler` at `MCP_PATH`
    (`/api/mcp`) for the methods the Streamable HTTP transport uses (GET/POST/DELETE).
    Mount **before** `/api/v1/*`, `/build/*`, and the `/*` SPA fallback, per the documented
    route order (health → mcp → api → build → spa). Ensure `/api/mcp` is never swallowed by
    the SPA fallback, and unknown `/api/*` still returns JSON 404.
  - Acceptance criteria:
    - [ ] `GET /api/mcp` without token → 401; with correct token → MCP response.
    - [ ] `POST /api/mcp` (initialize + tools/list) works with the bearer token.
    - [ ] Missing token → 401; wrong token → 401; oversized body → rejected.
    - [ ] `MCP_ENABLED=false` → `/api/mcp` falls through to the normal `/api/*` JSON 404.
    - [ ] `/healthz`, `/api/v1/*`, `/build/*`, and the SPA at `/` are unchanged.
  - Verification:
    - [ ] `go build ./...`, `go vet ./...`, `go test ./...` clean.
    - [ ] `curl -i localhost:8080/api/mcp` → 401 without token.
    - [ ] `curl -i -H "Authorization: Bearer <token>" localhost:8080/api/mcp` → non-401.
    - [ ] `curl -i localhost:8080/api/v1/nope` → JSON 404.
    - [ ] `curl -s localhost:8080/some/unknown/path | head` → SPA shell.
  - Dependencies: Task 5.
  - Files likely touched:
    - `internal/routes/mcp.go`
    - `internal/routes/api.go` or `internal/bootstrap/app.go` (registration order)
    - `internal/routes/mcp_test.go`
  - Estimated scope: Small (2–3 files).

### Checkpoint: Complete

- [ ] All six tasks' acceptance criteria met.
- [ ] `/api/mcp` reachable on the main port; tools `ping`, `whoami`, `run_sql` respond.
- [ ] Auth, size limit, and origin rules enforced by unit tests.
- [ ] Route order correct: MCP before SPA fallback; `/api/*` 404s unchanged.
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` clean.
- [ ] No changes under `docs/`, `tests/`, or `.github/`; `.env` untouched.
- [ ] Human reviews before `/build` starts.

## Risks and Mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| MCP Go SDK API names drift (still evolving, per docs 19.2) | High | Pin a specific version; verify `mcp.NewServer`, `mcp.AddTool`, `mcp.NewStreamableHTTPHandler` signatures before coding |
| Fiber `adaptor.HTTPHandler` rebuilds the request, so Fiber locals/token are unavailable | High | Do auth and actor injection on the `net/http` side in `WithAuth`, as docs 19.3 instructs |
| `/api/mcp` swallowed by `/*` SPA fallback | High | Register MCP before `RegisterWeb`; add a route-order test |
| No local Postgres → `run_sql` untestable end to end | Medium | Executor is injected; unit tests use a fake executor and classify/normalize helpers; real DB test deferred |
| SSE clients hit the global Fiber write timeout | Medium | Prefer stateless JSON response mode; document the SSE timeout exemption and proxy buffering setting |
| Token accidentally logged or committed | High | Never log the token; only `.env.example` gets a placeholder; hash stored, raw token not kept |
| `MCP_ACTOR_USER_ID` unset in production | Medium | `run_sql` fails with a safe tool error; document the required variable |
| `run_sql` can drop or alter data because any statement is allowed | High | Accepted by user; DB_* user should be a dedicated role, backups required, token treated as full DB access |

## Open Questions

1. **MCP file location**: docs 19.6 shows `internal/app/infrastructure/mcp/`, but this repo
   already has `internal/mcp/tools/`. Plan targets `internal/mcp/` to match existing code.
   Confirm this is acceptable (or move to the documented path).
2. **v1 tools**: user chose raw SQL (`run_sql`) instead of `search_users`; confirm DDL stays
   enabled in production.
3. **Audit storage**: user decided MCP must not create tables or persist tool calls. No
   `mcp_audit_logs` table; confirm again if audit becomes required.
4. **Actor identity**: `MCP_ACTOR_USER_ID` must point to an existing user in production.
   Plan assumes the operator sets it; no seeding is done.
5. **Rate limiting**: docs 19.3 mentions the PostgreSQL `rate_limits` table keyed by token
   fingerprint + IP. Plan defers it (no DB locally). Confirm deferral.

## Parallelization Opportunities

- Task 3 (repo method) can proceed in parallel with Task 4's tool skeleton, then
  Task 4 finishes against the repo interface.
- Tasks 5 and 6 are tightly coupled; do them together, and do not change routing in parallel
  with other `/api/*` work.
- Unit tests for Tasks 2 and 4 can be written alongside their implementations.


