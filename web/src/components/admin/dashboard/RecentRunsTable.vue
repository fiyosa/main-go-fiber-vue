<script setup lang="ts">
  import type { CronEntity, CronRunEntity } from '../../../api/cron'
  import { getCronName, runStatusBadgeClass } from '../../../utils/dashboard'
  import { formatDateTime, formatDuration } from '../../../utils/format'

  defineProps<{
    runs: CronRunEntity[]
    crons: CronEntity[]
    loading?: boolean
  }>()

  const emit = defineEmits<{
    view: [run: CronRunEntity]
  }>()
</script>

<template>
  <div class="overflow-x-auto rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] shadow-sm">
    <div v-if="loading" class="py-12 text-center text-xs text-[var(--admin-muted)]">
      Loading execution runs...
    </div>

    <div v-else-if="runs.length === 0" class="flex flex-col items-center justify-center p-12 text-center">
      <p class="text-sm font-semibold text-[var(--admin-fg)]">No execution runs recorded yet</p>
      <p class="mt-1 text-xs text-[var(--admin-muted)]">
        Cron execution logs and console outputs will appear here.
      </p>
    </div>

    <table v-else class="w-full text-left text-xs">
      <thead class="border-b border-[var(--admin-border)] bg-[var(--admin-border)]/20 text-[var(--admin-muted)] uppercase">
        <tr>
          <th class="px-4 py-3 font-semibold">Status</th>
          <th class="px-4 py-3 font-semibold">Cron Job</th>
          <th class="px-4 py-3 font-semibold">Trigger</th>
          <th class="px-4 py-3 font-semibold">Duration</th>
          <th class="px-4 py-3 font-semibold">Started At</th>
          <th class="px-4 py-3 font-semibold">Actions</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-[var(--admin-border)]">
        <tr v-for="run in runs" :key="run.id" class="transition-colors hover:bg-[var(--admin-border)]/10">
          <td class="px-4 py-3.5">
            <span
              class="rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase"
              :class="runStatusBadgeClass(run.status)"
            >
              {{ run.status }}
            </span>
          </td>
          <td class="px-4 py-3.5 font-semibold text-[var(--admin-fg)]">
            {{ getCronName(crons, run.cron_id) }}
          </td>
          <td class="px-4 py-3.5 font-medium capitalize">{{ run.trigger_type }}</td>
          <td class="px-4 py-3.5 font-mono">{{ formatDuration(run.duration_ms) }}</td>
          <td class="px-4 py-3.5 text-[var(--admin-muted)]">{{ formatDateTime(run.started_at) }}</td>
          <td class="px-4 py-3.5">
            <button
              type="button"
              class="rounded-lg border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 py-1 text-xs font-semibold transition-colors hover:bg-[var(--admin-border)]/50"
              @click="emit('view', run)"
            >
              View Output
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
