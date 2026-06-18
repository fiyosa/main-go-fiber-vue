# PRD — Go Fiber Portfolio Project

## 1. Project Overview

- **Name:** go-fiber-svelte
- **Module:** `go-fiber-svelte`
- **Stack:** Go 1.25+ (Fiber v2) + PostgreSQL (GORM) + Svelte 5 SPA (Vite + Tailwind CSS v4)
- **Project Layout:** [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- **Server:** Go Fiber (via `gofiber/fiber/v2`)
- **Deployment:** Docker (multi-stage: golang:alpine build + alpine runtime, port 8000)
- **Purpose:** Backend API with RBAC auth system, serving as a portfolio backend + SPA frontend

---

## 2. Project Structure

Project follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout) conventions:

| Directory   | Purpose                                                                                                                                           |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd/`      | Main applications — one subdirectory per binary (`cmd/app`, `cmd/migrate`, `cmd/seed`). Each has a small `main.go` that imports from `internal/`. |
| `internal/` | Private application code — **not importable by external projects** (enforced by Go compiler). All core logic lives here.                          |
| `web/`      | Web application components — Vite + Svelte 5 SPA frontend.                                                                                        |
| `public/`   | Build output and static assets served by Go Fiber.                                                                                                |
| `build/`    | Packaging — Dockerfile in `build/package/`.                                                                                                       |
| `logs/`     | Runtime log files — created automatically on first log write, gitignored.                                                                         |
| —           | `.env.example` at project root.                                                                                                                   |
| `docs/`     | Project documentation.                                                                                                                            |
| `scripts/`  | Build/deploy scripts.                                                                                                                             |

> **Key rule:** No `src/` directory at the project root. Go's workspace (`$GOPATH`) has its own `src/`, but project-level `src/` is a Java anti-pattern.

```
.
.
├── go.mod                       # Module: go-fiber-svelte
├── go.sum
├── .gitignore
├── .air.toml                    # Hot-reload with go build (Windows .exe)
│
├── cmd/                         # Main applications
│   ├── app/
│   │   └── main.go              # Entry point: init config, DB, Fiber app, serve SPA
│   └── migrate/
│       └── main.go              # DB migration runner (CLI)
│
├── internal/                    # Private application code (enforced by Go compiler)
│   ├── bootstrap/
│   │   └── app.go               # Init middleware, register routes
│   │
│   ├── config/
│   │   ├── app.go               # App-level config (APP_LOG, APP_URL, etc.)
│   │   ├── config.go            # Load .env → Config struct
│   │   └── cors.go              # CORS configuration
│   │
│   ├── http/                    # HTTP layer
│   │   ├── controllers/         # Thin HTTP handlers
│   │   │   ├── auth_controller.go
│   │   │   ├── guest_controller.go
│   │   │   ├── doc_controller.go
│   │   │   ├── logger_controller.go  # Log viewer endpoints
│   │   │   └── policy_controller.go
│   │   ├── repositories/        # Business logic (1 file per endpoint)
│   │   │   ├── auth/
│   │   │   │   ├── login_repository.go
│   │   │   │   ├── logout_repository.go
│   │   │   │   └── user_repository.go
│   │   │   ├── logger_repository/  # Log viewer business logic
│   │   │   │   ├── log_list_repository.go
│   │   │   │   ├── log_detail_repository.go
│   │   │   │   ├── log_download_repository.go
│   │   │   │   └── log_delete_repository.go
│   │   │   └── policy/
│   │   ├── request/             # Struct validasi (go-playground/validator)
│   │   │   ├── auth/
│   │   │   │   └── login_request.go
│   │   │   └── policy/
│   │   │       └── permission_store_request.go
│   │   └── resources/           # Data transformers (encodeId tiap field id)
│   │       ├── auth/
│   │       │   └── user_resource.go
│   │       └── policy/
│   │           ├── role_list_resource.go
│   │           └── permission_resource.go
│   │
│   ├── db/                      # Database layer (GORM)
│   │   ├── db.go                # Koneksi & init GORM (PostgreSQL)
│   │   ├── models/              # GORM model structs (7 tabel)
│   │   │   ├── user.go
│   │   │   ├── user_detail.go
│   │   │   ├── auth.go
│   │   │   ├── role.go
│   │   │   ├── permission.go
│   │   │   ├── role_has_permission.go
│   │   │   └── user_has_role.go
│   │   ├── migrations/          # AutoMigrate / SQL migration files
│   │   └── seed/
│   │       └── seed.go          # Seeder data awal
│   │
│   ├── helper/
│   │   └── response.go          # Res.Success, Res.SuccessData, Res.Error, Res.Paginate, Res.Catch, Res.Validate
│   │
│   ├── lang/                    # i18n
│   │   ├── lang.go              # t._("key"), t._("key", args) helper
│   │   └── locales/
│   │       └── en.json          # English translations
│   │
│   ├── lib/                     # Reusable library packages
│   │   ├── hash.go              # bcrypt generate/verify + hashid encodeId/decodeId
│   │   ├── jwt.go               # JWT create/verify
│   │   ├── logger.go            # Logger (zerolog)
│   │   └── validator.go         # go-playground/validator wrapper
│   │
│   ├── middleware/
│   │   ├── auth_middleware.go    # Validasi JWT cookie → c.Locals("user_id")
│   │   └── hash_middleware.go   # Decode Hashids params ke numeric ID
│   │
│   ├── openapi/
│   │   └── openapi.go           # Generate OpenAPI spec
│   │
│   ├── provider/                # Core engine
│   │   ├── error_provider.go    # Global error handler (ValidationError, fiber.Error)
│   │   ├── auth_provider.go     # RBAC: user_has_roles → role_has_permissions
│   │   └── app_provider.go      # Middleware registration
│   │
│   ├── routes/
│   │   └── api.go               # Route definitions & handler mapping
│   │
│   └── utils/
│       ├── date.go              # now(), formatByDate(), formatByStr()
│       └── uuid.go              # UUID v4 generate & validate
│
├── web/                         # Vite + Svelte 5 SPA
│   ├── index.html               # Entry HTML
│   ├── package.json
│   ├── vite.config.ts            # Vite + Tailwind v4 plugin + Svelte
│   ├── svelte.config.js
│   ├── tsconfig.json
│   ├── tsconfig.node.json
│   ├── .gitignore
│   ├── .prettierrc
│   ├── .vscode/
│   │   └── settings.json        # Prettier format-on-save for Svelte
│   │
│   ├── src/
│   │   ├── pages/
│   │   │   ├── routes.ts        # Route definitions (map with '*' catch-all)
│   │   │   └── guest/
│   │   │       ├── Home.svelte
│   │   │       ├── About.svelte
│   │   │       ├── Logger.svelte  # Log viewer UI (sidebar + entries + pagination)
│   │   │       └── NotFound.svelte
│   │   │
│   │   ├── lib/                 # Svelte components & logic
│   │   │   ├── components/
│   │   │   │   └── ...svelte
│   │   │   ├── stores/
│   │   │   │   └── auth.ts
│   │   │   ├── tanstackUtil.ts  # useQuery/useMutation wrapper for TanStack Query v6
│   │   │   └── utils/
│   │   │       ├── axiosLib.ts  # Axios instance + createQueryStr utility
│   │   │       └── ...
│   │   │
│   │   ├── api/                 # Frontend API hooks (TanStack Query)
│   │   │   ├── logger/          # Log viewer API hooks
│   │   │   │   ├── index.ts
│   │   │   │   ├── getLogs.ts
│   │   │   │   ├── getLogDetail.ts
│   │   │   │   └── deleteLog.ts
│   │   │   ├── auth/
│   │   │   │   ├── index.ts
│   │   │   │   ├── login.ts
│   │   │   │   ├── logout.ts
│   │   │   │   └── user.ts
│   │   │   ├── guest/
│   │   │   │   ├── index.ts
│   │   │   │   └── ping.ts
│   │   │   └── policy/
│   │   │       ├── index.ts
│   │   │       ├── role_list.ts
│   │   │       ├── permission_list.ts
│   │   │       ├── permission_store.ts
│   │   │       └── permission_destroy.ts
│   │   │
│   │   ├── App.svelte           # Root: Router + routes
│   │   ├── main.ts              # Entry: setHashRoutingEnabled(false) + mount
│   │   ├── app.css              # Global styles @import "tailwindcss"
│   │   └── vite-env.d.ts
│   │
│   └── node_modules/            # pnpm dependencies (gitignored)
│
├── public/                      # Build output & static assets
│   ├── build/                   # Vite build output (served by Go Fiber)
│   │   ├── index.html
│   │   ├── main.js
│   │   └── main.css
│   ├── favicon.svg
│   ├── openapi.html             # Scalar API docs UI (local)
│   └── scalar-standalone.js     # Vendored Scalar JS (~3.6MB)
│
├── build/                       # Packaging & CI
│   └── package/
│       └── Dockerfile           # Multi-stage: build Go + build SPA
│
├── .env.example                 # Environment variable template
├── docs/
│   └── PRD.md                   # ← This file
│
└── scripts/
    └── deploy.sh                # Deployment script
