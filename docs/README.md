# Laravel → Go + Fiber Rewrite Guide

> Status: Draft
> Scope: Backend rewrite only. The existing Vue 3 frontend must keep working without changes.
> Target stack: Go + Fiber (latest stable) + GORM + robfig/cron + official MCP SDK.
> Authentication target: JWT access token + rotated hashed refresh token in one `auths` table, behavior-compatible with Laravel Passport expiry and revocation.
> Database: new schema, no legacy data migration required.
> Naming rule: every `Alma` / `alma` folder, package, and identifier becomes `Core` / `core`.

## Documentation Index

| Section | File |
| --- | --- |
| 1. Purpose and Scope | [01-purpose-and-scope.md](01-purpose-and-scope.md) |
| 2. Locked Decisions | [02-locked-decisions.md](02-locked-decisions.md) |
| 3. Current Laravel Architecture | [03-current-laravel-architecture.md](03-current-laravel-architecture.md) |
| 4. Target Go Architecture | [04-target-go-architecture.md](04-target-go-architecture.md) |
| 5. Repository Layout and Naming | [05-repository-layout.md](05-repository-layout.md) |
| 6. Layer Mapping: Laravel → Go | [06-layer-mapping.md](06-layer-mapping.md) |
| 7. HTTP API Contract | [07-http-api-contract.md](07-http-api-contract.md) |
| 8. Authentication and Authorization | [08-authentication.md](08-authentication.md) |
| 9. Data Layer | [09-data-layer.md](09-data-layer.md) |
| 10. Background Processing | [10-background-processing.md](10-background-processing.md) |
| 11. External Integrations | [11-external-integrations.md](11-external-integrations.md) |
| 11.4 Casdoor reference | [11b-casdoor-reference.md](11b-casdoor-reference.md) |
| 12. Platform Services | [12-platform-services.md](12-platform-services.md) |
| 13. Frontend Delivery | [13-frontend-delivery.md](13-frontend-delivery.md) |
| 14. Security Checklist | [14-security-checklist.md](14-security-checklist.md) |
| 15. Observability and Operations | [15-observability.md](15-observability.md) |
| 16. Testing Strategy | [16-testing-strategy.md](16-testing-strategy.md) |
| 17. Migration Roadmap | [17-migration-roadmap.md](17-migration-roadmap.md) |
| 18. Risks and Open Decisions | [18-risks-and-open-decisions.md](18-risks-and-open-decisions.md) |
| 19. MCP Server | [19-mcp-server.md](19-mcp-server.md) |
| 20. Localization | [20-localization.md](20-localization.md) · [en](lang/en.json) · [id](lang/id.json) |
| Appendix A: Route Inventory | [appendices/a-route-inventory.md](appendices/a-route-inventory.md) |
| Appendix B: Environment Variable Names | [appendices/b-configuration.md](appendices/b-configuration.md) |
| Appendix C: Go Package List | [appendices/c-go-packages.md](appendices/c-go-packages.md) |
| Appendix D: Route Permission Audit | [appendices/d-route-permissions.md](appendices/d-route-permissions.md) |

## Final Notes

- Port behavior before improving design. Every intentional difference must be listed in [Section 18](18-risks-and-open-decisions.md).
- The fastest path to confidence is contract tests against Laravel responses plus route parity checks.
- The highest-risk areas are authentication, attendance date logic, Excel/PDF exports, and external integrations. Assign senior engineers to those modules.
- Keep the Laravel app running until the Go service has survived at least one full business cycle (monthly attendance and payroll-adjacent reports).
