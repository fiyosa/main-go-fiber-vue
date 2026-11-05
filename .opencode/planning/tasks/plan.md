# Implementation Plan: Slice 1 — Minimal Login API

## Overview

This plan supersedes the Phase 0 ordering in `SPEC.md`. The first deliverable is one public
endpoint, `POST /api/v1/auth/login`, backed only by the `users` and `auths` tables. There is no RBAC,
no roles/permissions/menus, no LDAP, no 2FA, and no SSO.

Login verifies the local `users.password` hash (bcrypt cost 12), then issues:

- an HS256 JWT access token (`JWT_ACCESS_TTL`, default 3h) with `iss`, `aud`, `sub`, `email`,
  `jti`, `iat`, `exp`;
- an opaque refresh token (`JWT_REFRESH_TTL`, default 30d), returned once and stored as a
  SHA-256 hex hash in `auths.refresh_token`.

The token response shape stays compatible with the SPA: top-level `token`, `token_type`,
`expires_in`, and `refresh_token`. The rest of Phase 0 (health, jobs, sample module) is reduced
to only what this slice needs.

Current state: `go.mod` exists (`module go-fiber-vue`, Go 1.27.3) and the folder spine exists.
No Go source files exist yet. `tests/` and `.github/` remain intentionally absent.

## Architecture Decisions

- New PostgreSQL schema; only `users` (`id`, `email`, `password`, timestamps) and `auths` tables
  (plus indexes). No RBAC tables.
- Revocation uses `auths.revoked_at`; there is no `users.token_version` column.
- Refresh tokens: 256-bit random, base64url; only the hex SHA-256 is stored.
- One `auths` row per login: `token` (JWT), `refresh_token` (hash), expiries, `meta`, `revoked_at`.
- Login is a public route. No auth or permission middleware is attached to it.
- JWT signing: HS256 with a shared `JWT_SECRET`. Issuer (`app-api`) and audience
  (`app-spa`) are fixed in code, not environment variables.
  Never commit keys.
- The login endpoint does not require `roles`/`permissions`; `GET /me` is out of scope for this
  slice.
- A minimal seeder creates one user from defaults in code (`admin@example.com`).
- `tests/` is deferred; use package-local `_test.go` files and manual `curl` until enabled.
- Postgres is required to run the endpoint. If no DSN is available, verification stops at
  build + unit tests with fake repositories.

## Dependency Graph

```text
config (env) ──> DB pool ──> migrations (users, auths) ──> seed user
     │                                                        │
     ├──> JWT + refresh token service ──┐                     │
     │                                  │                     │
     └──> bcrypt password verifier ─────┼──> login usecase <──┘ (user lookup by email)
                                        │         │
                                        │         └──> auths insert (refresh hash)
                                        │                     │
                                        └──> controller ──> route POST /api/v1/auth/login
                                                              │
                                                              └──> Fiber runtime
```

## Task List

### Phase 1: Runtime and Data Foundation

- [ ] Task 1: Minimal bootable runtime, config, and DB pool
  - Description: Add `cmd/app` and `cmd/console`, typed env config, slog, an `errgroup` runtime,
    a Fiber server with `GET /healthz` and `GET /readyz`, and a GORM/pgx PostgreSQL pool.
  - Acceptance:
    - [ ] `go build ./...` and `go vet ./...` pass.
    - [ ] `GET /healthz` returns 200.
    - [ ] `GET /readyz` returns 200 with a reachable DB and 503 JSON without one.
    - [ ] SIGINT/SIGTERM exits cleanly.
  - Verify: build/vet; `go run ./cmd/app`; `curl -i localhost:8080/healthz`; Ctrl+C.
  - Dependencies: none.
  - Files: `cmd/app/main.go`, `cmd/console/main.go`, `internal/config/config.go`,
    `internal/bootstrap/app.go`, `internal/bootstrap/runtime.go`,
    `internal/bootstrap/database.go`, `internal/platform/database/database.go`,
    `internal/routes/api.go`, `internal/app/presentation/core/routes/health_routes.go`.
  - Scope: Medium (5–8 files).

