# 2. Locked Decisions

| Area              | Decision                                                                                                    |
| ----------------- | ----------------------------------------------------------------------------------------------------------- |
| Rewrite scope     | Backend only; frontend Vue stays except the login screen drops LDAP and second-factor handling              |
| HTTP framework    | Fiber, latest stable version at implementation time                                                         |
| ORM               | GORM for CRUD and relations; raw SQL for complex queries                                                    |
| Jobs              | Goroutines spawned by use cases; no queue or broker                                                         |
| MCP               | Streamable HTTP at `/api/mcp` on the API port; bearer token from `.env`; read-only tools first             |
| Scheduler         | robfig/cron v3                                                                                              |
| Process model     | One process runs API and cron scheduler; jobs are goroutines; console commands are one-off                  |
| Auth tokens       | JWT access token + opaque rotating refresh token, stored in one `auths` table                              |
| Auth behavior     | Passport-compatible expiry and revocation; no scopes, OAuth clients, or personal access tokens              |
| RBAC              | Four user-only tables: `roles`, `permissions`, `user_has_roles`, `role_has_permissions`; no guard/model type |
| Database          | New schema; no legacy data or encryption compatibility                                                      |
| Document language | English                                                                                                     |
| API messages      | Every JSON response has a translated `message`; default `en`, `id` supported; catalogue in `lang/*.json`     |
| Domain rename     | `alma` → `core` in Go folders, packages, and identifiers                                                    |
| API compatibility | Paths, payloads, status codes, and response envelopes preserved except the login and second-factor removals |
