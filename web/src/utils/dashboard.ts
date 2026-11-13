import type { CronEntity, CronRunEntity } from '../api/cron'

export interface RunWindowStats {
  total: number
  success: number
  failed: number
  running: number
  skipped: number
  failureRate: number
}

/// Counts how many cron jobs are currently enabled.
export const countEnabledCrons = (crons: CronEntity[]): number =>
  crons.filter((cron) => cron.enabled).length

/// Aggregates run statuses. When `windowMs` is provided, only runs started
/// within that window (relative to `now`) are counted.
export const summarizeRuns = (
  runs: CronRunEntity[],
  windowMs?: number,
  now: number = Date.now()
): RunWindowStats => {
  const since = windowMs ? now - windowMs : null
  const scoped = since
    ? runs.filter((run) => new Date(run.started_at).getTime() >= since)
    : runs

  const total = scoped.length
  const count = (status: string) => scoped.filter((run) => run.status === status).length

  return {
    total,
    success: count('success'),
    failed: count('failed'),
    running: count('running'),
    skipped: count('skipped'),
    failureRate: total === 0 ? 0 : Math.round((count('failed') / total) * 100),
  }
}

/// Resolves a cron id to its display name, falling back to the raw id.
export const getCronName = (crons: CronEntity[], cronId: string): string =>
  crons.find((cron) => cron.id === cronId)?.name ?? cronId

/// Returns a newest-first copy of the runs list without mutating the input.
export const sortRunsByStartedAtDesc = (runs: CronRunEntity[]): CronRunEntity[] =>
  [...runs].sort(
    (a, b) => new Date(b.started_at).getTime() - new Date(a.started_at).getTime()
  )

/// Maps a run status to the badge classes used across the dashboard.
export const runStatusBadgeClass = (status: string): string => {
  switch (status) {
    case 'success':
      return 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
    case 'failed':
      return 'bg-rose-500/10 text-rose-600 dark:text-rose-400'
    case 'running':
      return 'bg-amber-500/10 text-amber-600 dark:text-amber-400'
    default:
      return 'bg-slate-500/10 text-[var(--admin-muted)]'
  }
}
