# 17. Migration Roadmap

Because the database is new and no legacy data must be preserved, use a phased build with a single cutover rather than dual-write. Keep the Laravel deployment available for rollback.

## Phase 0 — Foundation

- Repository scaffold, config, logging, health checks.
- Single-process runtime wiring: API and cron scheduler under `errgroup`, with graceful shutdown and job goroutines drained by `jobs.Wait`.
- PostgreSQL + `console migrate` and `console db seed`.
- Response envelope, validation error mapper, request ID, pagination.
- CI: lint (`golangci-lint`), test, build, `govulncheck`.

Exit criteria: `GET /healthz` and `GET /readyz` pass; one sample module with contract tests.

## Phase 1 — Auth and RBAC

- Users, roles, permissions, menus, settings, profile.
- JWT, refresh tokens, impersonation.
- Activity log for auth events.
- Login and `/me` parity tests.

Exit criteria: the existing SPA can log in with local credentials, load menus and permissions, and log out against the Go API.

## Phase 2 — Core read modules

- Dashboard, appointments, divisions, documents, forms, holidays, shifts, tasks, broadcasts, HRDWH read endpoints.
- Criteria filters and pagination parity.

Exit criteria: all read-only screens work against Go.

## Phase 3 — Core write modules

- Attendance check-in/out, aux, imports/exports.
- Task board mutations, forms, approvals, documents/media, notifications.
- Excel/PDF exports.

Exit criteria: end-to-end operational flows pass; exports match templates.

## Phase 4 — External contexts

Order by business criticality and dependency depth:

1. Kypas
2. Sygap
3. Tgkypas
4. Windigital
5. Intra
6. Kampiun
7. Rfc
8. Tgcast
9. Mailadmin
10. Legacy contexts (Diarium, Evolution, NewCuti, Lemon, Tgsso, Action, Universal, Insera, MicrosoftGraph)

Exit criteria: each context has route parity, contract tests, and its integrations working.

## Phase 5 — Async and integrations

- Port all jobs, schedules, and console commands.
- Verify job goroutines and retries.
- Verify the scheduler goroutine and job goroutines coexist with the API under load.
- Run jobs in parallel with Laravel in a staging environment for comparison.

Exit criteria: all scheduled jobs run in Go; no missing job functions; monitoring in place.

## Phase 6 — MCP server

- Add the official Go SDK and mount Streamable HTTP at `/api/mcp` on the existing Fiber port.
- Bearer token middleware that compares `MCP_BEARER_TOKEN` from `.env`, plus `mcp_audit_logs` and per-token rate limiting.
- Start with `ping`, `whoami`, and `search_users`; add domain tools as modules complete.
- Contract tests use the SDK client against the Fiber test server.

Exit criteria: a remote MCP client connects over HTTPS with a bearer token and calls the read-only tools.

## Phase 7 — Cutover

- Deploy the single Go process (API + scheduler).
- Point the reverse proxy `/api`, `/media`, `/file-upload`, and SPA fallback to Go.
- Keep Laravel deployed but idle for rollback.
- Validate login, attendance, approvals, exports, and jobs.
- Monitor error rates and latency for 1–2 weeks.
- Decommission Laravel after sign-off.

Rollback: switch the reverse proxy back to Laravel. The new database is separate, so no data conflict exists.
