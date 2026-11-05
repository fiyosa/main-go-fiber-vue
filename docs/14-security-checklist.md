# 14. Security Checklist

- [ ] Secrets only from environment variables; never commit them.
- [ ] JWT signing keys rotated with `kid` support.
- [ ] Access tokens short-lived (`JWT_ACCESS_TTL`); refresh tokens rotated and stored as SHA-256 hashes.
- [ ] Refresh token rotation invalidates the previous hash; row revocation takes effect on refresh.
- [ ] `auths.revoked_at` invalidates all sessions on logout.
- [ ] Password policy: minimum 12 characters, mixed case, number, symbol, and breached-password check (HIBP) if the policy remains.
- [ ] Login rate limit 5/minute.
- [ ] API rate limit 60/minute per user/IP.
- [ ] CORS allowlist; `supports_credentials` only if cookies are used cross-origin.
- [ ] CSRF: not required for Bearer-token APIs; required if cookies become the auth transport.
- [ ] SQL injection prevented by parameterized GORM queries.
- [ ] File upload validation: extension, MIME, size, and storage path traversal protection.
- [ ] Trusted proxy configuration for `X-Forwarded-For`.
- [ ] Security headers: `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, HSTS in production.
- [ ] Audit log for all auth and authorization events.
- [ ] No secrets or tokens in logs.
- [ ] Dependency scanning and `govulncheck` in CI.
