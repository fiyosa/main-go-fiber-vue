# 8. Authentication and Authorization

## 8.1 Goals

- Keep the client-visible behavior of Laravel Passport: Bearer tokens, expiry, and revocation.
- Use JWT access tokens for stateless validation, with server-side revocation.
- Use opaque rotating refresh tokens; store only their SHA-256 hash in the database.
- Keep the existing token response shape for the SPA.
- Login uses local username/password only; there is no LDAP credential bind and no second-factor challenge.
- LDAP is limited to directory synchronization.

Note: the current Laravel code uses Passport personal access tokens (`$user->createToken('userToken', ['*'])->accessToken`) with a 3-hour expiry and the `impersonation` scope. The Go design keeps the client-visible token behavior, but uses one `auths` table, a self-contained JWT access token, and an opaque rotating refresh token. Personal access tokens, OAuth clients, and scopes are not ported.

## 8.2 Token model

| Token         | Format                    | Lifetime            | Storage                            | Notes                        |
| ------------- | ------------------------- | ------------------- | ---------------------------------- | ---------------------------- |
| Access token  | JWT (HS256, shared `JWT_SECRET`) | `JWT_ACCESS_TTL`    | Raw in `auths.token`, unique index | Claims below                 |
| Refresh token | 256-bit random, base64url | `JWT_REFRESH_TTL`   | SHA-256 hash in `auths.refresh_token` | Rotated on every use      |

Both token values are created at login; the JWT and the refresh hash are stored in the same `auths` row. The access token is signed from the claims below; the refresh token is independent random bytes, not derived from the JWT.

Signing uses HS256 with one shared secret, read from `JWT_SECRET` (minimum recommended length: 32 characters). The issuer and audience are fixed constants (`app-api` / `app-spa`), not environment variables. There is no private/public key pair.

Access token claims:

```json
{
    "iss": "app-api",
    "aud": "app-spa",
    "sub": "123",
    "username": "jdoe",
    "roles": ["super admin"],
    "ver": 7,
    "jti": "01J...",
    "iat": 1770000000,
    "exp": 1770010800,
    "act": { "sub": "1", "username": "admin" }
}
```

- `ver` is the user's `token_version`. Incrementing it invalidates all existing access tokens for that user without a blacklist table.
- `jti` is a random claim that keeps every JWT unique; it is not a separate column.
- `act` is present only on impersonation tokens.
- There is no `scopes` claim; authorization uses roles and permissions.

## 8.3 Database table

Auth uses one table. There are no `oauth_*`, personal access token, or impersonation session tables.

```sql
CREATE TABLE auths (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token               TEXT NOT NULL,
    refresh_token       CHAR(64) NOT NULL,
    impersonator_id     BIGINT NULL REFERENCES users (id) ON DELETE SET NULL,
    meta                JSONB NULL,
    token_expires_at    TIMESTAMPTZ NOT NULL,
    refresh_expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at          TIMESTAMPTZ NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX auths_token_idx ON auths (token);
CREATE UNIQUE INDEX auths_refresh_token_idx ON auths (refresh_token);
CREATE INDEX auths_user_revoked_idx ON auths (user_id, revoked_at);
CREATE INDEX auths_refresh_expiry_idx ON auths (refresh_expires_at);
CREATE INDEX auths_impersonator_idx ON auths (impersonator_id)
    WHERE impersonator_id IS NOT NULL;
```

Column rules:

- `token` is the compact JWT string. The unique index allows a direct lookup for revocation and session management.
- `refresh_token` is the hex SHA-256 hash of the raw refresh token; the raw value is returned to the client once and never stored.
- `impersonator_id` is null for a normal session. When set, `user_id` is the impersonated user and `impersonator_id` is the admin.
- `meta` is optional JSON such as `ip`, `user_agent`, `device`, and `last_used_at`. A null value is treated as an empty object; no auth decision depends on it.
- `token_expires_at` is `now + JWT_ACCESS_TTL` at issue time and matches the JWT `exp` claim.
- `refresh_expires_at` is `now + JWT_REFRESH_TTL` at issue time and is used for cleanup and refresh validation.
- `revoked_at` marks a session as revoked without deleting the row.
- Changing `.env` TTL values affects only tokens issued after the change; existing rows keep their stored expiry.

Validation rules:

- Every request validates the JWT signature, `exp`, `iss`, `aud`, and `ver` against the user's `token_version`.
- Impersonation tokens (`act` claim) also require an `auths` row with matching `token` and `revoked_at IS NULL`, so stopping impersonation takes effect immediately.
- Refresh always requires an `auths` row with matching `refresh_token`, `revoked_at IS NULL`, and `refresh_expires_at > now()`.

Also add `token_version INTEGER NOT NULL DEFAULT 0` to `users`.

## 8.3.1 Initial auth data

Migrations create the schema; seeders fill it. The new database starts empty, so the seed step is part of the first install, not an optional convenience.

Seeded once (idempotent):

