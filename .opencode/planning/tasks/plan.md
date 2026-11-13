# Implementation Plan: Auth State & Redirect Guards

## Overview
Add `GET /api/v1/auth/me` so the SPA can probe the HTTP-only session, then guard `/login`, `/register`, and `/admin` with a Vue Router navigation guard.

## Architecture Decisions
- `/api/v1/auth/me` uses the existing `middleware.Authenticate` handler. It requires a valid access token but no permission.
- The route audit lists it as public (`publicAPIRoutes`) because it is an authentication probe, matching how `POST /api/v1/auth/login` is treated.
- `AuthController.Me` reads the identity from request locals set by `middleware.Authenticate`.
- The frontend guard lives in `web/src/pages/authGuard.ts` so it can be unit tested without mounting the whole router.
- `/register` stays unmatched for logged-out visitors, so the existing guest 404 route still renders.

## Task List

### Phase 1: Backend identity endpoint
- [x] Task 1: Add `AuthController.Me` and mount `GET /api/v1/auth/me`
  - **Acceptance criteria:**
    - Valid access cookie returns `200` with `{ success: true, data: { id, email }, message }`.
    - Missing or invalid token returns `401`.
    - Nil token service returns `503`.
    - Route audit test passes with the route listed as public.
  - **Verification:** `go test ./...` and `go vet ./...`
  - **Dependencies:** None
  - **Files:** `internal/app/presentation/core/controllers/auth_controller.go`, `internal/app/presentation/core/routes/auth_routes.go`, `internal/bootstrap/routes_test.go`, `internal/app/presentation/core/controllers/auth_controller_test.go`
  - **Scope:** Small

### Phase 2: Frontend API client
- [x] Task 2: Add `getMe()` API client
  - **Acceptance criteria:**
    - Calls `GET /api/v1/auth/me`.
    - Resolves with `{ id, email }` on success.
    - Rejects with the normalized axios error on failure.
    - Exported from `web/src/api/auth/index.ts`.
  - **Verification:** `npm --prefix web run test`
  - **Dependencies:** Task 1 contract
  - **Files:** `web/src/api/auth/me.ts`, `web/src/api/auth/index.ts`
  - **Scope:** XS

### Phase 3: Router guard
- [x] Task 3: Add `authGuard` and wire it into the router
  - **Acceptance criteria:**
    - `/admin` route has `meta: { requiresAuth: true }`.
    - Authenticated visits to `/login` or `/register` redirect to `/admin`.
    - Unauthenticated visits to `/admin` redirect to `/login`.
    - Public pages do not call `getMe()`.
    - Logged-out `/register` still resolves to the guest 404 route.
  - **Verification:** `npm --prefix web run test`
  - **Dependencies:** Task 2
  - **Files:** `web/src/pages/authGuard.ts`, `web/src/pages/routes.ts`, `web/src/pages/routes.test.ts`
  - **Scope:** Small

### Phase 4: Verification
- [x] Task 4: Run the full verification suite
  - **Acceptance criteria:**
    - All Go tests pass.
    - All Vitest tests pass.
    - Frontend production build succeeds.
    - Go build succeeds.
  - **Verification:** `go test ./...`, `npm --prefix web run test`, `npm --prefix web run build`, `go build -o ./tmp/app.exe .`
  - **Dependencies:** Tasks 1-3
  - **Files:** None
  - **Scope:** XS

## Follow-up: Register API & OpenAPI Refinement

### Phase 5: Register API
- [x] Task 5: Add `POST /api/v1/auth/register`
  - **Acceptance criteria:**
    - No session or cookies are issued.
    - Assigns role `user` when it exists.
    - Password validation enforces `min=8`.
    - Duplicate email returns `409`.
    - Route is public in the route audit.
  - **Verification:** `go test ./...`, `go vet ./...`
  - **Files:** `internal/app/application/core/repository/user/repository.go`, `internal/app/application/core/repository/auth/register_repository.go`, `internal/app/presentation/core/request/auth_register_request.go`, `internal/app/presentation/core/response/auth_register_response.go`, `internal/app/presentation/core/controllers/auth_controller.go`, `internal/app/presentation/core/routes/auth_routes.go`, `internal/bootstrap/wire_api.go`

### Phase 6: OpenAPI labels and groups
- [x] Task 6: Tidy OpenAPI labels and group tags
  - **Acceptance criteria:**
    - Every operation has a clean label; summaries omit the HTTP method prefix because Scalar renders the method badge separately.
    - Groups: `Auth` (login/register/logout/me), `Roles` (roles), `Health` (healthz/readyz); fallback stays `Default`.
    - `public/docs/openapi.json` regenerated with `go run ./cmd/console openapi`.
  - **Verification:** `go test ./internal/platform/docs/...`, `go test ./...`
  - **Files:** `internal/platform/docs/spec.go`, `internal/platform/docs/spec_test.go`, `public/docs/openapi.json`

### Phase 7: Locale parity
- [x] Task 7: Keep `en.json` and `id.json` in sync
  - **Acceptance criteria:**
    - Same key set in both catalogs.
    - `email_already_taken` present; no stray trailing commas.
  - **Verification:** `go test ./...` (locale load test)
  - **Files:** `lang/en.json`, `lang/id.json`

## Checkpoint: Complete
- [x] All acceptance criteria met.
- [x] Ready for `/review`.
