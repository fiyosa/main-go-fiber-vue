<script setup lang="ts">
  import type { CronEntity, CronRunEntity } from '../../../api/cron'
  import { getCronName, runStatusBadgeClass } from '../../../utils/dashboard'
  import { formatTimeAgo } from '../../../utils/format'

  defineProps<{
    runs: CronRunEntity[]
    crons: CronEntity[]
  }>()

  const emit = defineEmits<{
    view: [run: CronRunEntity]
  }>()
</script>

<template>
  <div class="rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm">
    <div v-if="runs.length === 0" class="flex flex-col items-center justify-center py-8 text-center">
      <div class="flex size-10 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-500">
        <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <p class="mt-3 text-sm font-semibold text-[var(--admin-fg)]">No failed runs in the recent history</p>
      <p class="mt-1 text-xs text-[var(--admin-muted)]">Failed cron executions will surface here for quick triage.</p>
    </div>

    <div v-else class="flex flex-col gap-3">
      <div
        v-for="run in runs"
        :key="run.id"
        class="flex flex-col justify-between gap-3 rounded-xl border border-rose-500/20 bg-rose-500/5 p-4 sm:flex-row sm:items-center"
      >
        <div class="flex min-w-0 items-start gap-3">
          <span
            class="mt-0.5 inline-block shrink-0 rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase"
            :class="runStatusBadgeClass(run.status)"
          >
            {{ run.status }}
          </span>
          <div class="min-w-0">
            <p class="truncate text-xs font-semibold text-[var(--admin-fg)]">
              {{ getCronName(crons, run.cron_id) }}
            </p>
            <p class="mt-0.5 truncate font-mono text-[11px] text-rose-600 dark:text-rose-400">
              {{ run.error || 'No error message captured.' }}
            </p>
            <p class="mt-0.5 text-[11px] text-[var(--admin-muted)]">{{ formatTimeAgo(run.started_at) }}</p>
          </div>
        </div>

        <button
          type="button"
          class="inline-flex min-h-9 shrink-0 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3 text-xs font-semibold transition-colors hover:bg-[var(--admin-border)]/50"
          @click="emit('view', run)"
        >
          Investigate
        </button>
      </div>
    </div>
  </div>
</template>