| Data                    | Source                                                                 |
| ----------------------- | ---------------------------------------------------------------------- |
| Roles and permissions   | Exported from the Laravel `roles`, `permissions`, and pivot tables      |
| Menus and menu access   | Exported from the Laravel menu tables and their permission mapping      |
| First admin user        | `SEED_ADMIN_USERNAME`, `SEED_ADMIN_EMAIL`, `SEED_ADMIN_PASSWORD`         |
| Settings and templates  | Defaults from the Laravel seeders                                       |

Rules:

- The admin password is hashed with bcrypt (cost 12) at seed time; it is never stored or logged in plaintext.
- Refuse to seed when `SEED_ADMIN_PASSWORD` is empty, shorter than 12 characters, or equal to a known default.
- `ON CONFLICT (...) DO NOTHING` for reference data and `DO UPDATE` only for fields that are safe to refresh.
- The admin seeder runs only when `users` is empty; it must not overwrite existing users.
- `token_version` starts at `0`; increment it on password change, forced logout, and role reset.

Tables that stay empty until runtime:

| Table   | Filled when                                              |
| ------- | -------------------------------------------------------- |
| `auths` | Login, refresh, or impersonation creates a session row    |

User population after install:

1. Seed the admin and reference data.
2. Run the LDAP directory sync to import employees, or create accounts through the admin API.
3. SSO-only accounts have `password = NULL`; local login rejects them until an admin sets a password.
4. Invalidate sessions by bumping `token_version` when an account is re-created or its password is reset.

Deployment sequence for a fresh database:

```bash
console migrate status
console migrate
console db seed                                  # all seeders
console db seed ReferenceDataSeeder              # roles, permissions, menus, settings
```

The `app` process never seeds or migrates. The seed command reads secrets from `.env` only; no credentials are committed to the repository.

## 8.3.2 Authorization tables

RBAC uses four tables with no guard or model type, because the application has one authentication guard.

```sql
CREATE TABLE roles (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(125) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE permissions (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(125) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_has_roles (
    user_id  BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id  BIGINT NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX user_has_roles_role_idx ON user_has_roles (role_id);

CREATE TABLE role_has_permissions (
    role_id        BIGINT NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id  BIGINT NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX role_has_permissions_permission_idx ON role_has_permissions (permission_id);
```

Rules:

- `roles.name` and `permissions.name` are unique; permission strings stay unchanged from Laravel.
- A user can have many roles; permissions are the union of all roles.
- There are no direct user permissions and no pivot timestamps.
- `super admin` bypasses every permission check in the `Authorizer`; it does not need every permission assigned.
- Deleting a user, role, or permission cascades to the pivot rows.

Effective permission query:

```sql
SELECT DISTINCT p.name
FROM user_has_roles ur
JOIN role_has_permissions rp ON rp.role_id = ur.role_id
JOIN permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1;
```

## 8.4 Endpoints

| Method | Path                  | Purpose                                         | Compatibility                                    |
| ------ | --------------------- | ----------------------------------------------- | ------------------------------------------------ |
| `POST` | `/api/v1/login`       | Local username/password login                   | Returns tokens directly; no second-factor branch |
| `POST` | `/api/v1/refresh`     | Rotate refresh token and issue new access token | New endpoint                                     |
| `POST` | `/api/v1/logout`      | Revoke all user sessions                        | Same                                             |
| `GET`  | `/api/v1/me`          | Current user + roles + permissions              | Same                                             |

Success messages use `lang.logged_in_successfully` and `lang.logged_out_successfully`; failures use `lang.invalid_credentials` and `auth.unauthenticated`. See [Section 20](20-localization.md).

Login response (success):

```json
{
    "success": true,
    "message": "Logged in successfully",
    "token": "<jwt>",
    "token_type": "Bearer",
    "expires_in": 10800,
    "refresh_token": "<opaque>"
}
```

The SPA reads `response.data.token` or `response.token`, so both top-level and nested shapes are safe. Keep `token` top-level to match `authenticatedResponse`. `expires_in` is `JWT_ACCESS_TTL` in seconds.

Invalid credentials:

```json
{ "success": false, "message": "Username or Password is invalid" }
```

Status `401`.

## 8.5 Login flow

```text
POST /api/v1/login
  ├─ validate username/password
  ├─ load local user by username or email
  │    └─ missing → 401 invalid_credentials
  ├─ verify users.password hash (bcrypt)
  │    └─ fail → 401 invalid_credentials
  ├─ issue JWT access token and random refresh token
  ├─ insert auths row: token, refresh hash, expiries, meta
  └─ update last_login_at
```

Token generation:

```go
func newRefreshToken() (raw, hash string, err error) {
    buf := make([]byte, 32)
    if _, err = rand.Read(buf); err != nil {
        return "", "", err
    }
    raw = base64.RawURLEncoding.EncodeToString(buf)
    sum := sha256.Sum256([]byte(raw))
    return raw, hex.EncodeToString(sum[:]), nil
}
```

- `token_expires_at = now + JWT_ACCESS_TTL`; `refresh_expires_at = now + JWT_REFRESH_TTL`.
- Only the hash is stored in `auths.refresh_token`; the raw value is returned to the client.
- On refresh, hash the incoming token, find the row, issue a new JWT and a new refresh token, and update the same row.

