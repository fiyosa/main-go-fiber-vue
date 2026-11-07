# Implementation Plan: MCP Server (Read-Only)

## Overview

Add a **read-only MCP server** to the existing Go/Fiber app per `docs/19-mcp-server.md`. The
server speaks the MCP **Streamable HTTP** transport on `/api/mcp` (same port as the HTTP API),
authorizes every request with a static **bearer token** (`MCP_BEARER_TOKEN`) compared via
SHA-256 `subtle.ConstantTimeCompare`, and exposes three tools: `ping`, `whoami`,
`search_users`. Tool calls are not persisted; MCP writes nothing to the database.

MCP is a delivery adapter like the HTTP routes: no separate process, port, or queue.

Out of scope: mutating/write tools, OAuth, per-token rate limiting on PostgreSQL, SSE
optimization, and testcontainers integration tests (no local Postgres). Write tools and rate
limits are deferred until the read-only surface is stable and a DB is available.

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
  `MCP_ACTOR_USER_ID` into context; `whoami`/`search_users` fail when the actor is unset.
  Never log the token.
- **Stateless**: `MCP_STATELESS=true`, JSON response mode preferred; no server-side session
  state, no sticky-session dependency.
- **File layout** follows docs 19.6 but lands under `internal/mcp/` (existing `internal/mcp/`
  package with `tools/`); bootstrap+route in `internal/bootstrap/mcp.go` and
  `internal/routes/mcp.go` as documented.
- **Config**: add `MCP_ENABLED`, `MCP_PATH`, `MCP_BEARER_TOKEN`, `MCP_ACTOR_USER_ID`,
  `MCP_SCOPES`, `MCP_MAX_BODY_BYTES`, `MCP_ALLOWED_ORIGINS`, `MCP_SESSION_TIMEOUT`,
  `MCP_STATELESS` to the envconfig config struct; `MCP_BEARER_TOKEN` required when
  `MCP_ENABLED=true` (startup fails if empty). Update `.env.example` only, never `.env`.
- **No audit persistence** (user decision): tool calls are not stored in the database and no
  `mcp_audit_logs` table or migration is created. The read-only tools only read the users
  table. DB remains optional at runtime.
- **`search_users`**: paginated user search by email, max page size 50 — added to the user
  repository interface and implemented for GORM. Verified at code time; GORM user repo
  nil-safe.
- **Scope**: `mcp:read` granted by default via `MCP_SCOPES`; `mcp:write` reserved for future
  mutating tools.

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
        ▼tools/users.go (search_users) ──> user repository interface + GORM impl (nil-safe)
                                       │
                                       ▼
                              wire_api.go / bootstrap DI (Dependencies)
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
### Phase 2: Domain Tools — search_users

- [x] Task 3: user repository `SearchUsers` (no MCP table, no migration)
  - Description: Add `SearchUsers(ctx, email, page, pageSize)` to the `UserRepository`
    interface and implement it in the GORM user repository (nil-safe, max 50 rows,
    ILIKE %email% query). No migration and no audit table: the user decided MCP must not
    create database tables.
  - Acceptance criteria:
    - [x] `UserRepository.SearchUsers` returns at most 50 rows; nil DB → empty slice, no panic.
    - [x] `SearchUsers` interface added without breaking existing callers.
    - [x] No new migration and no `mcp_audit_logs` table.
  - Verification:
    - [x] `go build ./...`, `go vet ./...` clean.
    - [x] Unit test for nil-DB path in `SearchUsers`.
  - Dependencies: Task 2 (tools need the repo interface to exist).
  - Files likely touched:
    - `internal/app/application/core/repository/user/repository.go`
    - `internal/app/application/core/repository/user/repository_test.go`
  - Estimated scope: Small (2 files).

- [ ] Task 4: `search_users` tool
  - Description: Create `internal/mcp/tools/users.go`. Tool input: `{ email: string,
    page: int }`; output: `{ users: [], total: int, page: int, pageSize: int }`. Calls
    `UserRepository.SearchUsers`; enforces max page size 50 and max
    90-day range (created_at filter). Actor must be set or the tool returns a safe error.
    No database writes: MCP does not persist tool calls.
  - Acceptance criteria:
    - [ ] Tool registered in `server.go` with scope `mcp:read`.
    - [ ] Returns paginated users matching email; empty list when no match.
    - [ ] Page size capped at 50; actor missing → safe tool error.
    - [ ] No raw SQL, table name, or file path accepted.
    - [ ] No database write happens during a tool call.
  - Verification:
    - [ ] `go build ./...`, `go vet ./...`, `go test ./...` clean.
    - [ ] Unit tests: normal call, missing actor, page-size cap.
  - Dependencies: Task 3.
  - Files likely touched:
    - `internal/mcp/tools/users.go`
    - `internal/mcp/server.go` (register the tool)
  - Estimated scope: Small (2 files).

### Checkpoint: Domain Tools

- [ ] `search_users` works through the handler; returns users as JSON.
- [ ] No new database table or migration from MCP.
- [ ] `go build ./...`, `go vet ./...`, `go test ./...` pass.
- [ ] Human reviews before bootstrap/routing work.

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
- [ ] `/api/mcp` reachable on the main port; tools `ping`, `whoami`, `search_users` respond.
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
| No local Postgres → `search_users` untestable end to end | Medium | Nil-DB paths return safe defaults; unit tests cover the logic; real DB test deferred |
| SSE clients hit the global Fiber write timeout | Medium | Prefer stateless JSON response mode; document the SSE timeout exemption and proxy buffering setting |
| Token accidentally logged or committed | High | Never log the token; only `.env.example` gets a placeholder; hash stored, raw token not kept |
| `MCP_ACTOR_USER_ID` unset in production | Medium | `whoami`/`search_users` fail with a safe tool error; document the required variable |

## Open Questions

1. **MCP file location**: docs 19.6 shows `internal/app/infrastructure/mcp/`, but this repo
   already has `internal/mcp/tools/`. Plan targets `internal/mcp/` to match existing code.
   Confirm this is acceptable (or move to the documented path).
2. **v1 tools**: Should `search_users` filter by role/status too, or email only? (Plan assumes
   email only, matching docs 19.5.)
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


