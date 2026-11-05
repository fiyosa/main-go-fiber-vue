# 4. Target Go Architecture

## 4.1 Stack

| Concern     | Choice                                                           | Notes                                                                                  |
| ----------- | ---------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| HTTP        | `github.com/gofiber/fiber/v3` (latest stable)                    | Pin the exact version in `go.mod`; verify v3 API before coding                         |
| MCP         | `github.com/modelcontextprotocol/go-sdk/mcp`                     | Streamable HTTP at `/api/mcp` on the same Fiber port; bearer token auth                |
| ORM         | `gorm.io/gorm` + `gorm.io/driver/postgres`                       | Explicit table/column tags; no `AutoMigrate` in production                             |
| Migrations  | `github.com/golang-migrate/migrate/v4`                           | SQL-first migrations                                                                   |
| Jobs        | Go goroutines (`internal/platform/jobs`)                          | Fire-and-forget; no queue, broker, or worker process                                   |
| Scheduler   | `github.com/robfig/cron/v3`                                      | Lives in `internal/app/infrastructure/scheduler`; started as a goroutine; one leader across replicas |
| JWT         | `github.com/golang-jwt/jwt/v5`                                   | HS256 with a single shared `JWT_SECRET`; issuer `app-api`, audience `app-spa`         |
| Validation  | `github.com/go-playground/validator/v10`                         | Request structs + custom rules                                                         |
| Config      | `github.com/kelseyhightower/envconfig`                           | 12-factor env configuration                                                            |
| Logging     | `log/slog`                                                       | JSON handler in production                                                             |
| Cache       | PostgreSQL `cache` table                                         | Shared across replicas; invalidate on write                                            |
| LDAP        | `github.com/go-ldap/ldap/v3`                                     | Directory sync only; not used for login                                                |
| Excel       | `github.com/xuri/excelize/v2`                                    | Import/export parity                                                                   |
| PDF         | `github.com/go-pdf/fpdf` or `chromedp`                           | See [Section 12.6](12-platform-services.md#126-pdf)                                                                       |
| Mail        | `github.com/wneessen/go-mail`                                    | SMTP + templates                                                                       |
| Telegram    | `github.com/go-telegram-bot-api/telegram-bot-api/v5`             | Existing bot tokens                                                                    |
| Pusher      | `github.com/pusher/pusher-http-go/v5`                            | Broadcast parity                                                                       |
| AWS         | `github.com/aws/aws-sdk-go-v2`                                   | S3 storage                                                                             |
| Hash IDs    | `github.com/speps/go-hashids/v2`                                 | Same salt/length via env                                                               |
| Testing     | `github.com/stretchr/testify`, `testcontainers-go`, `go-sqlmock` | See [Section 16](16-testing-strategy.md)                                                                         |
| Concurrency | `golang.org/x/sync/errgroup`                                     | Start API and scheduler in one process with shared cancellation                        |
| CLI         | `github.com/spf13/cobra`                                         | `app` (single process) and `console` (one-off commands: migrate, seed, imports, syncs) |

## 4.2 Runtime process

```text
┌───────────────────────────────────────────────────────┐
│ cmd/app (single process)                              │
│                                                       │
│  ┌──────────────┐  ┌──────────────────┐               │
│  │ Fiber HTTP   │  │ cron scheduler   │               │
│  │ (goroutine)  │  │ (goroutine)      │               │
│  └──────┬───────┘  └────────┬─────────┘               │
│         │                   │                         │
│         │            spawn job goroutines (on demand) │
└─────────┼───────────────────┼─────────────────────────┘
          │                   │
      PostgreSQL          PostgreSQL
      (API)               (advisory lock)
```

- `app`: one process that starts Fiber and the cron scheduler as goroutines under one context.
- `console`: one-off commands for imports, syncs, migrations, and maintenance; never a long-running service.
- Background jobs are goroutines spawned on demand by the owning use case; there is no queue or worker process.
- MCP reuses the Fiber port at `/api/mcp` (Streamable HTTP mounted through the `net/http` adaptor) and authenticates with a bearer token.
- All components share the same config, DB, and logger instances.
- Use `errgroup.WithContext` for the lifecycle. If any component returns a fatal error, the shared context is canceled and the process exits non-zero so the orchestrator can restart it.
- Job goroutines are not errgroup members; the `jobs` helper tracks them and `jobs.Wait` drains them during shutdown.
- Job handlers must recover panics; a bad job must not kill the API.
- Multiple replicas: all replicas serve HTTP; only the scheduler leader runs cron, guarded by a PostgreSQL advisory lock.

```go
// internal/bootstrap/runtime.go
func Run(ctx context.Context, deps Dependencies) error {
    group, ctx := errgroup.WithContext(ctx)

    group.Go(func() error { return runAPI(ctx, deps) })
    group.Go(func() error { return runScheduler(ctx, deps) })

    return group.Wait()
}
```

On graceful shutdown, `runAPI` treats the server-closed error (for example, `http.ErrServerClosed` or Fiber's equivalent) as a clean exit, the scheduler returns nil after its context is canceled, and `jobs.Wait` drains in-flight job goroutines with a bounded timeout.

Keep the `--disable-scheduler` flag for local debugging only; production defaults to both components enabled.