No LDAP bind. No second-factor challenge. No trusted-device cookie.

## 8.6 LDAP directory integration

LDAP is not used for login authentication. It remains for directory synchronization and employee lookup.

- Use `go-ldap/ldap/v3`.
- Support the existing server configs: `openldap`, `domino` (and any future servers).
- Bind with the service account for search and sync jobs only; never validate user credentials at login.
- Upsert local users from sync jobs:
    - `username` from `uid` or configured attribute.
    - `email` from `mail`.
    - `name` from `cn` or `displayName`.
- Local `users.password` (bcrypt) is the only login credential. Do not store or read LDAP passwords.
- Map LDAP errors to `503` for sync failures; login must never depend on LDAP availability.

## 8.7 Authorization

- Use the tables `roles`, `permissions`, `user_has_roles`, and `role_has_permissions` from [Section 8.3.2](#832-authorization-tables).
- Keep permission strings unchanged.
- `super admin` bypasses every permission check, matching `Gate::before`.
- Policies from `app/Domain/Alma/Policies` map to an `Authorizer` service:
    - `DocumentPolicy`
    - `MediaPolicy`
    - `TaskPolicy`
    - `TaskCardPolicy`
    - `UserPolicy`
- Middleware:
    - `RequirePermission("manage_roles")`
    - `RequireRole("super admin")`
    - `RequireImpersonator()` for impersonation stop
- `GET /me` returns the same `AuthUserResource` fields:
    - `id`, `name`, `username`, `email`, `role`, `roles`, `permissions`, `created_at`, `formatted_created_at`, `email_verified_at`, `photo`, `verified`, `division`.

## 8.8 Impersonation

Current behavior:

- Any authenticated admin can impersonate any other user except themselves.
- Super-admin restrictions are commented out in the Laravel controller.
- Start returns a new token with scope `impersonation`.
- Stop revokes the impersonation token and returns the original token.

Go implementation:

- `POST /api/v1/impersonation/{user}`:
    - Reject self-impersonation with `422`.
    - Create a new `auths` row: `user_id` is the target user, `impersonator_id` is the admin, `token` is a JWT for the target user with the `act` claim, and `refresh_token` is a new hashed refresh token.
    - Return `{token, refresh_token, token_type, expires_in, impersonating: true, impersonator_id, user}`.
- `POST /api/v1/impersonation/stop`:
    - Require a token with the `act` claim and a matching `auths` row where `impersonator_id` is not null.
    - Set `revoked_at` on the impersonation row.
    - Create a new normal `auths` row for `impersonator_id` and return its tokens.
    - The admin's original session is independent and remains valid; no raw JWT is stored for restoration.

No `impersonation` scope and no `impersonation_sessions` table are used. Revocation is immediate because impersonation requests require the `auths` row to be active.

## 8.9 SSO callback

Current `AuthSSO` uses a session login and redirects to `/admin/dashboard` without returning a token. The SPA stores tokens in `localStorage`, so this flow does not currently produce a usable SPA session.

Go decision required:

- Option A: disable the legacy SSO route during the rewrite.
- Option B: complete the flow by redirecting to the SPA with a one-time code that the SPA exchanges for tokens.
- Option C: set the access token in an HttpOnly cookie and adapt the frontend.

Document the chosen option before implementing. Do not silently keep a broken flow.

## 8.10 Throttling

| Limiter          | Key               | Limit     | Current   |
| ---------------- | ----------------- | --------- | --------- |
| `login`          | `lower(username)  | ip`       | 5/minute  | Yes |
| `api`            | `user_id` or `ip` | 60/minute | Yes       |
| Export endpoints | user              | 1/minute  | Per route |
| Temp sync        | IP/user           | 2/minute  | Yes       |

Use a PostgreSQL-backed Fiber limiter (atomic upsert into `rate_limits`) so limits are shared across API replicas.

## 8.11 Logout and revocation

Current Laravel logout deletes all tokens for the user. Preserve this behavior:

1. Increment `users.token_version`.
2. Set `auths.revoked_at = now()` for every active row where `user_id = user.id` or `impersonator_id = user.id`.
3. Return `{success: true, message: "Logged out successfully"}` from `lang.logged_out_successfully`.

Because `ver` is checked on every request, existing JWTs become invalid immediately without a growing blacklist. Per-session revocation of an impersonation token uses its `auths` row.

## 8.12 Frontend compatibility notes

- `plugins/adminApi.js` sends `Authorization: Bearer <token>`.
- On `401`, the frontend dispatches `auth:unauthorized` and clears storage.
- The access token TTL comes from `JWT_ACCESS_TTL`; keep the default at 3 hours so the current re-login behavior is unchanged.
- A `/refresh` endpoint is available, but the current frontend does not call it. Adding a refresh interceptor is a separate frontend task.
- The login screen submits username/password and stores the returned token directly; it must not branch on a second-factor response or render a QR/OTP form.
