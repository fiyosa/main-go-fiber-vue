# 18. Risks and Open Decisions

| Risk / decision                                          | Impact                                                                          | Recommendation                                                                                         |
| -------------------------------------------------------- | ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| Fiber v3 API changes                                     | Implementation churn                                                            | Pin the latest stable version and verify middleware signatures before coding                           |
| JWT revocation semantics                                 | Logged-out tokens could remain valid                                            | `token_version` + `auths.revoked_at`; impersonation rows checked on every request                      |
| 3-hour access token without frontend refresh             | Users re-login every 3 hours, same as today                                     | Keep parity; add frontend refresh later if desired                                                     |
| Login contract change                                    | SPA currently handles extra login branches (directory lookup and second factor) | Update the SPA login screen to local password only; remove the QR/OTP UI                               |
| 403 responses with `success: true`                       | Contract quirk                                                                  | Preserve exactly; add a comment in code                                                                |
| JSON delete changes `204` to `200` with `message`        | Frontend must accept a response body for delete                                 | Update contract tests; document the intentional change in Section 20                                    |
| Hard-coded controller messages                           | Mixed languages and duplicated text                                             | Move every message to `internal/lang/en.json` and `internal/lang/id.json`; CI checks key parity and non-empty messages     |
| SSO callback does not return a token                     | Flow is currently unusable from the SPA                                         | Decide before porting: disable, one-time code exchange, or cookie                                      |
| Impersonation authorization is commented out             | Any admin can impersonate any user                                              | Preserve parity for now; revisit with security owner                                                   |
| PDF visual parity                                        | Export differences                                                              | Prefer headless Chrome if layout matters                                                               |
| Excel import validation messages                         | User-visible errors                                                             | Port rule-by-rule and add contract tests                                                               |
| Multi-database contexts                                  | Complexity                                                                      | Merge owned data into core; keep external DBs as adapters only                                         |
| Legacy `voltProtect`/`voltAccess` and Laravel encryption | Not required per decision                                                       | Do not port; use new AES-GCM where encryption is needed                                                |
| Domain contract folder                                   | Empty placeholder with no real port                                             | Do not create `app/domain/<ctx>/contract`; keep contracts in the owning layer                              |
| Unused `SygapApiContract`                                | Dead code ported for no reason                                                  | Do not port; no implementation or caller exists                                                        |
| Integration client injection                             | Concrete client blocks test doubles                                             | Inject `contracts.Client` into resource clients                                                        |
| `sendError` 404 inference                                | Translated text must never decide the HTTP status                               | Do not port string-based inference; callers pass explicit status codes                                 |
| Timezone handling                                        | Date filters and attendance boundaries                                          | Store UTC; render `Asia/Jakarta`; test DST-free behavior                                               |
| Activity log volume                                      | Storage growth                                                                  | Retention policy + scheduled cleanup                                                                   |
| Rate limiting behind proxies                             | Incorrect client IP                                                             | Configure trusted proxies and use `X-Forwarded-For` safely                                             |
| Single-process runtime                                   | A job panic or memory spike can degrade the API                                 | Recover panics in job handlers, cap goroutine concurrency, set memory limits; keep split flags available |
| Duplicate schedules across replicas                      | Cron entries run more than once                                                 | PostgreSQL advisory lock; release on shutdown                                                          |
| Goroutine leak or blocking job                           | Process degradation or stuck shutdown                                           | Use context cancellation, timeouts, and bounded concurrency; test shutdown paths                       |

## Open decisions to confirm before Phase 1

1. Exact Fiber version to pin.
2. SSO callback token delivery.
3. Whether impersonation should remain unrestricted.
4. PDF renderer: headless Chrome or pure-Go layout.
5. Retention period for `activity_log` and `auths`.
6. Whether to keep the single-process model at high load or split the scheduler into a separate deployment.
