# 3. Current Laravel Architecture

## 3.1 Layered structure

| Layer          | Path                 | Responsibility                                                                                                              |
| -------------- | -------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| Presentation   | `app/Presentasion`   | HTTP controllers, requests, resources, routes, views                                                                        |
| Application    | `app/Application`    | Use cases, repositories, criteria, jobs, exports, imports, commands                                                         |
| Domain         | `app/Domain`         | Eloquent models, enums, policies, casts, repository and unused API contracts (legacy; not carried into the Go domain layer) |
| Infrastructure | `app/Infrastructure` | Laravel adapters, middleware, providers, external integrations, services                                                    |

> Note: `Presentasion` is a legacy spelling. The Go rewrite uses the correctly spelled `presentation` folder, but the Laravel source keeps `Presentasion`.

## 3.2 Bounded contexts

| Laravel context  | Go context       | Notes                                                                                                        |
| ---------------- | ---------------- | ------------------------------------------------------------------------------------------------------------ |
| `Alma`           | `core`           | Main domain: auth, users, roles, menus, attendance, shifts, tasks, documents, media, forms, approvals, HRDWH |
| `Sygap`          | `sygap`          | Incident, request fulfillment, task, organization, service                                                   |
| `Kypas`          | `kypas`          | Contractor/employee management, approvals, NDE, app tokens                                                   |
| `Tgkypas`        | `tgkypas`        | Contractor users, contracts, access requests                                                                 |
| `Kampiun`        | `kampiun`        | Document hits and recap                                                                                      |
| `Intra`          | `intra`          | Application approval and comments                                                                            |
| `Rfc`            | `rfc`            | RFC users and roles                                                                                          |
| `Windigital`     | `windigital`     | Document download and copy                                                                                   |
| `Tgcast`         | `tgcast`         | Mail broadcast                                                                                               |
| `Mailadmin`      | `mailadmin`      | Mail profile operations                                                                                      |
| `Evolution`      | `evolution`      | Legacy user data                                                                                             |
| `NewCuti`        | `newcuti`        | Leave data                                                                                                   |
| `Lemon`          | `lemon`          | HRDWH legacy                                                                                                 |
| `WebAuth`        | `webauth`        | Auth logs and token reminders                                                                                |
| `Action`         | `action`         | Legacy action user                                                                                           |
| `Universal`      | `universal`      | Universal inbox integration                                                                                  |
| `Insera`         | `insera`         | Insera incident sync                                                                                         |
| `MicrosoftGraph` | `microsoftgraph` | Microsoft Graph integration                                                                                  |
| `Tgsso`          | `tgsso`          | SSO user data                                                                                                |

## 3.3 Admin route loading

Laravel auto-loads every PHP file under `app/Presentasion/Admin/*/Http/Routes/` inside this group:

```php
Route::middleware(['auth:admin', 'subtitute.bindings'])
    ->prefix('admin/api')
    ->name('admin.')
    ->group(...);
```

The Go rewrite must register routes explicitly. Route files live in each context's `routes` package and are aggregated by `internal/routes/v1/apiv1.go`.

## 3.4 Authentication today

- Guard `admin` uses Passport (`config/auth.php`).
- `User` implements `OAuthenticatable` and uses `HasApiTokens`.
- Current Laravel login validates credentials against LDAP, then finds the local user and runs a TOTP second factor (removed in the rewrite).
- `verify-2fa` completes the second-factor challenge (not ported).
- Access token lifetime: 3 hours. Refresh token lifetime: 30 days. Personal access tokens: 3 hours.
- Scope defined today: `impersonation`.
- Logout deletes all tokens for the user.

The Go target removes LDAP credential checks and the second factor: login verifies the local `users.password` hash and issues tokens immediately ([Section 8.5](08-authentication.md#85-login-flow)). LDAP remains only for directory synchronization.

## 3.5 Runtime components

| Component    | Current implementation                      |
| ------------ | ------------------------------------------- |
| Database     | PostgreSQL; multiple PostgreSQL connections |
| Migrations   | 155 Laravel migrations                      |
| Cache        | File                                     |
| Queue        | Database queue                           |
| Scheduler    | Laravel console schedule                    |
| Broadcast    | Pusher                                      |
| Storage      | Local disks + S3                            |
| Mail         | Laravel mail + templates                    |
| Excel        | maatwebsite/excel                           |
| PDF          | barryvdh/laravel-dompdf + Blade views       |
| Telegram     | irazasyed/telegram-bot-sdk                  |
| Activity log | spatie/laravel-activitylog                  |
| RBAC         | spatie/laravel-permission                   |
| Hash IDs     | vinkla/hashids                              |