```

---

## 3. Architecture & Patterns

### 3.1 Request Flow

```
HTTP Request
  → Fiber router (internal/routes/api.go)
    → Middleware chain (auth, hash)
      → Policy check (RBAC via authProvider)
        → HTTP Controller (internal/http/controllers/)
          → HTTP Repository (internal/http/repositories/)
            → Resource (internal/http/resources/)
              → Helper response (JSON)
```

### 3.2 Middleware

Middleware is registered per-route in `routes/api.go`, not globally. Public routes skip auth; protected routes use `AuthMiddleware` via route grouping.

**Existing middleware:**
| Name | File | Purpose |
|------|------|---------|
| `auth` | `auth_middleware.go` | Validates JWT cookie, sets `c.Locals("user_id")` |
| `hash` | `hash_middleware.go` | Decodes Hashids params to numeric IDs |

### 3.3 Policy/RBAC

Policies are checked via `authProvider.go`. Uses `user_has_roles → role_has_permissions → permissions` chain.

Applied per-route. Returns 403 if user lacks permission.

### 3.4 Controller → Repository Pattern

```
Controller (thin, delegates)
  → Repository (business logic, db queries)
    → Resource (data transformation)
```

Controllers are exported functions that delegate to repositories. Repositories are exported async functions.

### 3.5 Response Helpers

All from `internal/helper/response.go`:
| Function | Purpose |
|----------|---------|
| `Res.Success(msg)` | Success message |
| `Res.SuccessData(data, msg)` | Success with data |
| `Res.Error(msg, errors)` | Error response |
| `Res.Paginate(data, meta, msg)` | Paginated response |
| `Res.Catch(err)` | Catch-block handler (generic error) |
| `Res.Validate(err)` | Validation error |

Semua response API menggunakan `helper.Res.*` — lihat penggunaan di `controllers/`, `repositories/`, dan `provider/error_provider.go`.

**Standard response format:**

```json
{
  "message": "...",
  "data": { ... },
  "errors": { ... },
  "meta": { "total": 0, "page": 1, "limit": 10 }
}
```

---

## 4. Database Schema

All tables use `bigserial` PKs and GORM with relations.

### Tables:

| Table                  | Columns                                                            | Relations                              |
| ---------------------- | ------------------------------------------------------------------ | -------------------------------------- |
| `users`                | id, email, username, password, created_at, updated_at, deleted_at  | → auths, user_details, user_has_roles  |
| `user_details`         | id, user_id, first_name, last_name, created_at, updated_at         | → users                                |
| `auths`                | id, user_id, token, revoke, ip, user_agent, created_at, updated_at | → users                                |
| `roles`                | id, name, notes, created_at, updated_at, deleted_at                | → user_has_roles, role_has_permissions |
| `permissions`          | id, name, notes, created_at, updated_at, deleted_at                | → role_has_permissions                 |
| `role_has_permissions` | role_id, permission_id (composite PK)                              | → roles, permissions                   |
| `user_has_roles`       | user_id, role_id (composite PK)                                    | → users, roles                         |

### Key relationships:

- User M:N Role via `user_has_roles`
- Role M:N Permission via `role_has_permissions`
- All junction tables use composite primary keys
- `deleted_at` used for soft-delete on users, roles, permissions

### GORM client:

```go
// internal/db/db.go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // local only
})
```

---

## 5. API Endpoints

| Method | Path                         | Auth | Policy | Controller                           | Description                    |
| ------ | ---------------------------- | ---- | ------ | ------------------------------------ | ------------------------------ |
| GET    | `/api/openapi.json`          | -    | -      | `docController.openapi`              | OpenAPI 3.0 spec (dynamic)     |
| GET    | `/api/docs`                  | -    | -      | `docController.docs`                 | Scalar API docs UI             |
| POST   | `/api/auth/login`            | -    | -      | `authController.login`               | Login, HTTP-only cookie        |
| DELETE | `/api/auth/logout`           | auth | -      | `authController.logout`              | Revoke token + clear cookie    |
| GET    | `/api/auth/user`             | auth | -      | `authController.user`                | Current user info              |
| GET    | `/api/policy/role`           | auth | -      | `policyController.roleList`          | List roles with permissions    |
| GET    | `/api/policy/permission`     | auth | -      | `policyController.permissionList`    | List permissions               |
| POST   | `/api/policy/permission`     | auth | -      | `policyController.permissionStore`   | Create permission              |
| DELETE | `/api/policy/permission/:id` | auth | -      | `policyController.permissionDestroy` | Delete permission (soft)       |
| GET    | `/api/guest/ping`            | -    | -      | `guestController.ping`               | Health check                   |
| GET    | `/api/log`                   | -    | -      | `loggerController.LoggerList`        | List log files with name + size |
| GET    | `/api/log/:filename`         | -    | -      | `loggerController.LoggerDetail`      | Log entries (filter, search, paginate) |
| DELETE | `/api/log/:filename`         | -    | -      | `loggerController.LoggerDelete`      | Delete log file                 |
| GET    | `/api/log/:filename/download`| -    | -      | `loggerController.LoggerDownload`    | Download log file               |
| \*     | `/api/*`                     | -    | -      | 404 JSON                             | Fallback for invalid API paths |

---

## 6. Libraries & Utilities

### Go Libraries (`internal/lib/`):

| File           | Exports                                      | Purpose                           |
| -------------- | -------------------------------------------- | --------------------------------- |
| `hash.go`      | `Generate`, `Verify`, `EncodeId`, `DecodeId` | bcrypt hashing + Hashids encoding |
| `jwt.go`       | `Create`, `Verify`                           | JWT sign/verify with APP_SECRET   |
| `logger.go`    | `Info`, `Error`, `Info`, `CloseFile`         | Zerolog logging + file handle closing |
| `validator.go` | `Validate`, `ValidateRequest`                | go-playground/validator wrapper   |

### Go Utilities (`internal/utils/`):

| File      | Exports                              | Purpose                         |
| --------- | ------------------------------------ | ------------------------------- |
| `date.go` | `Now`, `FormatByDate`, `FormatByStr` | Date formatting                 |
| `uuid.go` | `Create`, `Verify`                   | UUID v4 generation & validation |

### Frontend Dependencies (SPA):

| Package                        | Version | Purpose                                   |
| ------------------------------ | ------- | ----------------------------------------- |
| `svelte`                       | ^5.55.5 | UI framework (runes-based)                |
| `@keenmate/svelte-spa-router`  | ^5.1.0  | SPA router (history mode with `use:link`) |
| `@tailwindcss/vite`            | ^4.3.0  | Tailwind CSS v4 Vite plugin               |
| `tailwindcss`                  | ^4.3.0  | Utility-first CSS                         |
| `vite`                         | ^8.0.12 | Build tool                                |
| `@sveltejs/vite-plugin-svelte` | ^7.1.2  | Svelte Vite integration                   |
| `prettier`                     | ^3.8.3  | Code formatter                            |
| `prettier-plugin-svelte`       | ^4.1.0  | Svelte Prettier plugin                    |
| `prettier-plugin-tailwindcss`  | ^0.8.0  | Tailwind class sorting                    |
| `@tanstack/svelte-query`      | ^6.3.0  | TanStack Query (async state, cache, mutations) |
| `axios`                        | ^1.7.9  | HTTP client for API requests              |

---

## 7. Configuration & Environment

### Environment Variables (.env):

| Variable           | Example                     | Scope  |
| ------------------ | --------------------------- | ------ |
| `APP_ENV`          | `local`                     | Public |
| `APP_LOCALE`       | `en`                        | Public |
| `APP_SECRET`       | `secret`                    | Secret |
| `APP_JWT_DURATION` | `1d`                        | Public |
| `APP_URL`          | `http://localhost:3000`     | Public |
| `DB_URL`           | `postgresql://...`          | Secret |
| `API_URL`          | `http://localhost:8000`     | Public |
| `PORT`             | `8000`                      | Public |
| `APP_LOG`          | `true`                      | Public |

Config loaded via `internal/config/config.go` (using `github.com/joho/godotenv` + `os.Getenv`). Jika file `.env` tidak ditemukan (misal di container), fallback ke env vars dari Docker `--env-file`.

**`APP_URL`** digunakan oleh `internal/config/cors.go` sebagai `AllowOrigins` di CORS config. Dipisahkan dari `API_URL` karena frontend dev server (port 3000) bisa beda dengan backend (port 8000).

### Logger (`internal/lib/logger.go`)

- Logging dikontrol oleh `APP_LOG` (default `true`)
- `APP_LOG=false` → `zerolog.Nop()` (no-op, tidak ada file atau stdout)
- `APP_LOG=true` → tulis log ke file `./logs/{name}_YYYY-MM-DD.log`
- Pemanggil: `internal/db/db.go` (GORM query, name: `"db"`), `cmd/app/main.go` (server status, name: `"fiber"`)
- Di production (`APP_ENV != "local"`), GORM log level = `Warn`, jadi query normal tidak tercatat — hanya slow query & error
- Folder `./logs/` terbuat otomatis saat `lib.Log.Info/Error` pertama dipanggil
- File handle disimpan di memory (`logWriter.file`) — setiap panggilan `getOrCreate(name)` dengan nama prefix yang sama mengembalikan writer yang sudah ada
- `lib.Log.CloseFile(filePath string)` — method untuk menutup file handle sebelum file dihapus (diperlukan di Windows karena `os.Remove` gagal jika file masih terbuka). Memanggil `file.Close()` pada writer yang cocok dengan `filePath` dan menghapusnya dari map

---

## 8. Internationalization

Custom i18n system at `internal/lang/`:

- `t._("key")` → returns translated string
- `t._("key", map[string]any{"arg": "value"})` → replaces `:arg` in string
- Currently only `en` locale exists
- `APP_LOCALE` controls active language

---

## 9. Code Conventions

### Go:

- **HTTP layer** berada di `internal/http/`:
  - **Controllers:** `internal/http/controllers/` — exported functions that delegate to repositories.
  - **Repositories:** `internal/http/repositories/` — exported functions; **1 file per API endpoint** with specific name (`RoleListRepository`, `PermissionStoreRepository`).
  - **Resources:** `internal/http/resources/` — exported functions `Single()` dan `Collection()`, bukan struct method. Setiap field `id` wajib di-encode dengan `hash.EncodeId()`.
  - **Request schemas:** `internal/http/request/` — struct validation via go-playground/validator.
- **Middleware:** Exported function + registered via `appProvider`.
- **Handlers:** Return `c.JSON()` responses via helper functions.

### Frontend (Svelte/TypeScript):

- **Router:** `@keenmate/svelte-spa-router` with **history mode** (`setHashRoutingEnabled(false)`). Hash redirect fallback converts `/#/xyz` → `/xyz`.
- **Route definitions:** Map `{ '/': Home, '/about': About, '*': NotFound }` or `defineRoutes()` for type-safe named routes. Catch-all `'*'` matches unmatched paths.
- **Navigation links:** `<a href="/path" use:link>` for SPA navigation. Programmatic: `push(path)` / `nav.about.push()`.
- **Route params:** `{ path: '/user/:id' }` → component receives `let { routeParams = {} } = $props()` — access via `routeParams.id`.
- **Query string helpers:** `query<T>()` for reactive query params, `updateQuerystring()` for modifications.
- **Navigation guards:** `registerBeforeLeave()` for dirty-form protection.
- **Permissions (frontend):** Built-in RBAC with `createProtectedRoute()`, `hasPermission()`.
- **API hooks:** `web/src/api/` berisi wrapper TanStack Query per endpoint. Naming: `{method}{Name}` (getUser, postLogin, etc). Barrel per subfolder.
- **TanStack Query pattern (Svelte 5 runes):**
  - `@tanstack/svelte-query` v6 — `createQuery` dan `createMutation` dipanggil via `useQuery`/`useMutation` wrapper (`internal/lib/tanstackUtil.ts`) dengan `() => QueryClient` untuk shared client
  - Setiap API hook file menggunakan interface `IProps` dengan `param` (path params) dan `query` (query string params)
  - Query hooks menggunakan **getter pattern**: `() => IProps` — karena `createBaseQuery.svelte.js` memanggil `options()` di dalam `$derived.by`, getter diperlukan agar Svelte 5 bisa melacak pembacaan `$state`
  - `createQueryStr(props)` dari `axiosLib.ts` untuk membangun query string dari `IProps.query` secara otomatis
  - Mutation hooks (`useMutation`) dibuat via factory function: `export const deleteLog = () => useMutation(...)` — dipanggil sekali di komponen, hasilnya disimpan untuk digunakan `mutate()`
  - `queryKey` mencakup semua variable yang mempengaruhi hasil query (`file_name`, `search`, `page`, `levels`) — beda param = beda cache
  - `enabled: !!props.param.file_name` — disable query saat param belum siap (e.g. `selectedFile` masih `null`)
- **TypeScript:** Strict mode.
- **Formatting:** Prettier with `prettier-plugin-svelte` + `prettier-plugin-tailwindcss`.
- **CSS:** Tailwind CSS v4 via `@tailwindcss/vite` plugin — import `@import "tailwindcss"` di `app.css`, class-based styling, otomatis sort via Prettier.
- **Imports:** Use relative paths within `web/src/`.

### Tooling & IDE:

- **VS Code:** `editor.formatOnSave: true`, Svelte formatter set ke Prettier via `.vscode/settings.json`.
- **Go hot-reload:** Air (`air`) configured in `.air.toml` (Windows-compatible: `.exe` extension). Watches `cmd/` + `internal/`, builds `./cmd/app`.

---

## 10. NPM Scripts (SPA)

| Script    | Command                                                                    |
| --------- | -------------------------------------------------------------------------- |
| `dev`     | `vite --host` (port 5173 / configured)                                     |
| `build`   | `vite build` → output `../../public/build/`                                |
| `preview` | `vite preview --host`                                                      |
| `check`   | `svelte-check --tsconfig ./tsconfig.app.json && tsc -p tsconfig.node.json` |
| `format`  | `prettier --write 'src/**/*.{svelte,ts,css}'`                              |

