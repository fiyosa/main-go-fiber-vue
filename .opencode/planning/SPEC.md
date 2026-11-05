# Spec: Alma Backend Rewrite — Phase 0 Foundation (Scaffold-First)

Status: Draft
Date: 2026-11-02
Source of truth: `docs/01`–`docs/19` and appendices
Scope: Phase 0 only. Later phases get their own specs.
Current increment: Slice 1 — one public login API (`POST /api/v1/login`) using only `users`
and `auths`; no RBAC, LDAP, 2FA, or SSO.

## 1. Objective

Replace the Laravel backend of Alma with a Go + Fiber service while preserving the existing
Vue 3 SPA contract at `/api/v1`. The rewrite is parity-first: paths, payloads, status codes,
response envelopes, pagination, permission strings, and file URLs must not change unless the
difference is recorded in `docs/18-risks-and-open-decisions.md`.

Target users:

- Existing Alma SPA users (admins and employees).
- Remote MCP automation clients (Phase 6, out of scope here).

Phase 0 goal: build the repository skeleton and an observable foundation so later modules can
be added incrementally. The first deliverable is the project folder tree plus `go.mod`.

Why:

- A new schema is allowed, so no legacy data migration is required.
- A single-process Go service (API + scheduler) replaces Laravel Passport, queue, and cron.
- Contract tests plus route parity are the confidence mechanism before cutover.

Success (Phase 0):

- Folder tree matches the documented spine.
- `go build ./...`, `go vet ./...`, and `go test ./...` pass once code exists.
- `GET /healthz` and `GET /readyz` work.
- Response envelope and pagination shapes match Laravel.
- One sample module has a contract test (when `tests/` is enabled).
- CI gates (lint, test, build, govulncheck) are in place (when `.github/` is enabled).

## 2. Tech Stack

| Concern | Choice | Notes |
| --- | --- | --- |
| Language | Go 1.27.3 | Installed locally |
| HTTP | `github.com/gofiber/fiber/v3` | Pin exact latest stable at Task 2 |
| Config | `github.com/kelseyhightower/envconfig` | 12-factor env only |
| Logging | `log/slog` | JSON in production, text in dev |
| ORM | `gorm.io/gorm` + `gorm.io/driver/postgres` | Explicit tags; no AutoMigrate in prod |
| Driver | `github.com/jackc/pgx/v5` | PostgreSQL only |
| Migrations | `github.com/golang-migrate/migrate/v4` | SQL-first; `console migrate` |
| Jobs | `internal/platform/jobs` | Goroutines; panic recovery; `jobs.Wait` |
| Scheduler | `github.com/robfig/cron/v3` | Stub in Phase 0; advisory lock later |
| JWT | `github.com/golang-jwt/jwt/v5` | HS256 via `JWT_SECRET`, 3h access TTL |
| Validation | `github.com/go-playground/validator/v10` | Laravel-shaped errors |
| CLI | `github.com/spf13/cobra` | `app` and `console` |
| Concurrency | `golang.org/x/sync/errgroup` | Shared context lifecycle |

Explicitly excluded for now: SSO, PDF generation, LDAP login, 2FA, MCP tools, queue/broker.
Postgres wiring is deferred until a DSN is available (no Docker or psql installed locally).

## 3. Commands

Local toolchain: Go 1.27.3, Node 24, npm 11, git 2.55. No Docker, psql, make, or golangci-lint.

```bash
# Module
go mod init go-fiber-vue
go mod tidy

# Build and vet
go build ./...
go vet ./...

# Run
go run ./cmd/app
go run ./cmd/app --disable-scheduler
go run ./cmd/console migrate status

# Test (once tests/ is enabled)
go test ./... -race -count=1
go test ./tests/unit/... -race -count=1
go test ./tests/integration/... -race -count=1 -tags=integration

# Quality gates (install once, once CI is enabled)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
golangci-lint run ./...
govulncheck ./...

# Frontend (only when web/ source exists)
npm --prefix web run dev
npm --prefix web run build
```

Notes:

- No `make`; all commands are direct Go/npm commands.
- Integration tests are skipped until Docker or a Postgres DSN is available.
- The `app` process never runs migrations; `console migrate` is an init job.

## 4. Project Structure

Phase 0 creates the spine plus the `core` context. Other contexts are added when their phase
starts. `tests/` and `.github/` are intentionally deferred for now.

