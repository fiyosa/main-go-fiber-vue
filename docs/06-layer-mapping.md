# 6. Layer Mapping: Laravel → Go

## 6.1 Summary table

| Laravel                                       | Go equivalent                                            | Notes                                                                                                         |
| --------------------------------------------- | -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| Eloquent model                                | GORM struct in `app/domain/<ctx>/entity`                     | Explicit `TableName()` and column tags                                                                        |
| Eloquent relationship                         | GORM relation or explicit repository query               | Avoid lazy loading                                                                                            |
| Cast                                          | Custom Go type / `Scanner`+`Valuer`                      | No legacy encryption compatibility required                                                                   |
| Enum                                          | Typed string/int + constants                             | Values must match existing JSON                                                                               |
| Observer                                      | `AfterCreate`/`AfterUpdate` hooks or usecase events      | Keep side effects explicit                                                                                    |
| Controller                                    | Controller struct + method                                | One controller file per module                                                                                |
| FormRequest                                   | Request struct + validator tags + custom rules           | Errors mapped to Laravel validation shape                                                                     |
| API Resource                                  | Response DTO + mapper                                    | Preserve JSON keys and nested shapes                                                                          |
| Middleware                                    | Fiber middleware in `internal/platform/http/middleware`  | `auth`, `guest`, `can`, `throttle`, `verify.pin`, `verify.header`, `telegram.auth`                            |
| Repository (Prettus)                          | Repository interface + GORM implementation               | Both live in `internal/app/Application/<ctx>/repository`; domain has no repository package                                 |
| Integration client contract                   | `app/infrastructure/integrations/<name>/contracts/client.go` | Adapter-local port; inject it into resource clients                                                           |
| Unused domain API contract                    | Not ported                                               | `SygapApiContract` has no implementation or callers                                                           |
| Criteria                                      | Typed filter struct + GORM scopes                        | Applied explicitly by the usecase/controller                                                                  |
| Policy                                        | `Authorizer` service                                     | Super admin bypass preserved                                                                                  |
| Job                                           | Goroutine job function                                   | Started via `jobs.Go` from the owning use case                                                                |
| Console command                               | Cobra subcommand                                         | Same command names where practical                                                                            |
| Scheduler entry                               | robfig/cron entry                                        | Runs the handler directly in a goroutine; no enqueue                                                         |
| Event/Listener                                | In-process event bus or explicit goroutine job           | No retry machinery; handlers must be idempotent                                                              |
| Blade view (PDF/Excel)                        | `html/template` or excelize                              | See [Section 12.6](12-platform-services.md#126-pdf)                                                                                              |
| Config                                        | Env struct                                               | No `.env` parsing inside the app                                                                              |
| Helper functions                              | `app/infrastructure/support/helpers`                | `sendResponse`, `sendError`, `logActivity`, etc.                                                              |
| `Infrastructure/External/AuthSSO.php`         | `app/infrastructure/external/auth_sso.go`                    | Preserve the SSO behavior described in [Section 8.9](08-authentication.md#89-sso-callback)                                                            |
| `Infrastructure/Integrations/<Name>`          | `app/infrastructure/integrations/<name>`                     | Keep `clients`, `contracts`, `dto/requests`, `dto/responses`, `exceptions`, and `support`                     |
| `Infrastructure/Rules/ValidImageUrl.php`      | `app/infrastructure/rules/valid_image_url.go`                | Custom validation rule                                                                                        |
| `Infrastructure/Services/Approval`            | `app/infrastructure/services/approval`                       | Keep `concerns` and `handlers`                                                                                |
| `Infrastructure/Services/LdapService.php`     | `app/infrastructure/services/ldap_service.go`                | LDAP adapter                                                                                                  |
| `Infrastructure/Services/MailService.php`     | `app/infrastructure/services/mail_service.go`                | Mail adapter                                                                                                  |
| `Infrastructure/Services/TelegramService.php` | `app/infrastructure/services/telegram_service.go`            | Telegram adapter                                                                                              |
| `Infrastructure/Support/Helpers`              | `app/infrastructure/support/helpers`                         | Response, activity log, pagination, hash ID, and date helpers                                                 |
| `Infrastructure/Support/Traits`               | `app/infrastructure/support/traits`                          | Convert to helper functions or embedded structs                                                               |
| `Infrastructure/Laravel`                      | Not a Go package                                         | Framework glue maps to `bootstrap`, `internal/platform/http/middleware`, `internal/app/Application/.../jobs`, and support helpers |

`app/infrastructure` mirrors `app/Infrastructure`. Go-only runtime adapters (`auth`, `database`, `cache`, `http/middleware`, `jobs`, `storage`, `broadcast`) live in `internal/platform`; the cron scheduler and MCP tools live in `internal/app/infrastructure`.

`app/Infrastructure/Laravel` subfolder mapping:

| Laravel                                | Go                                        |
| -------------------------------------- | ----------------------------------------- |
| `Console`                              | `cmd/console`                             |
| `Events`, `Listeners`, `Notifications` | `internal/app/Application/<ctx>/jobs`           |
| `Exceptions`                           | `app/infrastructure/support/helpers`          |
| `Mail`                                 | `app/infrastructure/services/mail_service.go` |
| `Middleware`                           | `internal/platform/http/middleware`          |
| `Providers`                            | `bootstrap`                               |

## 6.2 Example: model and repository

```go
// app/domain/core/entity/user.go
package entity

import "time"

type User struct {
    ID        uint64    `gorm:"column:id;primaryKey"`
    Email     string    `gorm:"column:email;size:255;not null;uniqueIndex"`
    Password  string    `gorm:"column:password;size:255;not null"`
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string { return "users" }
```

```go
// internal/app/application/core/repository/user/repository.go
package user

import (
    "context"

    "go-fiber-vue/internal/app/domain/core/entity"
)

type Repository interface {
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
```

```go
// internal/app/application/core/repository/user/repository.go
package user

import (
    "context"
    "errors"

    "gorm.io/gorm"

    "go-fiber-vue/internal/app/domain/core/entity"
)

type gormRepository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    if db == nil {
        return nil
    }
    return &gormRepository{db: db}
}

func (r *gormRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var found entity.User
    err := r.db.WithContext(ctx).
        Where("email = ?", email).
        First(&found).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &found, err
}
```

## 6.3 Example: request, controller, route

```go
// internal/app/presentation/core/request/auth_login_request.go
package request

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
```

```go
// internal/app/presentation/core/routes/auth_routes.go
package route

import (
    "github.com/gofiber/fiber/v3"

    "go-fiber-vue/internal/app/presentation/core/controllers"
)

func RegisterAuthRoutes(group fiber.Router, c *controllers.AuthController, auth fiber.Handler, guest fiber.Handler) {
    group.Get("/me", auth, c.Me)
    group.Post("/login", guest, c.Login)
    group.Post("/refresh", c.Refresh)
    group.Post("/logout", auth, c.Logout)
}
```

```go
// internal/app/presentation/core/controllers/auth_controller.go
package controllers

import (
    "github.com/gofiber/fiber/v3"

    "go-fiber-vue/internal/app/presentation/core/request"
)

type AuthController struct {
    authService AuthService
    validator   Validator
}

func (c *AuthController) Login(fCtx fiber.Ctx) error {
    var req request.LoginRequest
    if err := c.validator.BindAndValidate(fCtx, &req); err != nil {
        return err
    }

    result, err := c.authService.Login(fCtx.Context(), req)
    if err != nil {
        return err
    }

    return fCtx.JSON(result)
}
```

## 6.4 Response envelope

The Laravel helpers in `app/Infrastructure/Support/Helpers/Function.php` define the base envelope. Go preserves the shape, adds a translated `message` to every response, and does not use `204 No Content` for JSON deletes.

```go
// internal/app/infrastructure/support/helpers/response.go
package helpers

import "github.com/gofiber/fiber/v3"

type Envelope struct {
    Success bool   `json:"success"`
    Data    any    `json:"data,omitempty"`
    Message string `json:"message"`
}

func Success(c fiber.Ctx, data any, key string, params map[string]string, code int) error {
    message := lang.From(c).T(key, params)
    return c.Status(code).JSON(Envelope{Success: true, Data: data, Message: message})
}

func Error(c fiber.Ctx, key string, data any, code int) error {
    message := lang.From(c).T(key, nil)
    return c.Status(code).JSON(Envelope{Success: false, Data: data, Message: message})
}
```

Message rules:

- Every JSON response has a non-empty `message`; see [Section 20](20-localization.md).
- Helpers accept a translation key, never a literal string.
- `lang.From(c)` reads the resolved locale and translator from Fiber locals.
- `204 No Content` is not used for JSON deletes because it cannot carry `message`.

Compatibility quirks to preserve:

- `sendResponse` always sets `success: true`, even when the HTTP status is `403`. Preserve this behavior for existing clients.
- Do not infer the status from translated text. Callers pass the HTTP status explicitly; missing resources return `404` through `common.not_found` or `lang.not_found`.
- `sendSuccess` returns `{success: true, message}` without `data`.
- Pagination is provided by `pagination($query, $limit)` and `sendResponsePaginate`. Preserve the Laravel paginator JSON shape (`data`, `links`, `meta`, `current_page`, `per_page`, `total`, ...).
