<script setup lang="ts">
  import type { CronRunEntity } from '../../../api/cron'
  import { runStatusBadgeClass } from '../../../utils/dashboard'
  import { formatDateTime, formatDuration } from '../../../utils/format'

  defineProps<{
    run: CronRunEntity | null
    cronName: string
  }>()

  const emit = defineEmits<{
    close: []
  }>()
</script>

<template>
  <div
    v-if="run"
    class="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-slate-950/60 p-4 backdrop-blur-xs"
    @click.self="emit('close')"
  >
    <div
      role="dialog"
      aria-modal="true"
      class="relative my-8 w-full max-w-2xl rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-2xl"
    >
      <div class="flex items-center justify-between border-b border-[var(--admin-border)] pb-4">
        <div class="flex items-center gap-3">
          <h2 class="text-base font-bold tracking-tight">{{ cronName }}</h2>
          <span
            class="rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase"
            :class="runStatusBadgeClass(run.status)"
          >
            {{ run.status }}
          </span>
        </div>
        <button
          type="button"
          class="text-[var(--admin-muted)] transition-colors hover:text-[var(--admin-fg)]"
          aria-label="Close"
          @click="emit('close')"
        >
          ✕
        </button>
      </div>

      <div class="mt-4 flex flex-col gap-4 text-xs">
        <div class="grid grid-cols-2 gap-4 rounded-xl border border-[var(--admin-border)] bg-[var(--admin-border)]/10 p-3">
          <div>
            <span class="text-[var(--admin-muted)]">Trigger Type:</span>
            <strong class="capitalize">{{ run.trigger_type }}</strong>
          </div>
          <div>
            <span class="text-[var(--admin-muted)]">Duration:</span>
            <strong class="font-mono">{{ formatDuration(run.duration_ms) }}</strong>
          </div>
          <div>
            <span class="text-[var(--admin-muted)]">Started At:</span> {{ formatDateTime(run.started_at) }}
          </div>
          <div>
            <span class="text-[var(--admin-muted)]">Finished At:</span> {{ formatDateTime(run.finished_at) }}
          </div>
        </div>

        <div v-if="run.error" class="flex flex-col gap-1.5">
          <label class="font-semibold text-rose-500">Error Message</label>
          <pre
            class="rounded-xl border border-rose-500/30 bg-rose-500/10 p-3 font-mono text-xs whitespace-pre-wrap text-rose-600 dark:text-rose-400"
            >{{ run.error }}</pre
          >
        </div>

        <div class="flex flex-col gap-1.5">
          <label class="font-semibold">Console Output</label>
          <pre
            class="max-h-60 overflow-y-auto rounded-xl border border-[var(--admin-border)] bg-slate-950 p-4 font-mono text-xs whitespace-pre-wrap text-emerald-400"
            >{{ run.output || '(No console output)' }}</pre
          >
        </div>
      </div>
    </div>
  </div>
</template>