```text
go-fiber-vue/                     # module: go-fiber-vue
├── cmd/
│   ├── app/                      # API + scheduler process
│   └── console/                  # one-off commands: migrate, seed, imports, syncs
├── internal/
│   ├── app/
│   │   ├── domain/
│   │   │   └── core/{entity,enum,policy}
│   │   ├── application/          # lowercase per naming rules
│   │   │   └── core/{usecase,repository,criteria,jobs,export,import,command}
│   │   ├── infrastructure/
│   │   │   ├── scheduler/
│   │   │   ├── integrations/     # per-integration folders added later
│   │   │   ├── rules/
│   │   │   ├── services/approval/
│   │   │   └── support/{helpers,traits}
│   │   └── presentation/
│   │       ├── core/{controllers,request,response,routes,rule,view}
│   │       ├── public/
│   │       └── apiv1/
│   ├── bootstrap/                # app, database, cache, runtime, scheduler, mcp
│   ├── mcp/tools/
│   ├── routes/{api.go,web.go,mcp.go} and routes/v1/apiv1.go
│   ├── config/
│   └── platform/{auth,database,cache,jobs,storage,broadcast}
├── migrations/
├── web/                          # Vue source placeholder; not in repo yet
├── public/build/                 # Vite build output
├── tests/                        # Deferred: not created yet
├── docs/                         # existing planning set; read-only
├── .github/                      # Deferred: CI not created yet
├── go.mod
└── README.md
```

Deliberate deviations from `docs/05-repository-layout.md`:

- `internal/app/Application` → `internal/app/application` (lowercase, per docs 5.3).
- `internal/app/infrastructure/external/` is omitted because SSO is not ported.
- Context folders other than `core` are added in their own phase, not as empty placeholders.
- `tests/` and `.github/` are deferred until explicitly enabled.

## 5. Code Style

One real snippet (the envelope is the contract):

```go
package helpers

import "github.com/gofiber/fiber/v3"

type Envelope struct {
    Success bool   `json:"success"`
    Data    any    `json:"data,omitempty"`
    Message string `json:"message,omitempty"`
}

func Success(c fiber.Ctx, data any, message string, code int) error {
    return c.Status(code).JSON(Envelope{Success: true, Data: data, Message: message})
}
```

Conventions:

- Directories: lowercase, singular concept (`controllers`, `usecase`, `entity`).
- Packages: lowercase, one word (`repository`, `controllers`).
- Files: `snake_case.go` (`user_controller.go`).
- Types: PascalCase. Functions: camelCase, verb-first (`FindByUsername`).
- Interfaces: noun + `Repository` / `Service`, defined in the consuming layer.
- Request DTOs: `<Action><Entity>Request`. Responses: `<Entity>Response`.
- Job functions: `<context>:<domain>:<action>` (`core:telegram:send`).
- Always `gofmt`/`goimports`; wrap errors with `%w`; `context.Context` is the first param.
- GORM models use explicit `TableName()` and column tags; never AutoMigrate in production.
- Store timestamps in UTC; render in `Asia/Jakarta`.
- Money uses `decimal.Decimal` or integer minor units, never `float64`.
- Comments, when present, are English and minimal.

Compatibility quirks to preserve:

- `sendResponse` returns `success: true` even on `403`.
- `sendError` may downgrade to `404` only where the frontend relies on it; otherwise document it.
- Pagination keeps the Laravel paginator JSON shape (`data`, `links`, `meta`, `current_page`, ...).

## 6. Testing Strategy

| Level | Scope | Tooling | Location |
| --- | --- | --- | --- |
| Unit | Envelope, mappers, pagination, config, token service | `testing`, `testify` | `tests/unit` (deferred) |
| Integration | Handlers + Postgres | `testcontainers-go`, `httptest` | `tests/integration` (deferred) |
| Contract | JSON shape/status vs Laravel golden files | Golden JSON | `tests/contract` (deferred) |
| E2E | Login → API → logout | Testcontainers + Fiber app | `tests/integration` (deferred) |
| Load | Login, check-in, lists | k6/vegeta | later |

Until `tests/` is enabled, verification uses package-local `_test.go` files and manual checks.

Phase 0 required tests (when enabled):

- Envelope success/error/validation shapes and pagination keys.
- `GET /healthz` and `GET /readyz` status and JSON.
- Route registry test: every registered non-public route has an access category.
- Sample module contract golden test.
- Config parsing from env; startup fails on missing required values.

Coverage targets:

- `internal/platform/auth`: ≥ 85%.
- Helpers, mappers, response envelope: ≥ 80%.
- No global coverage gate until the module surface stabilizes.

Blockers:

- Integration and E2E need Docker (not installed) or a Postgres DSN.
- Contract parity needs Laravel response captures (Laravel source not in this workspace).
- `tests/` is deferred by request; use `testing.Short()` or the `integration` build tag later.

## 7. Boundaries

Always:

