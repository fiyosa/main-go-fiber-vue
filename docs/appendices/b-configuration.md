# Appendix B: Environment Variable Names

Only names are listed. Never commit values.

## Application

- `APP_ENV`
- `APP_NAME`
- `APP_URL`
- `APP_TIMEZONE`
- `APP_LOCALE` (default `en`; supported `en`, `id`)
- `LOG_LEVEL`

## HTTP

- `HTTP_HOST`
- `HTTP_PORT`
- `HTTP_READ_TIMEOUT`
- `HTTP_WRITE_TIMEOUT`
- `HTTP_BODY_LIMIT`
- `CORS_ALLOWED_ORIGINS`
- `TRUSTED_PROXIES`

## Database

Set `DB_CONNECTION=pgsql`. All database credentials and endpoints come from `.env` only; commit `.env.example` with empty placeholders and never commit `.env`.

- `DB_CONNECTION` (`pgsql`)
- `DB_HOST`
- `DB_PORT` (`5432`)
- `DB_DATABASE`
- `DB_USERNAME`
- `DB_PASSWORD`
- `DB_SSLMODE` (`disable` locally, `require` in production)
- `DATABASE_URL` (optional; overrides the split `DB_*` values for migrations)
- `DB_MAX_OPEN_CONNS`
- `DB_MAX_IDLE_CONNS`
- `DB_CONN_MAX_LIFETIME`
- `DB_KYPAS_*`, `DB_TGKYPAS_*`, `DB_SYGAP_*`, `DB_RFC_*`, `DB_NEW_CUTI_*` (only if those connections remain)

```dotenv
# .env.example
DB_CONNECTION=pgsql
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=alma
DB_USERNAME=
DB_PASSWORD=
DB_SSLMODE=disable
```

## Scheduler and cache

- `SCHEDULER_ENABLED`
- `SCHEDULER_LOCK_KEY`
- `CACHE_TABLE`

## MCP

- `MCP_ENABLED`
- `MCP_PATH` (default `/api/mcp`)
- `MCP_BEARER_TOKEN`
- `MCP_ACTOR_USER_ID`
- `MCP_SCOPES` (default `mcp:read`)
- `MCP_MAX_BODY_BYTES`
- `MCP_ALLOWED_ORIGINS`
- `MCP_SESSION_TIMEOUT`
- `MCP_STATELESS`

## Authentication

- `JWT_ISSUER`
- `JWT_AUDIENCE`
- `JWT_PRIVATE_KEY`
- `JWT_PUBLIC_KEY`
- `JWT_ACCESS_TTL`
- `JWT_REFRESH_TTL`
- `SEED_ADMIN_USERNAME`
- `SEED_ADMIN_EMAIL`
- `SEED_ADMIN_PASSWORD`

## LDAP

- `LDAP_CONNECTION`
- `OPEN_LDAP_SERVER`
- `OPEN_LDAP_BASE_DN`
- `OPEN_LDAP_USER`
- `OPEN_LDAP_PASS`
- `OPEN_LDAP_PORT`
- `LDAP_ADMIN_SERVER`
- `LDAP_ADMIN_BASE_DN`
- `LDAP_ADMIN_USER`
- `LDAP_ADMIN_PASS`

## SSO

- `SSO_SERVER_URL`
- `SSO_CLIENT_URL`
- `SSO_CALLBACK_URL`
- `SSO_CLIENT_KEY`
- `SSO_CLIENT_SECRET`

## Storage

- `FILESYSTEM_DISK`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_DEFAULT_REGION`
- `AWS_BUCKET`
- `AWS_BUCKET_TGKYPAS`
- `AWS_ENDPOINT`

## Mail

- `MAIL_MAILER`
- `MAIL_HOST`
- `MAIL_PORT`
- `MAIL_USERNAME`
- `MAIL_PASSWORD`
- `MAIL_ENCRYPTION`
- `MAIL_FROM_ADDRESS`
- `MAIL_FROM_NAME`

## Broadcasting

- `BROADCAST_DRIVER`
- `PUSHER_APP_ID`
- `PUSHER_APP_KEY`
- `PUSHER_APP_SECRET`
- `PUSHER_HOST`
- `PUSHER_PORT`
- `PUSHER_SCHEME`
- `PUSHER_APP_CLUSTER`

## Telegram

- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_BOT_TOKEN_<NAME>` for each bot in `TelegramBotEnum`

## Integrations

- `<INTEGRATION>_URL`
- `<INTEGRATION>_USERNAME`
- `<INTEGRATION>_PASSWORD`
- `<INTEGRATION>_APP_ID`
- `<INTEGRATION>_SECRET_ID`

Replace `<INTEGRATION>` with: `CASDOOR`, `DIARIUM`, `DOMINO`, `INSERA`, `KYPAS`, `MAILADMIN`, `MICROSOFT_GRAPH`, `MILISTADMIN`, `MYTENS`, `NDE`, `ONE_PORTAL`, `PORTAL`, `SYGAP`, `TGKYPAS`, `UNIVERSAL`, `WINDIGITAL`.

## Hashing

- `HASHIDS_SALT`
- `HASHIDS_LENGTH`