---

## 11. Go Scripts

| Script    | Command                                 | Port |
| --------- | --------------------------------------- | ---- |
| `dev`     | `go run ./cmd/app` / `air` (hot-reload) | 8000 |
| `build`   | `go build -o ./tmp/app ./cmd/app`       | -    |
| `run`     | `./tmp/app`                             | 8000 |
| `migrate` | `go run ./cmd/migrate`                  | -    |
| `lint`    | `golangci-lint run ./...`               | -    |
| `test`    | `go test ./...`                         | -    |

Hot-reload via `air` (`.air.toml`): watches `cmd/` + `internal/`, rebuilds to `./tmp/main.exe` from `./cmd/app`.

Binary output ke `./tmp/` (terdaftar di `.gitignore`).

### CLI Commands

- **`go run ./cmd/app`** — Start API server
- **`go run ./cmd/migrate`** — Run database migrations (AutoMigrate)
- **`go run ./cmd/seed`** — Run database seeder

 ---

## 12. Log Viewer

### Overview

Log viewer untuk membaca log file JSON yang dihasilkan oleh `internal/lib/logger.go`. Backend membaca file langsung dari filesystem (`./logs/`) tanpa database. Frontend menyediakan UI untuk browsing, filter, search, pagination, expand/collapse log entries, download, dan delete.