- Keep `/api/v1`, envelopes, status codes, permission strings, and file URLs compatible.
- Read secrets from the environment only; never log tokens, passwords, or raw MCP arguments.
- Keep layering: `internal/app` for application code, `internal/platform` for runtime adapters.
- Register routes explicitly; no file scanning.
- Run `gofmt`, `go vet ./...`, `go build ./...`, and `go test ./...` before handoff.
- Write docs and code comments in English.
- Record every intentional behavior difference in `docs/18-risks-and-open-decisions.md`.
- Use parameterized queries and Postgres-compatible SQL.

Ask first:

- Adding any dependency outside the approved list.
- Any schema change, new table, or column rename.
- Changing an API path, envelope, status code, or permission string.
- Enabling SSO, PDF, LDAP login, 2FA, MCP write tools, or a queue/broker.
- Changing CI configuration or splitting the scheduler into a separate process.
- Editing the `docs/` planning set.
- Creating `tests/` or `.github/`.

Never:

- Commit `.env` files or secrets.
- Use GORM `AutoMigrate` in production or migrate from the `app` process.
- Port the LDAP credential check or the second-factor flow.
- Rename public API paths, database tables/columns, permission strings, or the Alma brand.
- Use `float64` for money.
- Add the global auth middleware twice.
- Invent permission strings.
- Delete or skip failing tests to make CI green.
- Use a non-Postgres database for tests.

## 8. Milestones

- M0.1 Project scaffold — DONE
  - `go.mod` (`module go-fiber-vue`, Go 1.27.3) and the folder spine.
  - `tests/` and `.github/` are deferred.
- M0.2 Bootable process
  - `cmd/app`, `cmd/console`, config, slog, errgroup lifecycle.
  - Verify: `go build ./...`; clean start and Ctrl+C exit.
- M0.3 HTTP health slice
  - Fiber v3, request ID middleware, `/healthz`, `/readyz`, graceful HTTP shutdown.
  - Verify: manual `curl` + handler tests.
- M0.4 Response contract primitives
  - Envelope, validation error mapper, pagination.
  - Verify: package-local tests.
- M0.5 Jobs and scheduler slice
  - `jobs.Go` with panic recovery/timeouts/`jobs.Wait`; cron stub; `--disable-scheduler`.
  - Verify: a debug job runs and drains within a bounded timeout.
- M0.6 Sample read-only module
  - Route → controller → usecase → stub repository for `GET /api/v1/settings`.
  - Verify: envelope JSON and explicit route registration; no DB required.
- M0.7 Database and migrations — deferred until DSN
  - Connection resolver, `console migrate`/`seed` skeleton, first SQL migration.
- M0.8 Tests and CI — deferred
  - Create `tests/` and `.github/`, then add gates.

## 9. Success Criteria

- M0.1 tree matches section 4; `docs/` unchanged; `tests/` and `.github/` absent.
- `go build ./...` and `go vet ./...` exit 0 from M0.2 onward.
- `GET /healthz` returns 200; `GET /readyz` returns 503 without a DB and 200 with one.
- SIGINT/SIGTERM stops the scheduler first and drains jobs within a bounded timeout.
- Envelope and pagination JSON keys match the Laravel contract (golden tests when enabled).
- `GET /api/v1/settings` returns the expected envelope through the full layer stack.
- No SSO, PDF, LDAP-login, or 2FA code exists.

## 10. Open Questions

1. Exact Fiber v3 version to pin (resolve at M0.3 after a compile smoke test).
2. JWT secret rotation: one `JWT_SECRET` now; key rotation and `kid` are out of scope.
3. Retention defaults: `activity_log` 90 days, `auths` pruned 30 days after `refresh_expires_at`.
4. Frontend source and `public/build` are not in this workspace; how will they be delivered?
5. MCP stays disabled until Phase 6; confirm no `/api/mcp` route in Phase 0.
6. `docs/06-layer-mapping.md` import examples omit the `internal/` prefix; section 4 here is authoritative.
7. Integration tests are blocked by missing Docker; DB tests are blocked by missing DSN.
8. Root `README.md` ownership and content.

## 11. Resolved Decisions

- Spec scope: Phase 0 foundation, phased; start with `go.mod` and folders.
- Laravel source: docs only for now; parity extraction is deferred.
- Module path: `go-fiber-vue`.
- Fiber: v3 latest stable, pinned.
- JWT: HS256 with `JWT_SECRET`.
- SSO: disabled / not ported.
- PDF: not ported now.
- Folder casing: `application` lowercase.
- Database: wiring deferred until a DSN is provided.
- `tests/` and `.github/`: deferred; do not create yet.
- Language: English docs and comments.
