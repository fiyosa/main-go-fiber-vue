<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { useDashboard } from './useDashboard'
  import DashboardStatCard from '../../components/admin/dashboard/DashboardStatCard.vue'
  import RecentRunsTable from '../../components/admin/dashboard/RecentRunsTable.vue'
  import FailedRunsList from '../../components/admin/dashboard/FailedRunsList.vue'
  import RunDetailModal from '../../components/admin/dashboard/RunDetailModal.vue'
  import type { CronRunEntity } from '../../api/cron'

  const router = useRouter()
  const {
    loading,
    error,
    healthStatus,
    readyStatus,
    crons,
    currentUser,
    cronStats,
    runsStats,
    recentRuns,
    failedRuns,
    loadDashboard,
  } = useDashboard()

  const selectedRun = ref<CronRunEntity | null>(null)

  onMounted(() => {
    loadDashboard()
  })

  const handleViewRun = (run: CronRunEntity) => {
    selectedRun.value = run
  }
</script>

<template>
  <div class="flex flex-col gap-8 pb-12">
    <!-- Welcome Banner -->
    <div
      class="flex flex-col justify-between gap-4 rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm sm:flex-row sm:items-center"
    >
      <div class="flex flex-col gap-1">
        <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Dashboard Overview</h1>
        <p class="text-sm text-[var(--admin-muted)]">
          Welcome back, <span class="font-semibold text-[var(--admin-fg)]">{{ currentUser?.email || 'Administrator' }}</span>. System monitoring and automation metrics.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          type="button"
          class="inline-flex min-h-10 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 text-xs font-semibold shadow-sm transition-colors hover:bg-[var(--admin-border)]/50"
          :disabled="loading"
          @click="loadDashboard"
        >
          {{ loading ? 'Refreshing...' : 'Refresh Data' }}
        </button>
      </div>
    </div>

    <!-- Error banner if any -->
    <div
      v-if="error"
      class="rounded-2xl border border-rose-500/30 bg-rose-500/10 p-4 text-xs font-semibold text-rose-600 dark:text-rose-400"
    >
      {{ error }}
    </div>

    <!-- Top Metrics Grid -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
      <DashboardStatCard
        label="System Health"
        :value="readyStatus === 'ready' ? 'Ready' : (readyStatus ?? 'Checking...') "
        :hint="`Liveness: ${healthStatus ?? 'unknown'}`"
        :tone="readyStatus === 'ready' ? 'ok' : 'danger'"
      >
        <template #icon>
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
        </template>
      </DashboardStatCard>

      <DashboardStatCard
        label="Cron Jobs"
        :value="cronStats.total"
        :hint="`${cronStats.enabled} enabled, ${cronStats.disabled} paused`"
        tone="neutral"
      >
        <template #icon>
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </template>
      </DashboardStatCard>

      <DashboardStatCard
        label="Runs (Past 24h)"
        :value="runsStats.total"
        :hint="`${runsStats.success} success, ${runsStats.failed} failed`"
        :tone="runsStats.failed > 0 ? 'warn' : 'ok'"
      >
        <template #icon>
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </template>
      </DashboardStatCard>

      <DashboardStatCard
        label="Failure Rate"
        :value="`${runsStats.failureRate}%`"
        :hint="`${runsStats.failed} failure(s) in window`"
        :tone="runsStats.failureRate > 0 ? 'danger' : 'ok'"
      >
        <template #icon>
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </template>
      </DashboardStatCard>
    </div>

    <!-- Issues / Failures Section -->
    <div class="flex flex-col gap-4">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-base font-bold tracking-tight">Active Issues & Failures</h2>
          <p class="text-xs text-[var(--admin-muted)]">Failed executions requiring attention</p>
        </div>
      </div>

      <FailedRunsList :runs="failedRuns" :crons="crons" @view="handleViewRun" />
    </div>

    <!-- Recent Activity Section -->
    <div class="flex flex-col gap-4">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-base font-bold tracking-tight">Recent Cron Executions</h2>
          <p class="text-xs text-[var(--admin-muted)]">Latest automation runs across all cron tasks</p>
        </div>

        <button
          type="button"
          class="text-xs font-semibold text-blue-500 hover:text-blue-600"
          @click="router.push({ name: 'admin-cron' })"
        >
          Manage Cron Jobs →
        </button>
      </div>

      <RecentRunsTable :runs="recentRuns" :crons="crons" :loading="loading" @view="handleViewRun" />
    </div>

    <!-- Run Detail Modal -->
    <RunDetailModal
      :run="selectedRun"
      :cronName="selectedRun ? crons.find((c) => c.id === selectedRun?.cron_id)?.name ?? selectedRun.cron_id : ''"
      @close="selectedRun = null"
    />
  </div>
</template>