### Backend — Log API (`/api/log`)

Semua endpoint di bawah grup `/api/log`, public (tanpa auth), tag OpenAPI `"Log"`.

| Method | Path | Repository | Description |
|--------|------|-----------|-------------|
| GET | `/api/log` | `log_list_repository.go` | List file `.log` di folder `./logs/`, return `[{name, size}]` sorted descending |
| GET | `/api/log/:filename` | `log_detail_repository.go` | Baca file, filter level, search, pagination via `helper.Res.Paginate` |
| DELETE | `/api/log/:filename` | `log_delete_repository.go` | Hapus file (close handle via `lib.Log.CloseFile` dulu, baru `os.Remove`) |
| GET | `/api/log/:filename/download` | `log_download_repository.go` | Download file via `c.SendFile` |

**Log Detail Query Params:**

| Param | Type | Description |
|-------|------|-------------|
| `levels` | string (comma-separated) | Filter by level (e.g. `info,error`) |
| `search` | string | Text search di field `message` |
| `page` | integer | Page number (default 1) |
| `limit` | integer | Items per page (default 50) |

**Response format (paginated):**
```json
{
  "message": "...",
  "data": [{ "level": "info", "time": "2026-06-09 12:00:00", "message": "..." }],
  "meta": { "total": 100, "page": 1, "limit": 50 }
}
```

