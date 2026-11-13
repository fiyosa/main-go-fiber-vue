<spec id="admin-dashboard-monitoring">
# Specification: Real-Data Admin Dashboard Monitoring

## Overview
Recompose `web/src/pages/admin/AdminHome.vue` into a robust, real-data monitoring dashboard following the hierarchy: **status → activity → issues → details**. No mock numbers, no trend percentages, no charts.

## Objectives
1. **System Health & Connectivity Status**: Fetch liveness (`/api/v1/healthz`) and readiness (`/api/v1/readyz` - DB ping).
2. **Real Cron Automation KPIs**: Fetch crons (`/api/v1/crons`) and execution runs (`/api/v1/crons/runs`). Compute real counts: Total Cron Jobs, Active/Paused count, Last 24h Runs breakdown (Success, Failed, Running).
3. **Identity Context**: Display current logged-in admin identity (`/api/v1/auth/me`).
4. **Recent Activity Feed**: Show the latest real execution runs (`/api/v1/crons/runs` sorted by started_at desc) with status badges and quick view trigger for run output/error details (reusing existing modal patterns from `CronManager.vue`).
5. **Issues / Failures Section**: Highlight any failed cron runs or unhealthy system states immediately, with clear action links to investigate.

## Component Layout & Sections
- **Welcome Banner**: "Dashboard Overview", dynamic greeting with logged-in user email and role badges, "Refresh Data" button.
- **Top Metrics Grid (4 Stat Cards)**:
  1. *System Health*: Ready / Degraded status based on `/api/v1/readyz`.
  2. *Cron Jobs*: Total count + active / paused breakdown.
  3. *Execution Runs (24h)*: Total executed runs in past 24h (or latest available 100 runs window).
  4. *Failure Rate / Issues*: Count of failed runs in the fetched history.
- **Recent Activity Table**: Real list of recent cron execution runs with status badge, trigger type, duration, started timestamp, and "View Output" button.
- **Quick Actions / Navigation Shortcuts**: Fast links to Cron Manager (`/admin/work/cron`) and Role Management / Register (`/admin/auth/register`).

## Verification Plan
- `npm run test` (Vitest)
- `npm run build` (Vue TSC & Vite bundle)
</spec>
