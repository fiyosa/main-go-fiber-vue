# 16. Testing Strategy

## 16.1 Test pyramid

| Level       | Scope                                               | Tools                           |
| ----------- | --------------------------------------------------- | ------------------------------- |
| Unit        | Enums, policies, mappers, validators, token service | `testing`, `testify`            |
| Integration | Handlers + DB                                       | `testcontainers-go`, `httptest` |
| Contract    | JSON shape and status codes vs Laravel responses    | Golden files                    |
| End-to-end  | Login → API → logout                                | Testcontainers + Fiber app      |
| Load        | Login, attendance check-in, list endpoints          | `k6` or `vegeta`                |

Integration and end-to-end tests use a PostgreSQL `testcontainers-go` container that matches the production major version; do not substitute another database engine.

## 16.2 Mandatory tests

- Login success, invalid credentials.
- JWT validation: signature, expiry, issuer, audience, `token_version`.
- Refresh rotation invalidates the previous hash; revoked or expired rows are rejected.
- Logout invalidates all sessions.
- Impersonation start/stop and `impersonator_id` enforcement.
- Permission middleware: allowed, denied, super-admin bypass.
- RBAC: multi-role permission union, cascade delete of pivot rows.
- Route permission audit: every route matches [Appendix D](appendices/d-route-permissions.md); non-public routes declare an access category; permission routes return `403` without permission and `2xx` with permission.
- PIN middleware: wrong PIN attempts and lockout.
- Response envelope parity for success, error, validation error, and pagination.
- Localization: `en.json` and `id.json` have identical keys; missing key returns the key and logs a warning; unsupported locale resolves to `en`; every JSON response has a non-empty `message`.
- Default locale is `en`; `Accept-Language: id` and `?lang=id` return Indonesian messages.
- MCP bearer token from `.env`: missing, malformed, wrong token, missing actor where required.
- MCP tool `whoami` and `search_users` return the expected structured content; an `mcp_audit_logs` row is written; rate limit is enforced.
- Every route in the inventory has at least a smoke test.

## 16.3 Route parity harness

1. Export the Laravel route list in a working environment (`php artisan route:list --json`).
2. Export the Go route table from a test-only registry.
3. Compare method + path + middleware names.
4. Compare each route with the access category and permission in [Appendix D](appendices/d-route-permissions.md).
5. Fail CI on missing or extra routes unless explicitly allowlisted.

## 16.4 Contract tests

Capture representative Laravel responses for each module (list, detail, create, validation error, unauthorized). Store them as golden JSON. The Go handler must produce the same JSON keys and value types.
