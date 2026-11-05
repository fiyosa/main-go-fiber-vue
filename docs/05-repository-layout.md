# 5. Repository Layout and Naming

## 5.1 Target tree

```text
go-fiber-vue/          # module name: go-fiber-vue
├── cmd/
│   ├── app/main.go          # single process: API + scheduler
│   └── console/main.go      # one-off commands
├── internal/
│   ├── app/
│   │   ├── domain/
│   │   │   ├── core/
│   │   │   │   ├── entity/          # GORM models
│   │   │   │   ├── enum/
│   │   │   │   └── policy/
│   │   │   ├── sygap/
│   │   │   ├── kypas/
│   │   │   ├── tgkypas/
│   │   │   ├── kampiun/
│   │   │   ├── intra/
│   │   │   ├── rfc/
│   │   │   ├── windigital/
│   │   │   ├── tgcast/
│   │   │   ├── mailadmin/
│   │   │   └── ...
│   │   ├── Application/
│   │   │   ├── core/
│   │   │   │   ├── usecase/
│   │   │   │   ├── repository/      # interfaces + GORM implementations
│   │   │   │   ├── criteria/        # typed filters
│   │   │   │   ├── jobs/            # job functions (run as goroutines)
│   │   │   │   ├── export/
│   │   │   │   ├── import/
│   │   │   │   └── command/
│   │   │   └── ...
│   │   ├── infrastructure/
│   │   │   ├── mcp/
│   │   │   │   ├── server.go    # MCP server + tool registration
│   │   │   │   ├── handler.go   # Streamable HTTP handler + bearer wrapper
│   │   │   │   ├── auth.go      # env bearer token + constant-time compare
│   │   │   │   └── tools/       # one file per domain
│   │   │   ├── scheduler/       # robfig/cron wrapper + advisory lock + entry registry
│   │   │   ├── external/
│   │   │   │   └── auth_sso.go
│   │   │   ├── integrations/
│   │   │   │   ├── casdoor/
│   │   │   │   │   ├── client.go
│   │   │   │   │   ├── clients/
│   │   │   │   │   │   └── user_client.go
│   │   │   │   │   ├── contracts/
│   │   │   │   │   │   └── client.go
│   │   │   │   │   ├── dto/
│   │   │   │   │   │   ├── requests/
│   │   │   │   │   │   │   └── update_user_request.go
│   │   │   │   │   │   └── responses/
│   │   │   │   │   │       └── get_user_response.go
│   │   │   │   │   ├── exceptions/
│   │   │   │   │   │   └── casdoor_exception.go
│   │   │   │   │   └── support/
│   │   │   │   │       └── endpoint.go
│   │   │   │   ├── domino/
│   │   │   │   ├── insera/
│   │   │   │   ├── kypas/
│   │   │   │   ├── mailadmin/
│   │   │   │   ├── microsoftgraph/
│   │   │   │   ├── milistadmin/
│   │   │   │   ├── nde/
│   │   │   │   ├── oneportal/
│   │   │   │   ├── sygap/
│   │   │   │   ├── tgkypas/
│   │   │   │   ├── universal/
│   │   │   │   └── windigital/
│   │   │   ├── rules/
│   │   │   │   └── valid_image_url.go
│   │   │   ├── services/
│   │   │   │   ├── approval/
│   │   │   │   │   ├── concerns/
│   │   │   │   │   │   └── form_response.go
│   │   │   │   │   ├── handlers/
│   │   │   │   │   │   └── ...
│   │   │   │   │   ├── approval_notification_service.go
│   │   │   │   │   ├── approval_service.go
│   │   │   │   │   └── form_request_handler.go
│   │   │   │   ├── ldap_service.go
│   │   │   │   ├── mail_service.go
│   │   │   │   └── telegram_service.go
│   │   │   └── support/
│   │   │       ├── helpers/
│   │   │       │   ├── date_helper.go
│   │   │       │   ├── function.go
│   │   │       │   ├── media_helper.go
│   │   │       │   ├── string_helper.go
│   │   │       │   ├── user_helper.go
│   │   │       │   ├── response.go       # Go-only envelope
│   │   │       │   ├── activity_log.go   # Go-only
│   │   │       │   ├── pagination.go     # Go-only
│   │   │       │   └── hashids.go        # Go-only
│   │   │       └── traits/
│   │   │           ├── create_user_request.go
│   │   │           └── register_employee.go
│   │   └── presentation/
│   │       ├── core/
│   │       │   ├── controllers/
│   │       │   ├── request/
│   │       │   ├── response/
│   │       │   ├── routes/          # one file per module, exported for aggregation
│   │       │   ├── rule/
│   │       │   └── view/    # embedded SPA shell (spa.html)
│   │       ├── sygap/ ...
│   │       ├── public/
│   │       └── apiv1/
│   ├── bootstrap/
│   │   ├── app.go           # fiber app + dependency wiring; calls routes.RegisterAPI + routes.RegisterWeb
│   │   ├── database.go
│   │   ├── cache.go         # PostgreSQL cache store
│   │   ├── mcp.go           # wires internal/app/infrastructure/mcp
│   │   ├── runtime.go       # starts the HTTP server and scheduler
│   │   └── scheduler.go     # wires internal/app/infrastructure/scheduler
│   ├── routes/
│   │   ├── api.go               # /api group mounted as /api/v1
│   │   ├── mcp.go               # /api/mcp via net/http adaptor
│   │   ├── web.go               # /build/* assets + /* SPA fallback
│   │   └── v1/
│   │       └── apiv1.go         # aggregates exports from internal/app/presentation/core/routes
│   ├── config/
│   │   └── config.go
│   ├── platform/                # Go-only runtime adapters shared across layers
│   │   ├── auth/                # jwt, refresh tokens, impersonation
│   │   ├── database/
│   │   ├── cache/
│   │   ├── http/middleware/     # Fiber middleware: locale, maintenance
│   │   ├── jobs/                # goroutine helper: panic recovery, timeout, wait group
│   │   ├── lang/                # translator + locale loader
│   │   ├── storage/
│   │   └── broadcast/
│   └── lang/                    # en.json + id.json; loaded at startup
├── migrations/
├── web/                         # Vue 3 + TypeScript frontend source (built to public/build)
├── tests/
│   ├── contract/
│   ├── integration/
│   └── unit/
├── docs/
├── go.mod
├── go.sum
└── README.md
```

