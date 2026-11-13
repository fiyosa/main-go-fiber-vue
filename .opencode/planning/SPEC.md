# Spec: Auth State & Redirect Guards

## Objective
Prevent authenticated users from visiting `/login` or `/register` by redirecting them to `/admin`, and protect `/admin` by redirecting unauthenticated users to `/login`.

## Scope

### Backend
1. Add `GET /api/v1/auth/me`:
   - Validate the HTTP-only access cookie using the existing `middleware.Authenticate` handler.
   - Return `200` with `{ success: true, data: { id, email }, message }` when the token is valid.
   - Return `401` when the cookie is missing or invalid, `503` when the token service is unavailable.
   - Register the route as public in the route audit (`publicAPIRoutes`) because it requires authentication but no permission.
2. Add `AuthController.Me` which reads the authenticated identity from request locals set by `middleware.Authenticate`.
3. Wire the route through `middleware.Authenticate(deps.Tokens)` in `RegisterAuthRoutes`.

### Frontend
1. Add `getMe()` API client under `web/src/api/auth/me.ts` and export it from `web/src/api/auth/index.ts`.
2. Add `web/src/pages/authGuard.ts` with a Vue Router navigation guard:
   - Guest-only paths: `/login` and `/register`. If `getMe()` succeeds, redirect to `/admin`; if it fails, allow the navigation.
   - Admin paths: if `getMe()` succeeds, allow; if it fails, redirect to `/login`.
   - Public pages: no auth check.
3. Attach the guard with `router.beforeEach(authGuard)` in `web/src/pages/routes.ts`.
4. Mark the `/admin` route with `meta: { requiresAuth: true }`.
5. Do not add a register page. `/register` remains an unmatched guest path and renders the existing 404 page for logged-out visitors.

### Testing
1. Backend: add a controller/route test proving `/api/v1/auth/me` returns the identity for a valid token and 401 for a missing/invalid token.
2. Frontend: add unit tests for `authGuard` covering login, register, admin, and public paths.
3. Update route tests to assert the admin route requires auth and `/register` still resolves to the guest 404 route.

## Tech Stack
- Go + Fiber v3
- Vue 3 + Vue Router
- Vitest
- Go testing

## Commands
- Test Go: `go test ./...`
- Vet Go: `go vet ./...`
- Build Go: `go build -o ./tmp/app.exe .`
- Test Frontend: `npm --prefix web run test`
- Build Frontend: `npm --prefix web run build`

## Success Criteria
- Visiting `/login` or `/register` while authenticated redirects to `/admin`.
- Visiting `/admin` while unauthenticated redirects to `/login`.
- Logged-out `/register` still shows the guest 404 page.
- All unit tests pass and both builds succeed.

## Follow-up: Register API & OpenAPI Refinement

### Scope
1. Add `POST /api/v1/auth/register`:
   - No session or cookies are issued.
   - Assigns role `user` when it exists.
   - Password validation enforces `min=8`.
   - Duplicate email returns `409` with `auth.email_taken`.
2. Tidy OpenAPI labels and group tags:
   - Summaries omit the HTTP method prefix; Scalar renders the method badge separately.
   - Groups: `Auth` (login/register/logout/me), `Roles` (roles), `Health` (`healthz`/`readyz`); fallback `Default`.
   - Regenerate `public/docs/openapi.json` with `go run ./cmd/console openapi`.
3. Keep `lang/en.json` and `lang/id.json` key parity.

### Success Criteria
- Register returns `201` with `{id, email}` and never issues cookies.
- Every OpenAPI operation has a clean label and correct group.
- Locale catalogs stay in sync and all Go tests, vet, and build pass.