### Backend — Delete File Handling

Di Windows, `os.Remove` gagal jika file masih terbuka oleh proses yang sama. Log writer (`lib/logger.go`) menyimpan `*os.File` handle di memory. Saat delete:
1. `LogDeleteRepository` panggil `lib.Log.CloseFile("./logs/" + filename)` — search di `l.files` map, jika ada writer dengan `filePath` cocok, `file.Close()` dan `delete(l.files, name)`
2. `os.Remove("./logs/" + filename)` — berhasil karena handle sudah ditutup

### Frontend — API Hooks (`web/src/api/logger/`)

Menggunakan `@tanstack/svelte-query` v6 dengan pattern getter untuk Svelte 5 reactivity.

| File | Type | IProps | Query Key |
|------|------|--------|-----------|
| `getLogs.ts` | useQuery | `{ param: {}, query: {} }` | `["log"]` |
| `getLogDetail.ts` | useQuery | `{ param: { file_name }, query: { search, page, limit }, levels }` | `["log", "detail", file_name, search, page, levels]` |
| `deleteLog.ts` | useMutation (factory) | `{ param: { file_name } }` | — |

**Pattern:**
- Query hooks menerima getter `() => IProps`, dipanggil oleh `createQuery` di dalam `$derived.by` → Svelte 5 track `$state` reads
- `enabled: !!props.param.file_name` — query tidak jalan jika `file_name` kosong
- `queryKey` mencakup semua variable yang mempengaruhi hasil → cache terisolasi per kombinasi filter
- `createQueryStr(props)` dari `axiosLib.ts` membangun query string dari `IProps.query`

