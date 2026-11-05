# 9. Data Layer

## 9.1 Principles

- New schema, owned by Go migrations.
- Keep Laravel table and column names where practical to reduce mapping cost and preserve API field names.
- Use explicit GORM tags. Never rely on GORM's pluralization guesses.
- No `AutoMigrate` in production. Migrations are applied by `golang-migrate`.
- All timestamps are stored in UTC; the application renders in `Asia/Jakarta`.
- All money/decimal values use `decimal.Decimal` or integer minor units, never `float64`, unless the frontend already receives floats.

## 9.2 Context connections

The Laravel app has 17 named connections, all PostgreSQL. The Go rewrite should not copy all of them blindly.

| Connection                                                                                                                                 | Recommendation                                                              |
| ------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------- |
| `pgsql` (main)                                                                                                                             | Core database                                                               |
| `kypas`, `tgkypas`, `tgcast`, `intra`, `kampiun`, `windigital`, `newpresensi`, `lemon`, `evolution`, `tgsso`, `auth`, `action`, `kypasbiz` | Evaluate per context; merge into core when the data is owned by the rewrite |
| `rfc`, `new_cuti`                                                                                                                          | Keep as read/write adapters if the feature remains                          |
| `sygap_sit`, `sygap_prod`, `sygap_view`                                                                                                    | External Sygap databases; use dedicated read/write clients                  |

All connections use PostgreSQL; only the database name, schema, and credentials differ.

Implement a connection resolver:

```go
type Resolver struct {
    core *gorm.DB
    byName map[string]*gorm.DB
}

func (r *Resolver) Core() *gorm.DB
func (r *Resolver) Named(name string) (*gorm.DB, error)
```

Place the resolver in `internal/platform/database`.

## 9.3 Repository and criteria

- Repository interfaces live in `internal/app/Application/<ctx>/repository`, next to their GORM implementations.
- Domain packages contain entities, enums, and policies only; there is no `app/domain/<ctx>/repository` or `app/domain/<ctx>/contract` package.
- Contracts live in the layer that owns them: repository interfaces in `internal/app/Application/<ctx>/repository`, integration client interfaces in `app/infrastructure/integrations/<name>/contracts`.
- Filters live in `internal/app/Application/<ctx>/criteria` and are applied explicitly.
- Example criteria:

```go
type UserFilter struct {
    Search   *string
    Role     *string
    Division *uint64
}

func (f UserFilter) Apply(db *gorm.DB) *gorm.DB {
    if f.Search != nil {
        like := "%" + *f.Search + "%"
        db = db.Where("name LIKE ? OR username LIKE ? OR email LIKE ?", like, like, like)
    }
    if f.Role != nil {
        db = db.Where("role = ?", *f.Role)
    }
    if f.Division != nil {
        db = db.Where("division_id = ?", *f.Division)
    }
    return db
}
```

## 9.4 Transactions

