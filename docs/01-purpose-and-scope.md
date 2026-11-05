# 1. Purpose and Scope

This document is the implementation guide for rewriting the Laravel backend of Alma into a Go service that uses Fiber as the HTTP framework.

## In scope

- HTTP API used by the existing Vue 3 SPA under `/api/v1`.
- Public routes used by the SPA and external clients (`/media/*`, `/file-upload/*`, `/api/v1/*`, SSO callback).
- Authentication with local credentials, authorization, impersonation, and token lifecycle.
- LDAP directory synchronization; LDAP is not used for login authentication.
- Domain logic for all bounded contexts currently under `app/Domain`, `app/Application`, and `app/Presentasion`.
- Background jobs, scheduled commands, and job groups.
- An MCP server at `/api/mcp` for remote automation clients.
- Infrastructure adapters: database, cache, LDAP, storage, mail, Telegram, Pusher, Excel, PDF, and external HTTP integrations.

## Out of scope

- Rewriting the Vue frontend, AdminLTE styling, or Vite build pipeline.
- Migrating existing production data or preserving Laravel encryption formats.
- Renaming public API paths, database tables/columns, permission strings, or the "Alma" product name.
- Changing the business rules while porting. Behavior parity comes first; cleanup comes later.

## Compatibility rule

The frontend is the contract. Every request path, HTTP status, response envelope, field name, enum value, pagination shape, and file URL must stay compatible with the current frontend unless the document explicitly marks it as a change.

**Auth exception:** the login endpoint intentionally drops LDAP credential checks and second-factor verification. The SPA must submit local username/password and store the returned token directly. All other paths, payloads, and envelopes remain compatible.
