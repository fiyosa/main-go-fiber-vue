# Appendix D: Route Permission Audit

## D.1 Purpose

This appendix is the source of truth for checking every route registered by `internal/app/presentation/{domain}/routes/*.go`. When an API is hit, the route must have the same access decision as Laravel:

- public route: no token required;
- authenticated route: valid JWT required, no permission required;
- permission route: valid JWT plus the named permission required;
- super admin: bypasses permission checks through `Authorizer`;
- impersonation route: valid impersonation token plus an active `auths` row.

The matrix below is derived from the current Laravel route files and FormRequest `authorize()` methods. Laravel behavior wins unless a difference is explicitly recorded in [Section 18](../18-risks-and-open-decisions.md).

Do not invent permission strings during the Go rewrite. If a route has no Laravel permission, keep it authenticated-only until a product decision assigns one.

## D.2 Access categories

| Category | Go middleware | Expected anonymous response | Expected authenticated response |
| --- | --- | --- | --- |
| Public | none | `2xx` | `2xx` |
| Authenticated | inherited `auth:admin` / JWT check | `401` | `2xx` unless the controller rejects the request |
| Permission | `deps.Auth.Permission("<name>")` | `401` | `403` without permission, `2xx` with permission |
| Super admin | `deps.Auth.Permission(...)` plus `Authorizer` bypass | `401` | `2xx` for `super admin` |
| Impersonation | JWT `act` claim plus active `auths` row | `401` | `2xx` only while the session is active |

`auth:admin` is applied once at the API group. A module route file must not add a second global auth middleware.

## D.3 Confirmed Laravel route-level permissions

Only the Alma route files currently contain explicit route middleware. Other contexts are authenticated-only at route level.

| Go route file | Prefix | Protected endpoint | Permission | Laravel source |
| --- | --- | --- | --- | --- |
| `attendance_routes.go` | `/attendances` | `GET /shift-options` | `attendances.can-edit` | route middleware |
| `attendance_routes.go` | `/attendances` | `GET /export-template` | `attendances.download-template` | route middleware |
| `attendance_routes.go` | `/attendances` | `POST /import` | `attendances.import-attendance` | route middleware |
| `attendance_routes.go` | `/attendances` | `POST /` | `attendances.can-edit` | route middleware |
| `attendance_routes.go` | `/attendances` | `DELETE /{attendance}` | `attendances.can-edit` | route middleware |
| `attendance_routes.go` | `/attendances` | `PUT /{attendance}/toggle-auto-wip` | `attendances.toggle-auto-wip` | route middleware |
| `aux_type_routes.go` | `/manage-aux-types` | all CRUD endpoints | `manage_aux_types` | group middleware |
| `email_template_routes.go` | `/manage-email-templates` | all CRUD endpoints | `manage_email_templates` | group middleware |
| `form_routes.go` | `/manage-forms` | all CRUD endpoints | `manage_forms` | group middleware |
| `holiday_routes.go` | `/holidays` | all CRUD endpoints | `manage_holidays` | group middleware |
| `menu_routes.go` | `/menus/manage` | all management endpoints | `manage_menus` | group middleware |
| `role_routes.go` | `/roles` | all management endpoints | `manage_roles` | group middleware |
| `shift_routes.go` | `/manage-shifts` | shifts, patterns, and options | `manage_shifts` | group middleware |

Unprotected routes in the same files remain authenticated-only unless listed elsewhere in this appendix:

- `GET /attendances`, `/attendances/me`, check-in/check-out, aux, user options, date options, detail, export, update, and reset;
- `GET /menus`;
- `GET /roles/all`;
- `PUT /shifts/{id}/toggle-automation`;
- `GET /shift-schedules/shift-today`;
- `GET /shift-baselines/organization-option`;
- `GET /forms/form-options` and `GET /forms/{id}`.

## D.4 Permission enforced only by FormRequest

Laravel also enforces permissions inside FormRequest `authorize()` methods. These checks duplicate the route middleware above and must not be the only enforcement in Go.

| Permission | Request classes | Endpoints |
| --- | --- | --- |
| `attendances.can-edit` | `AttendanceStoreRequest`, `AttendanceShiftOptionsRequest` | `POST /attendances`, `GET /attendances/shift-options` |
| `manage_aux_types` | `AuxTypeStoreRequest`, `AuxTypeUpdateRequest` | `POST` and `PUT /manage-aux-types` |
| `manage_email_templates` | `EmailTemplateStoreRequest`, `EmailTemplateUpdateRequest` | `POST` and `PUT /manage-email-templates` |
| `manage_forms` | `FormStoreRequest`, `FormUpdateRequest` | `POST` and `PUT /manage-forms` |
| `manage_holidays` | `HolidayStoreRequest`, `HolidayUpdateRequest` | `POST` and `PUT /holidays` |
| `manage_menus` | `MenuStoreRequest`, `MenuUpdateRequest` | `POST` and `PUT /menus/manage` |
| `manage_roles` | `RoleStoreRequest`, `RoleUpdateRequest` | `POST` and `PUT /roles` |
| `manage_shifts` | `ShiftStoreRequest`, `ShiftUpdateRequest`, `ShiftPatternStoreRequest`, `ShiftPatternUpdateRequest` | shifts and shift-pattern writes |

