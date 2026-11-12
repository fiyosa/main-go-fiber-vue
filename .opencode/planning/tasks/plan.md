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
- [ ] Task 1: Add `AuthController.Me` and mount `GET /api/v1/auth/me`
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
- [ ] Task 2: Add `getMe()` API client
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
- [ ] Task 3: Add `authGuard` and wire it into the router
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
- [ ] Task 4: Run the full verification suite
  - **Acceptance criteria:**
    - All Go tests pass.
    - All Vitest tests pass.
    - Frontend production build succeeds.
    - Go build succeeds.
  - **Verification:** `go test ./...`, `npm --prefix web run test`, `npm --prefix web run build`, `go build -o ./tmp/app.exe .`
  - **Dependencies:** Tasks 1-3
  - **Files:** None
  - **Scope:** XS

## Checkpoint: Complete
- [ ] All acceptance criteria met.
- [ ] Ready for `/review`.