All Go application code lives under `internal/app/`; shared runtime adapters live under `internal/platform/`. Only `cmd/`, `migrations/`, `public/`, `web/`, and `tests/` remain at the repository root.

Every folder under `internal/app/infrastructure/integrations` uses the same subfolders as `casdoor`; add `token_manager.go` only when the upstream requires token exchange.

`app/Infrastructure/Laravel` is framework glue and is not a Go package; see [Section 6.1](06-layer-mapping.md#61-summary-table) for its folder mapping.

## 5.2 `alma` → `core` mapping

| Laravel                                                       | Go                                                   |
| ------------------------------------------------------------- | ---------------------------------------------------- |
| `app/Domain/Alma`                                             | `internal/app/domain/core`                      |
| `app/Application/Alma`                                        | `internal/app/Application/core`                 |
| `app/Presentasion/Admin/Alma`                                 | `internal/app/presentation/core`                |
| `resources/js/domain/alma`                                    | Out of scope; frontend stays `alma`                  |
| `App\Domain\Alma\Models\User`                                 | `internal/app/domain/core/entity/user.go` (`type User`)  |
| `App\Application\Alma\Repositories\...`                       | `internal/app/Application/core/repository/...`  |
| `App\Presentasion\Admin\Alma\Http\Controllers\UserController` | `internal/app/presentation/core/controllers/user_controller.go` |

## 5.3 Naming rules

| Kind                 | Rule                                                                    | Example                                       |
| -------------------- | ----------------------------------------------------------------------- | --------------------------------------------- |
| Folder               | lowercase, singular concept                                             | `controllers`, `usecase`, `entity`          |
| Package              | lowercase, one word                                                     | `controllers`, `repository`, `core`          |
| File                 | `snake_case.go`                                                         | `user_controller.go`, `attendance_repository.go` |
| Interface            | noun + `Repository` / `Service`                                         | `UserRepository`, `ApprovalService`           |
| Repository interface | `internal/app/Application/<ctx>/repository`, next to its GORM implementation         | `repository.UserRepository`                   |
| Integration contract | `internal/app/infrastructure/integrations/<name>/contracts`, next to its HTTP client | `contracts.Client`                            |
| Domain contract      | Not allowed; define the port in the consuming layer                     | —                                             |
| Contract test        | `tests/contract`                                                        | `TestAuthLoginContract`                       |
| Struct               | PascalCase                                                              | `User`, `AttendanceController`               |
| Function             | camelCase, verb-first                                                   | `FindByUsername`, `StartAttendance`           |
| Route file           | `<module>_routes.go`                                                    | `user_routes.go`, `attendance_routes.go`      |
| Request DTO          | `<Action><Entity>Request`                                               | `StoreUserRequest`, `LoginRequest`            |
| Response DTO         | `<Entity>Response`                                                      | `UserResponse`, `AuthUserResponse`            |
| Job function         | `<context>:<domain>:<action>`                                           | `core:telegram:send`, `sygap:incident:update` |
| Job group            | lowercase, hyphen-free                                                  | `notification`, `email`, `processing`         |

## 5.4 Identifier rename exceptions

Do **not** rename these even though they contain `alma`:

- Public API paths such as `/api/v1/...`.
- Database tables and columns (`users`, `attendances`, `activity_log`, `auths`).
- Permission strings (`manage_roles`, `attendances.can-edit`, ...).
- Product branding in UI, emails, and PDF output.
- Existing environment variable names where deployment already uses them; new variables use `CORE_` only when they are genuinely new.
