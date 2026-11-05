# 15. Observability and Operations

## 15.1 Logging

- `slog` JSON in production, text in development.
- Fields: `request_id`, `user_id`, `route`, `method`, `status`, `latency_ms`, `job`, `group`, `duration_ms`.
- MCP tool calls log `mcp_user_id`, `tool`, `status`, `duration_ms`, and `args_hash`; never log the token or raw arguments.
- Add a request ID middleware and return it in `X-Request-ID`.

## 15.2 Health and metrics

- `GET /healthz`: process is alive and the API and scheduler goroutines are running.
- `GET /readyz`: DB is reachable and the scheduler has started.
- Prometheus metrics: HTTP duration, status codes, in-flight job goroutines, job failures, token refresh failures, integration latency, `mcp_tool_calls_total`, `mcp_tool_duration_seconds`, `mcp_auth_failures_total`.
- Protect `/metrics` and `/debug/pprof` from public access.

## 15.3 Graceful shutdown

- Single process: on `SIGTERM` or `SIGINT`, cancel the shared context.
- Stop the scheduler first so it cannot start new jobs.
- Stop the HTTP server from accepting new requests and wait for in-flight requests.
- Wait for in-flight job goroutines with `jobs.Wait` and a bounded timeout.
- Close DB after `errgroup.Wait()` returns.
- Bound the whole shutdown with a timeout; force-exit if it expires.

## 15.4 Deployment

- Multi-stage Docker build; ship one `app` binary and one `console` binary.
- Inject database credentials from `.env` or the orchestrator's secret store at runtime; never bake them into the image.
- One deployment runs API + scheduler in the same process; job goroutines share that process.
- Run `console migrate` as an init job before the app replicas start; never migrate inside the `app` process.
- Fail the init job when `console migrate status` reports dirty; seed only on first install with `console db seed`.
- One readiness probe covers the whole process.
- Horizontal scaling: scale the app replicas; each replica serves HTTP and spawns its own job goroutines, while the scheduler uses a PostgreSQL advisory lock so only one replica runs cron.
- Use `--disable-scheduler` only for debugging; do not use it to create separate production deployments unless the single-process model is revisited.