### Frontend — Logger Page (`web/src/pages/guest/Logger.svelte`)

Layout split: sidebar (file list) + main area (entries).

| Component | Description |
|-----------|-------------|
| **Sidebar** | List file `.log` dengan nama + size (format human-readable via `formatSize`). Klik untuk select file |
| **Search bar** | Input text dengan debounce 300ms. Clear button. Memicu `currentPage = 1` |
| **Level filter** | Tombol `info` (blue) dan `error` (red). `Record<string, boolean>` — toggle tiap level independent. Semua mati → tidak ada filter (semua level tampil) |
| **Log entries** | Tampilkan `level` badge, `time`, `message` (clamp 2 baris). Klik untuk expand/collapse — `expandedIndex` state |
| **Expanded content** | `<pre>` dengan `whitespace-pre-wrap` di dalam div `flex-1` message column, border left untuk visual indentasi |
| **Pagination** | Prev/Next + page numbers. `pageNumbers()` function — format compact untuk banyak halaman (with `...`). Detail dihitung client-side: `Math.ceil(total / limit)` |
| **Action buttons** | Refresh (`detailQuery.refetch()`), Download (`<a>` dengan `instance.defaults.baseURL`), Delete (`deleteLog.mutate()` dengan `onSuccess` callback) |

### Level Filter State