- [ ] Task 2: Schema and seed user
  - Description: Add SQL migrations for `users` and `auths`, `console migrate`/`migrate status`,
    and an idempotent admin seeder that hashes the password with bcrypt cost 12.
  - Acceptance:
    - [ ] Up/down migrations create and drop `users` and `auths` with the documented indexes.
    - [ ] `users` stores only `id`, `email`, `password`, and timestamps.
    - [ ] `console migrate status` reports the current version.
    - [ ] `console db seed` creates exactly one admin and refuses empty/short passwords.
  - Verify: run against a Postgres DSN; inspect `users` and `auths`; rerun seed (no duplicates).
  - Dependencies: Task 1.
  - Files: `migrations/000001_create_users_table.up.sql` (+ `.down.sql`),
    `migrations/000002_create_auths_table.up.sql` (+ `.down.sql`),
    `internal/platform/database/migrate.go`, `internal/platform/database/seed.go`,
    `cmd/console/main.go`.
  - Scope: Medium (5–6 files).

### Checkpoint A (after Tasks 1–2)

- [ ] Build and vet pass.
- [ ] Migrations apply and roll back cleanly.
- [ ] One seeded user exists; `auths` is empty until login.
- [ ] Human reviews the schema before auth code.

### Phase 2: Auth Primitives and Use Case

- [ ] Task 3: Password, JWT, and refresh-token services
  - Description: Add a bcrypt verifier (cost 12), an HS256 JWT issuer/verifier with the claims
    above, and a refresh-token generator that returns raw + SHA-256 hash.
  - Acceptance:
    - [ ] Correct password verifies; wrong password fails.
    - [ ] JWT round-trips; wrong `aud`/`iss`/expired token is rejected.
    - [ ] Refresh token is 256-bit and its hash is stable.
  - Verify: package-local `_test.go` with a fixed test secret.
  - Dependencies: Task 1.
  - Files: `internal/platform/auth/password.go`, `internal/platform/auth/jwt.go`,
    `internal/platform/auth/refresh.go`, `internal/platform/auth/*_test.go`.
  - Scope: Medium (4–5 files).

- [ ] Task 4: Login use case with repositories
  - Description: Add `User` and `Auth` entities, `UserRepository` and `AuthRepository`
    (interfaces + GORM implementations), and a `LoginUseCase` that validates credentials,
    issues tokens, and inserts an `auths` row.
  - Acceptance:
    - [ ] Valid credentials return access + refresh tokens and create one `auths` row.
    - [ ] Unknown user or wrong password returns `ErrInvalidCredentials`.
    - [ ] The JWT contains the user's `email` and no permission/scope claims.
    - [ ] No LDAP or 2FA branch exists.
  - Verify: table-driven use-case tests with fake repositories.
  - Dependencies: Tasks 2, 3.
  - Files: `internal/app/domain/core/entity/user.go`,
    `internal/app/domain/core/entity/auth.go`,
    `internal/app/application/core/repository/user_repository.go`,
    `internal/app/application/core/repository/auth_repository.go`,
    `internal/app/application/core/usecase/login_usecase.go`,
    `internal/app/application/core/usecase/login_usecase_test.go`.
  - Scope: Medium (5–6 files).

### Checkpoint B (after Tasks 3–4)

- [ ] Auth primitives and use-case tests pass.
- [ ] No RBAC or permission code exists.
- [ ] Human reviews the token claims and `auths` row shape.

### Phase 3: HTTP Login Slice