Use a single transaction per use case:

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&attendance).Error; err != nil {
        return err
    }
    return tx.Model(&user).Update("last_login_at", time.Now()).Error
})
```

## 9.5 Migrations

- Baseline: port the existing migrations to SQL files, replacing the Passport `oauth_*`, personal access token, and impersonation session tables with the single `auths` table in [Section 8.3](08-authentication.md#83-database-table).
- Naming: `000001_create_users_table.up.sql` / `.down.sql`.
- Keep migration order and foreign keys consistent.
- Seeders: roles, permissions, menus, settings, email templates, and admin user. Seeders are idempotent and run through `console db seed`; initial auth data is defined in [Section 8.3.1](08-authentication.md#831-initial-auth-data).
- Create an ERD section or link to one in `docs/`.
- Add indexes for:
    - `users.username`, `users.email`
    - `attendances.user_id`, `attendances.date`
    - `activity_log.subject_type/subject_id`
    - `auths.token`, `auths.refresh_token`, `auths.user_id`
    - `user_has_roles.role_id`, `role_has_permissions.permission_id`
    - `mcp_audit_logs.user_id`, `mcp_audit_logs.created_at`
    - Every foreign key used by list filters

See [Section 9.7](09-data-layer.md#97-database-commands) for the command surface.

## 9.6 Data ownership

The Go service is the only writer for the new core database. External PostgreSQL databases (`sygap_*`, `rfc`, `new_cuti`) are accessed through adapters, never through cross-database joins.

## 9.7 Database commands

All database changes run as one-off `console` commands. The `app` process never migrates or seeds on startup.

| Laravel                             | Go console command                 | Behavior                                                                           |
| ----------------------------------- | ---------------------------------- | ---------------------------------------------------------------------------------- |
| `php artisan migrate`               | `console migrate`                  | Apply all pending migrations; fail if the schema is dirty                          |
| `php artisan migrate:status`        | `console migrate status`           | Print the current version, pending count, and dirty flag; exit non-zero when dirty |
| `php artisan migrate:rollback`      | `console migrate rollback [steps]` | Roll back the last N migrations (default 1)                                        |
| `php artisan migrate:reset`         | `console migrate reset`            | Roll back all migrations in the current database                                   |
| `php artisan migrate:refresh`       | `console migrate refresh`          | Reset then migrate; development and staging only                                   |
| `php artisan migrate:fresh`         | `console migrate fresh`            | Drop all tables then migrate; requires `--force` outside development and test      |
| `php artisan migrate --seed`        | `console migrate --seed`           | Migrate, then run the seeders                                                      |
| `php artisan make:migration <name>` | `console migrate create <name>`    | Create the next `NNNNNN_<name>.up.sql` and `.down.sql` pair                        |
| `php artisan migrate:install`       | Not needed                         | `golang-migrate` creates `schema_migrations` automatically                         |
| `php artisan migrate --pretend`     | Not supported                      | Review the SQL files and run `console migrate status` instead                      |
| `php artisan db:seed`               | `console db seed [seeder]`         | Run all seeders or one named seeder                                                |
| `php artisan db:wipe`               | `console db wipe`                  | Drop all tables; requires `--force` outside development and test                   |
| `php artisan db:show`               | `console db show`                  | Print connection, PostgreSQL version, and table count                              |
| `php artisan db:monitor`            | `console db monitor`               | Check the core connection and connection pool limits                               |
| Spatie `backup:run --only-db`       | `console db backup`                | Write a compressed `pg_dump` custom-format dump                                    |
| Manual restore                      | `console db restore <file>`        | Restore a `pg_dump` dump; manual and guarded by `--force`                          |
| `php artisan tinker`                | Not ported                         | Use tests or a dedicated `console` command                                         |

Implementation rules:

- Read all database settings from `.env` only: `DB_CONNECTION=pgsql`, `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`, and `DB_SSLMODE`. Never hardcode credentials in code, migrations, Dockerfiles, or committed files.
- If `DATABASE_URL` is set, it overrides the split `DB_*` values when building the migration database URL; otherwise build the URL from `DB_*`.
- Use `golang-migrate` with the PostgreSQL driver and `-path ./migrations`; do not use GORM `AutoMigrate` in production.
- Register the commands in `cmd/console` and keep migration logic in `internal/platform/database`.
- Run exactly one migration writer at a time. During deployment, run `console migrate` as an init job before the app replicas start; never migrate from the `app` process.
- PostgreSQL supports transactional DDL; keep each migration atomic where possible. If a migration fails, `golang-migrate` marks the version dirty (`schema_migrations.dirty = 1`). Fix the database manually, run `console migrate force <version>`, then run `console migrate` again.
- `console migrate force <version>` is a manual recovery command; never call it from automation.
- Destructive commands (`fresh`, `reset`, `refresh`, `wipe`, `restore`) require `--force` outside `development` and `test`.
- Seeders must be idempotent: use `INSERT ... ON CONFLICT (...) DO UPDATE SET ...` or existence checks so repeated runs do not duplicate roles, permissions, menus, settings, email templates, or the admin user.
- `console db show` and `console db monitor` are read-only diagnostics.
- Write both the `.up.sql` and `.down.sql` side for every migration; an empty down migration is only acceptable for irreversible data fixes and must be noted in the file.
- Typical deployment sequence:

```bash
console migrate status
console migrate
console db seed        # first install only
```

## 9.8 MCP audit table

The MCP bearer token lives in `.env` (`MCP_BEARER_TOKEN`); no token table is needed. Only tool calls are persisted.

```sql
CREATE TABLE mcp_audit_logs (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users (id) ON DELETE SET NULL,
    tool        VARCHAR(100) NOT NULL,
    status      VARCHAR(20) NOT NULL,
    duration_ms INTEGER NOT NULL,
    ip          INET,
    args_hash   CHAR(64),
    error       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX mcp_audit_logs_user_created_idx ON mcp_audit_logs (user_id, created_at);
```

- `user_id` is `MCP_ACTOR_USER_ID` when configured, otherwise null.
- `args_hash` is a hash of the normalized tool arguments; audit logs do not store raw values that may contain personal data.