```typescript
let selectedLevels = $state<Record<string, boolean>>({
  info: true,
  error: true,
})

function toggleLevel(level: string) {
  selectedLevels = { ...selectedLevels, [level]: !selectedLevels[level] }
  currentPage = 1
}
```

Nilai levels dikirim sebagai query param: `Object.keys(selectedLevels).filter(k => selectedLevels[k]).join(',')`.

---

## 13. Production

### Build & Deploy

```bash
# Deploy via Docker (SPA + Go build di dalam multi-stage build)
docker build -t portfolio -f build/package/Dockerfile .

# Run container
docker run -itd \
  --name portfolio \
  --restart always \
  -p 8000:8000 \
  --network proxy \
  --env-file /path/to/.env \
  portfolio
```

Atau via `scripts/deploy.sh` yang menjalankan `docker build` + `docker run` + cleanup.

### main.go production logic (`cmd/app/main.go`)

```go
// Global error handler (ValidationError, fiber.Error, generic)
app := fiber.New(fiber.Config{
    ErrorHandler: provider.NewErrorHandler(),
})

// CORS — allow frontend dev server (APP_URL)
app.Use(cors.New(config.CORSConfig()))

// Static assets (CSS/JS from Vite build)
app.Static("/", "public")

// SPA catch-all — serve index.html for all non-API paths
app.Get("/api/*", func(c *fiber.Ctx) error {
    return c.Status(fiber.StatusNotFound).JSON(helper.Res.Error("API endpoint not found", nil))
})
app.Get("/*", func(c *fiber.Ctx) error {
    return c.SendFile("public/index.html")
})
```