- [ ] Task 5: `POST /api/v1/auth/login` end-to-end
  - Description: Add the login request DTO, controller, response DTO, and public route; wire the
    real GORM repositories and the login use case through bootstrap.
  - Acceptance:
    - [ ] Valid login returns 200 with `success`, `message`, `token`, `token_type: "Bearer"`,
      `expires_in`, `refresh_token`.
    - [ ] Invalid credentials return 401 `{success:false, message:"Invalid credentials."}`.
    - [ ] Missing fields return 422 with field errors.
    - [ ] The route is public and has no permission middleware.
    - [ ] An `auths` row is written.
  - Verify: `curl -i -X POST localhost:8080/api/v1/auth/login` with valid and invalid bodies;
    focused package tests.
  - Dependencies: Task 4.
  - Files: `internal/app/presentation/core/request/login_request.go`,
    `internal/app/presentation/core/controllers/auth_controller.go`,
    `internal/app/presentation/core/response/login_response.go`,
    `internal/app/presentation/core/routes/auth_routes.go`,
    `internal/routes/v1/apiv1.go`, `internal/bootstrap/app.go`.
  - Scope: Medium (5–6 files).

### Phase 4: Localization Contract (docs/20)

- [x] Task 6: Translated response messages
  - Description: Copy the reviewed catalogues from `docs/lang/` to `lang/`, add the runtime
    translator, the locale middleware, and make every response helper accept a translation
    key instead of a literal.
  - Acceptance:
    - [ ] `internal/lang/en.json` and `internal/lang/id.json` have identical key sets (enforced at load).
    - [x] `?lang=` beats `Accept-Language`, which beats `APP_LOCALE`; unsupported values fall
      back to `en`.
    - [x] `APP_LOCALE` must be `en` or `id`, otherwise startup fails.
    - [x] `message` is present and non-empty on success, validation, auth, and 404/500 errors.
    - [x] Validation field messages use `validation.*` keys with `:attribute`.
  - Verify: `go test ./...`; manual `curl` with `Accept-Language: id`.
  - Dependencies: Task 1.
  - Files: `internal/lang/en.json`, `internal/lang/id.json`, `internal/platform/lang/translator.go`,
    `internal/app/presentation/core/middleware/locale.go`,
    `internal/app/infrastructure/support/helpers/response.go`,
    `internal/app/infrastructure/support/helpers/validation.go`, `internal/bootstrap/app.go`.
  - Scope: Medium (8 files incl. tests).

### Checkpoint C (after Tasks 5–6)

- [ ] Login works end-to-end against Postgres.
- [x] `go build ./...`, `go vet ./...`, and `go test ./...` pass.
- [x] Every response path returns a translated, non-empty `message`.
- [ ] Human confirms the response shape before any other endpoint is added.

## Risks and Mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| No Postgres/Docker locally | High | Migrations + use-case tests first; manual login needs a DSN |
| JWT secret handling | High | Single `JWT_SECRET`; long random value; never commit it |
| `tests/` deferred | Medium | Package-local tests + curl; migrate later |
| Token response parity | Medium | Match docs/08 exactly; golden test when `tests/` exists |
| `internal/lang/` missing at runtime | High | Startup fails with an explicit error; deployment must copy `lang/` |
| Seeder password safety | Medium | Require ≥12 chars; reject known defaults; bcrypt cost 12 |
| bcrypt cost 12 latency | Low | Accept for login; do not raise without measurement |
| Single-process shutdown | Low | `errgroup` + bounded shutdown from Task 1 |

## Parallelization Opportunities

- Task 3 (auth primitives) can start once Task 1 config is stable; it does not need the DB.
- Task 2 (migrations) and Task 3 can run in parallel after Task 1.
- Task 4 depends on both; Task 5 depends on Task 4.
- Do not parallelize schema changes with repository code without a shared schema contract.

## Open Questions

1. Should `JWT_SECRET` rotation (two active secrets) be added before production?
2. Add `POST /api/v1/auth/refresh` now or defer until after login works?
3. Retention defaults for `auths`.
4. Where will the Postgres DSN come from for manual verification?
5. Seed credentials: keep defaults in code (`admin@example.com`) or move them to env?
6. When should `tests/` and `.github/` be enabled?
