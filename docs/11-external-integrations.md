# 11. External Integrations

## 11.1 Pattern

Every integration package contains:

```text
app/infrastructure/integrations/<name>/
├── client.go          # base HTTP client
├── clients/
│   └── <resource>_client.go
├── contracts/
│   └── client.go      # adapter-local client interface
├── dto/
│   ├── requests/
│   └── responses/
├── exceptions/
│   └── <name>_exception.go
├── support/
│   └── endpoint.go
└── token_manager.go   # only when the upstream requires a token
```

This mirrors `app/Infrastructure/Integrations/<Name>`; do not collapse the folders into a flat package.

Rules:

- Base client injects the token, timeouts, and user agent.
- Token manager, when present, caches in the PostgreSQL `cache` table with a 5-minute safety buffer.
- Use `singleflight` to prevent token stampedes.
- Map non-2xx responses to typed errors that preserve status and response body.
- Never log credentials or raw tokens.
- Casdoor is the reference implementation; see [Section 11.4](11b-casdoor-reference.md#114-reference-implementation-casdoor) for the full layout, file mapping, and Go examples.

## 11.2 Inventory

| Integration     | Clients today                                                                | Go package                    |
| --------------- | ---------------------------------------------------------------------------- | ----------------------------- |
| Casdoor         | `UserClient`                                                                 | `integrations/casdoor`        |
| Domino          | `EmployeeClient`                                                             | `integrations/domino`         |
| Insera          | `TicketClient`                                                               | `integrations/insera`         |
| Kypas           | `UserClient`                                                                 | `integrations/kypas`          |
| Mailadmin       | `UserClient`                                                                 | `integrations/mailadmin`      |
| Microsoft Graph | `UserClient`                                                                 | `integrations/microsoftgraph` |
| Milistadmin     | `MemberClient`                                                               | `integrations/milistadmin`    |
| Nde             | `EmployeeClient`                                                             | `integrations/nde`            |
| OnePortal       | `UserClient`                                                                 | `integrations/oneportal`      |
| Sygap           | `EmployeeClient`, `IncidentClient`, `RequestFulfillmentClient`, `TaskClient` | `integrations/sygap`          |
| Tgkypas         | `ContractorClient`, `FulfillmentClient`, `UserClient`                        | `integrations/tgkypas`        |
| Universal       | `InboxClient`                                                                | `integrations/universal`      |
| Windigital      | `DocumentClient`                                                             | `integrations/windigital`     |

Casdoor is the reference implementation. Every integration follows the same layout; add `token_manager.go` only when the upstream requires a bearer-token exchange.

## 11.3 Resilience

- Connect timeout: 5 seconds. Total timeout: 15 seconds, unless the operation is a file transfer.
- Retry only idempotent operations.
- Circuit breaker optional per integration (`sony/gobreaker`).
- Emit metrics: request count, latency, error rate, token refresh count.
- Log request ID from the incoming context.