### Docker

Multi-stage build (`build/package/Dockerfile`):

```
Stage 1 (node:24-alpine) → npm install && npm build (from web/) → public/build/
Stage 2 (golang:1.25-alpine) → go build ./cmd/app → app
Stage 3 (alpine:3.19) → copy app + public/ (index.html, favicon, dll) + public/build/ → ./app (port 8000)
```

SPA build mendukung `ARG VITE_API_URL` / `VITE_APP_URL` untuk env vars yang di-inject via `--build-arg`.

`.dockerignore` tersedia untuk memperkecil build context (ignore node_modules, tmp, .git, dll).

### SPA Development Mode

When `APP_ENV=local`, SPA runs separately on `:3000` (Vite dev server). Production uses built files served by Go Fiber from `public/`. To start SPA dev: `cd web && pnpm dev`.

### SPA Router

| Package             | `@keenmate/svelte-spa-router` (fork with Svelte 5 runes)                          |
| ------------------- | --------------------------------------------------------------------------------- |
| **Mode**            | History mode via `setHashRoutingEnabled(false)` + `setBasePath('/')` in `main.ts` |
| **Hash redirect**   | `/#/xyz` → `history.replaceState` to `/xyz` on init                               |
| **Link navigation** | `<a href="/path" use:link>` or programmatic `push(path)`                          |
| **Params**          | `{ path: '/user/:id' }` → `routeParams.id` in component                           |
| **Catch-all**       | `'*': NotFound` in route map for 404 pages                                        |

### GitHub Actions Deploy

Workflow: `.github/workflows/deploy.yml`

Trigger: push ke `master` atau manual `workflow_dispatch`.

```yaml
# Single job deploy:
# 1. SSH ke VPS
# 2. cd /root/docker/go-fiber-svelte
# 3. git reset --hard + git pull
# 4. cd web && pnpm i --prod && cd ..
# 5. chmod +x & bash ./scripts/deploy.sh
```

Deploy script (`scripts/deploy.sh`):

1. `source .env` (baca VITE\_\* vars)
2. `docker build --build-arg VITE_* -t portfolio -f build/package/Dockerfile .`
3. Stop + remove container lama
4. `docker run --env-file .env ... portfolio`
5. `docker image prune -f`

### API Documentation (OpenAPI / Scalar)

- OpenAPI 3.0 spec generated dynamically at `/api/openapi.json`
- Per-controller `*OpenAPIPaths()` functions merged in `internal/openapi/openapi.go`
- Scalar UI served at `/api/docs` using vendored `public/scalar-standalone.js` (offline-capable)

### Go Middleware

**Auth middleware** (`internal/middleware/auth_middleware.go`): validates JWT cookie, sets `c.Locals("user_id")`. Applied **per-route** (not globally) via `internal/routes/api.go`:

```go
api := app.Group("/api")
authMw := AuthMiddleware{...}

// Public routes (no auth)
api.Post("/auth/login", authController.Login)
api.Get("/guest/ping", guestController.Ping)

// Protected routes (with auth)
auth := api.Group("/auth", authMw.Handle)
auth.Get("/user", authController.User)
```

Public API routes: `/api/openapi.json`, `/api/docs`, `/api/guest/ping`, `/api/auth/login`, `/api/log/*`. All others return 401 if unauthenticated.
