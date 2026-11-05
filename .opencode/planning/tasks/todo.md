# Todo: Slice 1 — Minimal Login API

Scope: one public endpoint, `POST /api/v1/auth/login`, using only `users` and `auths`.
No RBAC, no roles/permissions/menus, no LDAP, no 2FA, no SSO.
Do **not** create `./tests` or `./.github` yet.

Localization contract: `docs/20-localization.md`. Every JSON response carries a translated
`message`; catalogues live in `internal/lang/en.json` and `internal/lang/id.json`.

## Done

- [x] Task 0 — Module and folder scaffold
  - `go.mod` = `module go-fiber-vue`, Go 1.27.3
  - Folder spine created; `tests/` and `.github/` absent

- [x] Task 1 — Minimal bootable runtime, config, and DB pool
  - [x] `cmd/app/main.go` runs Fiber under signal-aware runtime; `cmd/console/main.go` compiles
  - [x] `internal/config/config.go` reads `APP_ENV`, `APP_PORT`, `APP_LOCALE`, `DB_*`, `JWT_*`
  - [x] `APP_LOCALE` validated against `en`/`id`; invalid value fails startup
  - [x] `.env` loaded at startup (real env wins); `.env.example` documents every variable
  - [x] `internal/platform/database/database.go` opens a Postgres pool (no AutoMigrate)
  - [x] `GET /healthz` → 200
  - [x] `GET /readyz` → 200 with DB, 503 JSON without
  - [x] `recover` middleware so a handler panic cannot kill the process
  - [x] Ctrl+C shuts down cleanly
  - Verify: `go build ./...`, `go vet ./...` pass; manual run confirmed 200/503

- [x] Task 2 — Schema and seed user
  - [x] `users` migration with `id`, `email`, `password`, and timestamps
  - [x] `auths` migration with unique `token`, unique `refresh_token`, expiry/revoked indexes
  - [x] `console migrate`, `migrate status`, `migrate down` implemented
  - [x] `console db seed` creates one admin (bcrypt cost 12), idempotent
  - [x] Seed refuses empty/short passwords (min 12 chars)
  - [ ] Apply up/down against a real Postgres — blocked: no DSN

- [x] Task 3 — Password, JWT, and refresh-token services
  - [x] bcrypt hash/verify (cost 12)
  - [x] HS256 JWT issue/verify with `iss`, `aud`, `sub`, `email`, `jti`, `iat`, `exp`
  - [x] Reject wrong `aud`, expired tokens, and invalid signatures
  - [x] Refresh token: 256-bit random, base64url raw + SHA-256 hex hash
  - Verify: `go test ./internal/platform/auth` passes

- [x] Task 4 — Login use case with repositories
  - [x] `User` and `Auth` entities
  - [x] `UserRepository` + `AuthRepository` interfaces and GORM implementations
  - [x] `LoginUseCase` verifies password, issues tokens, and inserts `auths`
  - [x] Unknown user / wrong password → `ErrInvalidCredentials`
  - [x] No permission/scope claims; no LDAP/2FA branch
  - Verify: table-driven use-case tests with fakes pass

- [x] Task 5 — `POST /api/v1/auth/login` end-to-end
  - [x] `LoginRequest`, controller, `LoginResponse`, public route
  - [x] Valid login → 200 with `success`, `message`, `token`, `token_type`, `expires_in`, `refresh_token`
  - [x] Invalid login → 401 with `lang.invalid_credentials`
  - [x] Missing fields → 422 with translated field errors
  - [x] Route has no permission middleware
  - [x] Writes one `auths` row (covered by use-case test)
  - [ ] Live login against Postgres — blocked: no DSN

- [x] Task 6 — Localization (`docs/20-localization.md`)
  - [x] `internal/lang/en.json` + `internal/lang/id.json` copied from `docs/lang/`, key parity enforced at load
  - [x] `internal/platform/lang` translator: nested catalogues, dot keys, placeholders, fallback, missing-key warning
  - [x] Locale middleware: `?lang=` → `Accept-Language` → `APP_LOCALE`
  - [x] Response helpers accept translation keys, never literals; `message` never omitted
  - [x] Validation messages translated with `:attribute` / `:min` / `:max`
  - [x] JSON error handler localizes 400/401/403/404/405/415/500
  - [x] Tests: translator, catalogues, middleware, controller messages, config validation

### Checkpoint A/B/C

- [x] Build/vet/tests pass
- [x] Login response shape and localized messages approved by tests
- [ ] Human review of schema and token claims
- [ ] End-to-end login against Postgres

## Deferred (do not start yet)

- [ ] `POST /api/v1/refresh` — wait until login is stable.
- [ ] `GET /api/v1/me` — out of scope for this slice.
- [ ] Roles/permissions/menus — not part of this slice.
- [ ] `tests/` and `.github/` — create only after approval.
- [ ] Scheduler/jobs — not needed for login.
- [ ] `console db seed` integration test — needs a DSN.

## Standing Rules

- Run `gofmt`, `go vet ./...`, and `go test ./...` before handoff.
- Never commit secrets or log tokens/passwords.
- Every response message comes from `lang/*.json`; both files change together.
- Do not add SSO, PDF, LDAP login, 2FA, or MCP tools.
- Do not create `tests/` or `.github/` until approved.
- No permission strings or RBAC tables in this slice.
