# 10. Background Processing

## 10.1 Execution model

There is no queue, worker pool, or broker. A job is a plain Go function in `internal/app/Application/<domain>/jobs/*.go` that the owning use case or repository runs as a goroutine.

- Fire-and-forget: no persistence, no automatic retry, no dead-letter store.
- The caller starts the goroutine through one helper (`internal/platform/jobs`) so panic recovery, timeout, and shutdown tracking are consistent.
- HTTP handlers never spawn raw goroutines; they call the use case that owns the job.
- Each job group has a semaphore that caps concurrent goroutines; an unbounded job does not spawn unbounded memory.
- The scheduler runs handlers directly, not by enqueueing.

```go
// internal/platform/jobs/run.go
package jobs

var wg sync.WaitGroup

func Go(name string, timeout time.Duration, fn func(context.Context) error) {
    wg.Add(1)
    go func() {
        defer wg.Done()
        defer func() {
            if r := recover(); r != nil {
                slog.Error("job panic", "job", name, "panic", r)
            }
        }()
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        defer cancel()
        if err := fn(ctx); err != nil {
            slog.Error("job failed", "job", name, "error", err)
        }
    }()
}

func Wait(ctx context.Context) error {
    done := make(chan struct{})
    go func() { wg.Wait(); close(done) }()
    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

The process starts no worker; job goroutines appear on demand and are drained by `jobs.Wait` during shutdown.

## 10.2 Job inventory

| Laravel job                               | Go job function                     | Group          |
| ----------------------------------------- | ----------------------------------- | -------------- |
| `AddMemberCompanyTif`                     | `core:member:company-tif:add`       | `processing`   |
| `AddMemberGlobalProtectJob`               | `core:member:global-protect:add`    | `processing`   |
| `BackupHrdwhRpt01Job`                     | `core:hrdwh:rpt01:backup`           | `background`   |
| `GenerateHrdwhBaruJob`                    | `core:hrdwh:baru:generate`          | `processing`   |
| `GenerateHrdwhJabatanJob`                 | `core:hrdwh:jabatan:generate`       | `processing`   |
| `GenerateHrdwhKeluarJob`                  | `core:hrdwh:keluar:generate`        | `processing`   |
| `GenerateHrdwhMutasiJob`                  | `core:hrdwh:mutasi:generate`        | `processing`   |
| `HrdwhRpt01EncryptJob`                    | `core:hrdwh:rpt01:encrypt`          | `processing`   |
| `HrdwhRpt01StatusUpdateJob`               | `core:hrdwh:rpt01:status-update`    | `processing`   |
| `LdapMutationGlobalProtectJob`            | `core:ldap:global-protect:mutation` | `background`   |
| `ProcessMutationEmployeeGlobalProtectJob` | `core:ldap:global-protect:process`  | `background`   |
| `SendTelegramMessageJob`                  | `core:telegram:send`                | `notification` |
| `UpdateLdapCasdoorJob`                    | `core:ldap:casdoor:update`          | `background`   |
| `UpdateLdapEmployeeJob`                   | `core:ldap:employee:update`         | `background`   |
| `UpdateLdapEmployeeTifJob`                | `core:ldap:employee-tif:update`     | `background`   |
| `InactiveRespondenJob`                    | `kypas:responden:inactive`          | `background`   |
| `UpdateMailJob`                           | `mailadmin:mail:update`             | `background`   |
| `ProcessUpdatePhotoMicrosoftGraphJob`     | `microsoftgraph:photo:update`       | `processing`   |
| `UpdateIncidentJob`                       | `sygap:incident:update`             | `processing`   |
| `UpdateRequestFulfillmentJob`             | `sygap:request:update`              | `processing`   |
| `UpdateTaskJob`                           | `sygap:task:update`                 | `processing`   |
| `AddLdapMitraTgkypasJob`                  | `tgkypas:ldap:mitra:add`            | `background`   |
| `DownloadToLocalJustificationsJob`        | `tgkypas:justification:download`    | `document`     |
| `SyncMitraNcxJob`                         | `tgkypas:mitra:sync`                | `processing`   |
| `PostUpdateUniversalInboxJob`             | `universal:inbox:update`            | `processing`   |
| `DocumentCopyToFolderJob`                 | `windigital:document:copy`          | `document`     |
| `DownloadDocumentWindigitalJob`           | `windigital:document:download`      | `document`     |
| `DownloadDocumentWindigitalUserJob`       | `windigital:document:download-user` | `document`     |

## 10.3 Scheduler entries

| Context                       | Command                                             | Schedule                   | Notes                                    |
| ----------------------------- | --------------------------------------------------- | -------------------------- | ---------------------------------------- |
| Core                          | `app:activate-automate-wip`                         | Daily                      | `withoutOverlapping`                     |
| Core                          | `app:reminder-new-rpt-info`                         | Hourly at `:45`            | `withoutOverlapping`                     |
| Core                          | `app:attendance-generate --force`                   | Monthly                    | `withoutOverlapping`                     |
| Core                          | `app:form-request-expired`                          | Monthly                    | `withoutOverlapping`                     |
| Core                          | `app:reminder-aux-info`                             | Every 10 minutes           | `withoutOverlapping`                     |
| Core                          | `activitylog:clean --force`                         | Daily                      | Retention cleanup                        |
| Core                          | `backup:run --only-db`                              | Daily at `21:00`           | Replace with `pg_dump` or managed backup |
| Insera                        | `app:insera-update-incident --force`                | Every 10 minutes           |                                          |
| Insera                        | `app:insera-reopen-incident --force`                | Every 30 minutes           |                                          |
| Insera                        | `app:insera-auto-flag-incident --force`             | Daily                      |                                          |
| Kypas                         | `app:reminder-approval-admin-kypas`                 | Every 30 minutes           |                                          |
| Sygap                         | `app:sygap-incident-update-wip`                     | Every 3 minutes            |                                          |
| Sygap                         | `app:sygap-task-update-wip`                         | Every 2 hours              |                                          |
| Sygap                         | `app:sygap-reminder-request-fulfillment-inprogress` | Hourly at `:15`            |                                          |
| Sygap                         | `app:sygap-reminder-task-inprogress`                | Hourly                     |                                          |
| Sygap                         | `app:sygap-reminder-incident-inprogress`            | Hourly                     |                                          |
| WebAuth                       | `app:auth-token-reminder-expiration`                | Daily at `09:00`           |                                          |
| WebAuth                       | `model:prune --model=AuthLog`                       | Monthly                    | Replace with retention command           |
| Windigital                    | `app:windigital-download-document`                  | Daily                      |                                          |
| MicrosoftGraph                | Current schedule duplicates Insera entries          | Every 10/30 minutes, daily | Verify intent before porting             |
| Mailadmin, Tgkypas, Universal | No active entries                                   | -                          | Universal entry is commented out         |

Implementation:

- One scheduler goroutine started from `internal/bootstrap/scheduler.go`, which uses `internal/app/infrastructure/scheduler` (robfig/cron v3).
- On startup, acquire a PostgreSQL advisory lock (`pg_try_advisory_lock`); only the leader registers cron entries. Release the lock on shutdown.
- Each entry runs its handler directly inside a goroutine (`jobs.Go`), not through a queue.
- Use a per-entry PostgreSQL advisory lock or an in-process mutex to enforce `withoutOverlapping`.
- Record last run and last error in `scheduled_runs` for observability.
- The scheduler must never run inside the `console` binary.

## 10.4 Console commands

| Laravel command                              | Go console command                               |
| -------------------------------------------- | ------------------------------------------------ |
| `app:activate-automate-wip`                  | `console activate-automate-wip`                  |
| `app:attendance-generate`                    | `console attendance-generate`                    |
| `app:attendance-run-import`                  | `console attendance-run-import`                  |
| `app:domino-letter-import`                   | `console domino-letter-import`                   |
| `app:form-request-expired`                   | `console form-request-expired`                   |
| `app:initiate-member-casdoor`                | `console initiate-member-casdoor`                |
| `app:initiate-member-company-tif`            | `console initiate-member-company-tif`            |
| `app:initiate-member-global-protect`         | `console initiate-member-global-protect`         |
| `app:ldap-update-employee`                   | `console ldap-update-employee`                   |
| `app:reminder-aux-info`                      | `console reminder-aux-info`                      |
| `app:reminder-new-rpt`                       | `console reminder-new-rpt`                       |
| `app:insera-*`                               | `console insera-*`                               |
| `app:reminder-approval-admin-kypas`          | `console reminder-approval-admin-kypas`          |
| `app:reminder-approval-admin-kypas-telegram` | `console reminder-approval-admin-kypas-telegram` |
| `app:microsoft-graph-update-photo`           | `console microsoft-graph-update-photo`           |
| `app:sygap-*`                                | `console sygap-*`                                |
| `app:initiate-ldap-group-mitra-tgkypas`      | `console initiate-ldap-group-mitra-tgkypas`      |
| `app:download-to-local-justifications`       | `console download-to-local-justifications`       |
| `app:sync-mitra-ncx`                         | `console sync-mitra-ncx`                         |
| `app:universal-post-update-inbox`            | `console universal-post-update-inbox`            |
| `app:auth-token-reminder-expiration`         | `console auth-token-reminder-expiration`         |
| `app:windigital-download-document`           | `console windigital-download-document`           |
| `SyncGraphCalendarCommand`                   | `console sync-graph-calendar`                    |

Database commands (`console migrate`, `console db seed`, `console db backup`) are one-off console commands defined in [Section 9.7](09-data-layer.md#97-database-commands); they are not background jobs.

## 10.5 Reliability

- Retries: none automatic. Jobs that call flaky integrations retry inside the job with a bounded attempt count; everything else fails once and is logged.
- Backoff: implemented inside the job when it retries.
- Failed jobs: no dead-letter store. Failures are logged with job name and error; rerun by calling the owning use case or the scheduled entry again.
- Idempotency: still required. A job can be started twice by two requests or two scheduler ticks, so writes must be safe to repeat.
- Timeouts: every job gets a context deadline from the helper; handlers must accept and respect the context.
- Concurrency: a per-group semaphore caps in-flight goroutines so a burst cannot exhaust memory.
- Logging: structured logs with job name, group, duration, and correlation ID.
