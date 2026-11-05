# Appendix A: Route Inventory

## A.1 Public routes

| Method | Path                                  | Purpose                        |
| ------ | ------------------------------------- | ------------------------------ |
| `GET`  | `/admin/login`                        | SPA shell                      |
| `GET`  | `/media/{media:file_name}`            | Public media                   |
| `GET`  | `/file-upload/{fileUpload:file_name}` | Public uploaded file           |
| `GET`  | `/form-requests/{id}/generate-pdf`    | Public PDF                     |
| `GET`  | `/auth/login/oauth2`                  | SSO redirect                   |
| `GET`  | `/callback`                           | SSO callback                   |
| `GET`  | `/api/v1/kypas/sync-password/{username}` | Temporary sync endpoint, 2/min |
| `GET`  | `/*`                                  | SPA fallback (`internal/routes/web.go`)  |

## A.2 Core routes

See [Section 7.2](../07-http-api-contract.md). One Go route file per Laravel file:

| Laravel route file       | Go route file                    |
| ------------------------ | -------------------------------- |
| `appointment.php`        | `appointment_routes.go`          |
| `attendance.php`         | `attendance_routes.go`           |
| `auxType.php`            | `aux_type_routes.go`             |
| `broadcast.php`          | `broadcast_routes.go`            |
| `division.php`           | `division_routes.go`             |
| `document.php`           | `document_routes.go`             |
| `emailTemplate.php`      | `email_template_routes.go`       |
| `form.php`               | `form_routes.go`                 |
| `formRequest.php`        | `form_request_routes.go`         |
| `holiday.php`            | `holiday_routes.go`              |
| `hrdwh.php`              | `hrdwh_routes.go`                |
| `ldap.php`               | `ldap_routes.go`                 |
| `legacyDominoLetter.php` | `legacy_domino_letter_routes.go` |
| `mailadmin.php`          | `mailadmin_routes.go`            |
| `menu.php`               | `menu_routes.go`                 |
| `profile.php`            | `profile_routes.go`              |
| `role.php`               | `role_routes.go`                 |
| `shift.php`              | `shift_routes.go`                |
| `task.php`               | `task_routes.go`                 |
| `user.php`               | `user_routes.go`                 |

## A.3 Other context route files

All module route files live in `internal/app/presentation/<ctx>/routes/` and are aggregated by `internal/routes/v1/apiv1.go`.

| Laravel                         | Go                                                        |
| ------------------------------- | --------------------------------------------------------- |
| `Intra/.../applications.php`    | `internal/app/presentation/intra/routes/intra_routes.go`           |
| `Kampiun/.../kampiun.php`       | `internal/app/presentation/kampiun/routes/kampiun_routes.go`       |
| `Kypas/.../appstoken.php`       | `internal/app/presentation/kypas/routes/appstoken_routes.go`       |
| `Kypas/.../kypas.php`           | `internal/app/presentation/kypas/routes/kypas_routes.go`           |
| `Rfc/.../roles.php`             | `internal/app/presentation/rfc/routes/rfc_role_routes.go`          |
| `Rfc/.../user.php`              | `internal/app/presentation/rfc/routes/rfc_user_routes.go`          |
| `Sygap/.../employee.php`        | `internal/app/presentation/sygap/routes/sygap_employee_routes.go`  |
| `Sygap/.../sygap.php`           | `internal/app/presentation/sygap/routes/sygap_routes.go`           |
| `Tgcast/.../tgcast.php`         | `internal/app/presentation/tgcast/routes/tgcast_routes.go`         |
| `Tgkypas/.../tgkypas.php`       | `internal/app/presentation/tgkypas/routes/tgkypas_routes.go`       |
| `Windigital/.../windigital.php` | `internal/app/presentation/windigital/routes/windigital_routes.go` |
| `routes/api.php`                | `internal/routes/api.go` + `internal/routes/v1/apiv1.go`                    |
| `routes/auth.php`               | `internal/app/presentation/core/routes/auth_routes.go` + `core/routes/sso_routes.go` |
| `routes/guest.php`              | `internal/app/presentation/public/routes.go`                       |
| `routes/web.php`                | `internal/routes/web.go`                                           |
