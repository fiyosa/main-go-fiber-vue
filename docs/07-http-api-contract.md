# 7. HTTP API Contract

## 7.1 Global routing rules

Every JSON API lives under a single prefix: `/api/v1`. There is no separate admin API path. The `/api` group is mounted as `/api/v1` in `internal/routes/api.go`, so module routes declare only their relative path (for example `/auth/login` becomes `/api/v1/auth/login`).

| Group      | Prefix                            | Middleware                            | Notes                                     |
| ---------- | --------------------------------- | ------------------------------------- | ----------------------------------------- |
| API (all)  | `/api/v1`                         | `auth:admin` + bindings               | Every context route file, admin or not    |
| API public | `/api/v1`                         | Bearer token validation                 | External integrations under the same root |
| Public     | `/`                               | none                                  | Media, file upload, PDF, SSO, SPA         |
| SSO        | `/auth/login/oauth2`, `/callback` | none                                  | See [Section 8.9](08-authentication.md#89-sso-callback)                           |
| SPA        | `/*` fallback                     | none                                  | `internal/routes/web.go` serves the FE shell       |

Binding middleware: Laravel uses implicit route-model binding plus `subtitute.bindings`. In Go, controllers explicitly parse IDs and load entities. Return `404` when not found, matching Laravel behavior.

## 7.2 Route inventory — Core (`/api/v1`)

| Module            | Endpoints                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Auth              | `POST /login`, `POST /logout`, `GET /me` (plus new `POST /refresh`)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| Users             | `GET /users`, `GET /users/status-count`, `GET /users/notifications`, `PUT /users/notifications/read-all`, `PUT /users/notifications/{id}`, `GET /users/{user}`, `POST /users`, `PATCH /users/{user}/change-role`, `PATCH /users/{user}/change-division`, `PUT /users/{user}`, `DELETE /users/{user}`, `DELETE /users`                                                                                                                                                                                                                                                                                                             |
| Impersonation     | `POST /impersonation/{user}`, `POST /impersonation/stop`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| Profile           | `GET /profile`, `PUT /profile`, `POST /profile/upload-profile-image`, `POST /profile/change-user-password`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| Menus             | `GET /menus`, `GET /menus/manage`, `GET /menus/manage/parents`, `POST /menus/manage`, `PUT /menus/manage/{menu}`, `DELETE /menus/manage/{menu}`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| Roles             | `GET /roles/all`, `GET /roles`, `GET /roles/permissions`, `POST /roles`, `PUT /roles/{role}`, `DELETE /roles/{role}`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| Settings          | `GET /api/v1/settings`, `POST /api/v1/settings` (web group, absolute path)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| Dashboard         | `GET /api/v1/dashboard/change-log` (web group, absolute path)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| Media             | `DELETE /api/v1/media/{media}`, `GET /api/v1/media/{media}/download`; public `GET /media/{media:file_name}`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| Attendance        | `GET /attendances`, `GET /attendances/me`, `POST /attendances/checkin`, `POST /attendances/checkout`, `POST /attendances/start-aux`, `PUT /attendances/stop-aux`, `GET /attendances/aux-options`, `GET /attendances/user-options`, `GET /attendances/shift-options`, `GET /attendances/date-options`, `GET /attendances/export-template`, `POST /attendances/import`, `GET /attendances/{attendance}/detail`, `POST /attendances/export`, `POST /attendances`, `PUT /attendances/{attendance}/update`, `PUT /attendances/{attendance}/reset`, `DELETE /attendances/{attendance}`, `PUT /attendances/{attendance}/toggle-auto-wip` |
| Appointments      | `GET /appointments`, `GET /appointments/get-status-with-count`, `POST /appointments/create`, `GET /appointments/{appointment}/edit`, `PUT /appointments/{appointment}`, `DELETE /appointments/{appointment}`, `GET /appointments/dashboard/status-count`; `GET /clients/get-clients`                                                                                                                                                                                                                                                                                                                                              |
| Broadcasts        | `apiResource /broadcasts`, `GET /broadcast-template-export`, `GET /broadcast/{broadcast}/send`, `GET /broadcast/messages`, `GET /broadcast-status`, `POST /broadcast-export`                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| Divisions         | `apiResource /divisions`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| Documents         | `GET /documents`, `POST /documents`, `DELETE /documents/{document}`, `GET /documents/{document}/download`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| Email templates   | `GET/POST /manage-email-templates`, `PUT/DELETE /manage-email-templates/{emailTemplate}`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| Forms             | `GET /forms/form-options`, `GET /forms/{id}`; `GET/POST /manage-forms`, `PUT/DELETE /manage-forms/{form}`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| Form requests     | `GET /form-requests`, `GET /form-requests/{id}`, `POST /form-requests`, `POST /form-requests/approve`, `POST /form-requests/reject`, `POST /form-requests/export-history`, `DELETE /form-requests/{id}`; public `GET /form-requests/{id}/generate-pdf`                                                                                                                                                                                                                                                                                                                                                                            |
| Holidays          | `GET/POST /holidays`, `PUT/DELETE /holidays/{holiday}`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| HRDWH             | `GET /hrdwh/baru`, `GET /hrdwh/jabatan`, `GET /hrdwh/keluar`, `GET /hrdwh/keluar/status`, `PATCH /hrdwh/keluar/{nik}/deactivate`, `POST /hrdwh/keluar/export-excel`, `GET /hrdwh/mutasi`, `POST /hrdwh/mutasi/export`, `GET /hrdwh/mutasi/status`, `GET /hrdwh/rpt`, `GET /hrdwh/check-hrdwh-update`, `POST /hrdwh/generate-hrdwh`, `GET /hrdwh/rpt/status`                                                                                                                                                                                                                                                                       |
| Personal workbook | `GET /personal-workbook`, `GET /personal-workbook/{nik}`, `POST /personal-workbook/{nik}/sync-and-encrypt`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| Domino letters    | `GET /domino/letters`, `POST /domino/letters/export`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| Mailadmin         | `GET /mailadmin/profile`, `PUT /mailadmin/profile`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| Shifts            | `PUT /shifts/{id}/toggle-automation`; `GET/POST /manage-shifts`, `PUT/DELETE /manage-shifts/{shift}`, `GET/POST /manage-shifts/patterns`, `PUT/DELETE /manage-shifts/patterns/{shiftPattern}`, `GET /manage-shifts/options`, `GET /manage-shifts/organization-options`; `GET /shift-schedules/shift-today`; `GET /shift-baselines/organization-option`                                                                                                                                                                                                                                                                            |
| Tasks             | `GET/POST /tasks`, `GET/DELETE /tasks/{task}`, `GET/POST /tasks/{task}/list`, `PUT /tasks/{card}/card`, `POST /tasks/{list}/card`, `DELETE /tasks/{card}/card`, `POST /tasks/{card}/upload-file`, `GET /tasks/{card}/download`, `GET /tasks/{task}/board`                                                                                                                                                                                                                                                                                                                                                                         |
| Aux types         | `GET/POST /manage-aux-types`, `PUT/DELETE /manage-aux-types/{auxType}`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| Tgpass            | `POST /api/v1/tgpass/apps`, `POST /api/v1/tgpass/reset` (web group, absolute path)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| Temp sync         | `GET /api/v1/kypas/sync-password/{username}` (throttle 2/min)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |

Auth removals: the second-factor verification endpoint, MFA management endpoints, trusted-device cookie, and LDAP login-check endpoint are not ported.

## 7.3 Route inventory — Other contexts

| Context           | Base path                       | Key endpoints                                                                                                                                                                                                     |
| ----------------- | ------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Kypas             | `/api/v1/kypas`                    | `approvals`, `approval-export`, `approval-status`, `employees`, `contractors`, `register-contractors`, `generate-id`, `inactive-user`, `upload-user`, `nde`, `companies`, `logs`, `search-employee`, `change-pju` |
| Kypas tokens      | `/api/v1/appstokens`               | CRUD, `status`, `send-email-auth-token-expired`                                                                                                                                                                   |
| Kypas public      | `/api/v1/kypas/file-upload`     | `POST` with `verify.header`                                                                                                                                                                                       |
| Kypas public file | `/file-upload/{file:file_name}` | Serves uploaded file                                                                                                                                                                                              |
| Sygap             | `/api/v1/sygap`                    | `employees`, `incidents`, `requests`, `dashboard`, `organizations`, `services`, `service-types`, `tasks`                                                                                                          |
| Intra             | `/api/v1/intra`                    | `applications`, `approval-status`, `applications/{id}/comments`, comments export                                                                                                                                  |
| Kampiun           | `/api/v1/kampiun`                  | `hits-access-documents`, export, `recap/export`                                                                                                                                                                   |
| Rfc               | `/api/v1/rfc`                      | `roles/all`, `users`, `users/{user}`                                                                                                                                                                              |
| Tgcast            | `/api/v1/tgcast`                   | `mail`, `status-count`, `mail-export`                                                                                                                                                                             |
| Tgkypas           | `/api/v1/tgkypas`                  | `dashboard`, `user-access-request`, `sync-*`, `extend-contract`, `change-pju`, `contractors`, `applications`, `users`, `user-contracts`                                                                           |
| Windigital        | `/api/v1/windigital`               | `documents`, `documents/{document}/download`                                                                                                                                                                      |
| Windigital public | `/api/v1/windigital/documents`  | `POST` with client-owner middleware                                                                                                                                                                               |

## 7.4 Route registration files

The `internal/routes/` files own the HTTP wiring. `internal/routes/api.go` mounts the `/api` group as `/api/v1`; `internal/routes/v1/apiv1.go` only re-exports the module route registrations from `internal/app/presentation/core/routes`; `internal/routes/web.go` serves the build assets and the `/*` SPA fallback. `internal/bootstrap/app.go` calls `routes.RegisterAPI` and then `routes.RegisterWeb`.

```go
// internal/routes/api.go
package routes

import (
    "github.com/gofiber/fiber/v3"

    v1 "go-fiber-vue/internal/routes/v1"
)

func RegisterAPI(app *fiber.App, deps v1.Dependencies) {
    api := app.Group("/api")
    v1Group := api.Group("/v1")
    v1.Register(v1Group, deps)
}
```

```go
// internal/routes/v1/apiv1.go
package v1

import (
    "github.com/gofiber/fiber/v3"

    coreroutes "go-fiber-vue/internal/app/presentation/core/routes"
)

// Dependencies carries the controllers and middleware shared by every module route.
type Dependencies struct {
    Core DependenciesCore
    Auth DependenciesAuth
}

func Register(group fiber.Router, deps Dependencies) {
    coreroutes.RegisterAuthRoutes(group, deps.Core.AuthController, deps.Auth.Admin, deps.Auth.Guest)
    coreroutes.RegisterRoleRoutes(group, deps.Core.RoleController, deps.Auth.Admin, deps.Auth.Permission("manage_roles"))
    coreroutes.RegisterAttendanceRoutes(group, deps.Core.AttendanceController, deps.Auth.Admin)
    // ... one call per module route file exported by internal/app/presentation/core/routes
}
```

```go
// internal/routes/web.go
package routes

import (
    "path/filepath"
    "strings"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/static"

    "go-fiber-vue/internal/app/presentation/core/view"
)

func RegisterWeb(app *fiber.App, publicDir string) {
    buildDir := filepath.Join(publicDir, "build")

    // /build/* holds the Vite output: index.js, index.css and assets/*.
    app.Use("/build", static.New(buildDir, static.Config{MaxAge: 3600, Compress: true}))

    // Every other non-API path renders the embedded SPA shell.
    app.Get("/*", func(c fiber.Ctx) error {
        switch {
        case strings.HasPrefix(c.Path(), "/api"):
            return c.Next()
        case strings.HasPrefix(c.Path(), "/build"):
            return fiber.ErrNotFound
        }
        c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
        c.Set(fiber.HeaderCacheControl, "no-cache")
        return c.SendString(view.SPA)
    })
}
```

`internal/routes/web.go` serves the build output from `public/build` (`index.js`, `index.css`, `assets/*`) and installs the `/*` SPA fallback. The shell itself is `internal/app/presentation/core/view/spa.html`, embedded with `//go:embed` as `view.SPA`, so the Vite build does not need to emit an `index.html`. `internal/bootstrap/app.go` calls `routes.RegisterWeb` last, after `routes.RegisterAPI`, so the fallback never shadows API routes. Unknown `/api/*` paths and missing `/build/*` files still return the JSON 404.

`Dependencies` is declared in `internal/routes/v1/apiv1.go` and consumed by `internal/routes/api.go`. Keeping it in `v1` preserves a one-way dependency chain (`internal/bootstrap` → `routes` → `routes/v1` → `internal/app/presentation`) with no import cycle and no extra file.

Rules:

- One module route file per Laravel route file, same module name, exported from `internal/app/presentation/core/routes`.
- Every API path is under `/api/v1`; module route files declare only relative segments.
- `internal/routes/v1/apiv1.go` contains only the aggregation of the presentation route files.
- `internal/routes/web.go` serves `/build/*` and owns the `/*` SPA fallback; `internal/bootstrap/app.go` calls it last so API routes are matched first.
- Register explicit middleware per route, matching the Laravel route file.
- Do not rely on automatic file scanning.

## 7.5 Contents of `internal/routes/`

`internal/routes/` is the single entry point for wiring the HTTP surface. It contains these files:

```text
internal/routes/
├── api.go       # owns /api, mounts /api/v1, delegates to v1.Register
├── web.go       # /build/* static assets + /* SPA fallback
└── v1/
    └── apiv1.go # Dependencies struct + aggregation of module route files
```

| File | Package | Responsibility | Called by |
| --- | --- | --- | --- |
| `internal/routes/api.go` | `routes` | Declare the `/api` group, nest `/v1`, and delegate to `v1.Register` | `internal/bootstrap/app.go` |
| `internal/routes/web.go` | `routes` | Serve `/build/*` (`index.js`, `index.css`, `assets/*`) and install the `/*` SPA fallback that renders the embedded `view.SPA` shell | `internal/bootstrap/app.go` |
| `internal/routes/v1/apiv1.go` | `v1` | Declare `Dependencies` and call every `Register*Routes` exported by `internal/app/presentation/*/routes` | `internal/routes/api.go` |

Request flow:

```text
/api/v1/...
  └─ internal/routes/api.go               (group /api → /v1)
       └─ internal/routes/v1/apiv1.go     (call module registrations)
            └─ internal/app/presentation/<ctx>/routes/<module>_routes.go
                 └─ internal/app/presentation/<ctx>/controllers/*

/build/*  → internal/routes/web.go → public/build/* assets
anything else (not /api and not /build)
  → internal/routes/web.go (SPA fallback) → embedded view.SPA (spa.html)
```

Rules specific to `internal/routes/`:

- `internal/routes/` lives under `internal/`, not under `app/`, and is the only place that mounts the API and the SPA.
- `api.go` and `web.go` are called by `internal/bootstrap/app.go`; bootstrap does not define its own route files.
- `v1/apiv1.go` holds no endpoint logic; it only aggregates the route files from `internal/app/presentation`.
- Module route files stay in `internal/app/presentation/<ctx>/routes/`, named `<module>_routes.go`, and declare only relative path segments.

## 7.6 Route permission audit

Every route registered by `internal/app/presentation/{domain}/routes/*.go` must have an explicit access category: public, authenticated-only, permission, or impersonation. The complete matrix, Laravel source, backlog, and verification procedure live in [Appendix D](appendices/d-route-permissions.md).

Rules:

- `auth:admin` is inherited from the API group; module route files must not repeat it.
- A route that needs a permission must attach `deps.Auth.Permission("<name>")` to the route or its group.
- A controller or request validator must not be the only place a permission is enforced.
- A route with no Laravel permission stays authenticated-only until an intentional difference is recorded in [Section 18](18-risks-and-open-decisions.md).
- A new route without an entry in Appendix D fails the route permission test.

## 7.7 Response messages

Every JSON API response must include a non-empty `message`. The message is translated through the localization layer described in [Section 20](20-localization.md).

- Success and error response helpers receive a translation key, not a literal.
- The default locale is `en`; `id` is supported per request.
- `Accept-Language` selects the locale; `?lang=` may override it for testing.
- Missing keys return the key and log a warning.
- JSON delete endpoints return `200` with a message instead of `204 No Content`.
- Contract tests must assert `message` for `2xx`, `401`, `403`, `404`, `422`, and `500`.
