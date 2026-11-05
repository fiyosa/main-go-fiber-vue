# 12. Platform Services

Implementations live in `internal/platform` (auth, database, cache, jobs, storage, broadcast), `internal/app/infrastructure/scheduler`, and `app/infrastructure/services` (LDAP, mail, Telegram, approval).

## 12.1 Storage and media

- Disks: local (`storage/app`, `storage/pvc`, `public`) and S3.
- Media model keeps the same columns (`file_name`, `file_path`, `disk`, `visibility`, `mediable_*`).
- Public route: `GET /media/{file_name}`.
- Protected routes: download and delete under `/api/v1/media`.
- Use signed URLs for S3 where applicable.
- Validate MIME type and size on upload.
- Detect remote picture size only if the existing feature still requires it; otherwise simplify.

## 12.2 Mail

- SMTP client with connection pooling.
- Template source: `email_templates` table plus `html/template`.
- Render and send asynchronously in a goroutine.
- Attachments supported.
- Preserve sender names and subjects used today.

## 12.3 Telegram

- One bot client per configured bot (`TelegramBotEnum`).
- `SendTelegramMessageJob` becomes `core:telegram:send`.
- Preserve Markdown/HTML parse mode and chat IDs.
- Retry on network errors, not on `400` responses.

## 12.4 Broadcast

- Pusher client with the same app ID/key/secret/cluster.
- Event names and channel names are part of the frontend contract; do not rename.
- If Pusher is replaced later, keep an adapter interface.

## 12.5 Excel

- Use `excelize`.
- Port import validation and error messages carefully; the UI displays row-level validation errors.
- Export templates currently stored as files or generated dynamically must keep the same column order and headers.

## 12.6 PDF

Current exports use Blade views:

- `exports/attendances.blade.php`
- `exports/attendance-aux.blade.php`
- `exports/shift-changes.blade.php`

Options:

1. `chromedp`: render the existing HTML/CSS in headless Chrome. Best visual parity, but adds a browser dependency.
2. `go-pdf/fpdf`: pure Go, but does not render HTML/CSS. Requires reimplementing the layout.
3. External `wkhtmltopdf` or a sidecar service.

Recommendation: start with `chromedp` if visual parity matters; otherwise use `fpdf` for simple tables and accept layout differences only with stakeholder approval.

## 12.7 Settings

- `settings` table with `key`/`value`.
- Cache in the PostgreSQL `cache` table with invalidation on update.
- Helper: `settings.Get(ctx, "pagination_limit", 10)`.
- Defaults from `config/settings.php`: `app_name`, `date_format`, `pagination_limit`.

## 12.8 Activity log

- Table name: `activity_log` (keep Laravel-compatible columns).
- Columns: `id`, `log_name`, `description`, `subject_type`, `subject_id`, `causer_type`, `causer_id`, `properties` (JSON), `event`, `batch_uuid`, `created_at`, `updated_at`.
- Enrich `properties.agent` with IP, user agent, URL, device, browser, and platform if the frontend or reports consume it.
- Log at least: login success/failure, impersonation start/stop, user CRUD, approval actions, and integration failures.
- Write asynchronously through a goroutine only if the request latency budget requires it; otherwise write synchronously for reliability.

## 12.9 Hash IDs

- `speps/go-hashids/v2`.
- Salt and length from env (`HASHIDS_SALT`, `HASHIDS_LENGTH=12`).
- Use only where the current API exposes hash IDs.