In Go:

- enforce the permission at the route or route group;
- use the request/controller only for ownership, status, or business-rule checks;
- never rely on a request validator to produce `403` for a permission failure.

## D.5 Route files with no route-level permission

These files are authenticated-only in Laravel. They are the audit backlog: every route must be explicitly classified before it is ported.

| Go route file | Sensitive endpoints to review | Current Laravel state |
| --- | --- | --- |
| `user_routes.go` | user CRUD, role and division changes, bulk delete | authenticated only |
| `appointment_routes.go` | create, update, delete, client lookup | authenticated only |
| `broadcast_routes.go` | broadcast CRUD, send mail, export | authenticated only |
| `division_routes.go` | CRUD | authenticated only |
| `document_routes.go` | upload, delete, download | authenticated only |
| `form_request_routes.go` | approve, reject, delete, export history | authenticated only |
| `hrdwh_routes.go` | employee, mutation, report, personal workbook | authenticated only |
| `legacy_domino_letter_routes.go` | list, export | authenticated only |
| `mailadmin_routes.go` | profile read and update | authenticated only |
| `profile_routes.go` | own profile update, image upload, password change | authenticated only; ownership check required |
| `task_routes.go` | task, list, card, upload, download | authenticated only |
| `intra_routes.go` | application comments and export | authenticated only |
| `kampiun_routes.go` | access-document report and export | authenticated only |
| `kypas_routes.go` | approvals, employees, contractors, imports, changes | authenticated only; two lookups use `verify.pin` |
| `appstoken_routes.go` | token CRUD and email | authenticated only |
| `rfc_role_routes.go` | role list | authenticated only |
| `rfc_user_routes.go` | user list and update | authenticated only |
| `sygap_routes.go` | incidents, requests, tasks, mass assignment | authenticated only |
| `sygap_employee_routes.go` | employee list | authenticated only |
| `tgcast_routes.go` | mail list and export | authenticated only |
| `tgkypas_routes.go` | access requests, contract changes, sync, exports | authenticated only |
| `windigital_routes.go` | document list and download | authenticated only |

Special cases:

- `verify.header` on `POST /kypas/file-upload` and `verify.pin` on Kypas `showmore` endpoints are not RBAC permissions. Keep them as dedicated middleware.
- `GET /api/v1/kypas/sync-password/{username}` stays public with its own throttle.
- `GET /form-requests/{id}/generate-pdf` stays public.
- `ldap.php` is not ported; do not add `ldap_routes.go`.

If a permission is added to one of these routes during the rewrite, record the difference in [Section 18](../18-risks-and-open-decisions.md).

## D.6 Per-route audit checklist

For every route in a Go route file:

1. Confirm the method and relative path match the Laravel route and the inventory in [Appendix A](a-route-inventory.md).
2. Confirm the route inherits `auth:admin`; do not repeat the global auth middleware.
3. If the route is public, confirm it appears in the public allowlist.
4. If the route is authenticated-only, confirm there is no accidental `Permission(...)` wrapper.
5. If the route requires a permission, confirm `deps.Auth.Permission("<exact string>")` is attached to the route or its group.
6. Confirm the permission string exactly matches the matrix. Typos are authorization bugs.
7. Confirm the route is not protected by role-only checks where Laravel uses a permission.
8. Confirm `super admin` bypass works through `Authorizer`, not through a separate route callback.
9. Confirm impersonation routes also require an active `auths` row.
10. Confirm the route file declares no absolute `/api/v1` prefix.
11. Confirm removed routes are not registered again.
12. Confirm the route has at least a smoke test and one authorization test.

## D.7 Runtime API verification

Run the following personas against every matrix row:

| Persona | Token | Expected for authenticated-only | Expected for permission route |
| --- | --- | --- | --- |
| Anonymous | none | `401` | `401` |
| No permission | valid JWT, role without permission | `2xx` | `403` |
| With permission | valid JWT, role with permission | `2xx` | `2xx` |
| Super admin | valid JWT | `2xx` | `2xx` |
| Revoked session | valid-looking JWT, `auths` row revoked | `401` | `401` |

Record the actual status for each method and path. A `500` is a failed audit, even when the expected result is a denial.

## D.8 Automated tests

The Go test suite must enforce this appendix:

1. **Route registry test**: enumerate all registered routes and fail when a non-public route has no explicit access category.
2. **Matrix test**: compare every registered method and path with this appendix; fail on missing entries and on extra routes unless allowlisted.
3. **Permission test**: for each permission route, assert `403` without permission and `2xx` with permission.
4. **Super-admin test**: assert the bypass on every permission route.
5. **Public-route test**: assert no token is required and that no accidental auth middleware is attached.
6. **Route parity test**: compare the Go route table with `php artisan route:list --json`; middleware names may differ, but the access decision may not.

Suggested locations:

```text
tests/integration/route_permissions_test.go
tests/integration/route_registry_test.go
tests/contract/route_parity_test.go
```

Keeping the matrix current is part of adding any new route. A new route without a matrix entry is a failing test, not a documentation task for later.
