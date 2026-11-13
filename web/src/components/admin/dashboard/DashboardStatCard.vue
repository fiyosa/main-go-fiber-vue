<script setup lang="ts">
  defineProps<{
    label: string
    value: string | number
    hint?: string
    tone?: 'ok' | 'warn' | 'danger' | 'neutral'
  }>()

  const toneIconClass: Record<string, string> = {
    ok: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400',
    warn: 'bg-amber-500/10 text-amber-600 dark:text-amber-400',
    danger: 'bg-rose-500/10 text-rose-600 dark:text-rose-400',
    neutral: 'bg-blue-500/10 text-blue-500',
  }
</script>

<template>
  <div
    class="flex flex-col justify-between rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm transition-all hover:shadow-md"
    :data-tone="tone ?? 'neutral'"
  >
    <div class="flex items-center justify-between">
      <span class="text-xs font-medium tracking-wider text-[var(--admin-muted)] uppercase">{{ label }}</span>
      <div
        v-if="$slots.icon"
        class="flex size-10 items-center justify-center rounded-xl"
        :class="toneIconClass[tone ?? 'neutral']"
      >
        <slot name="icon" />
      </div>
    </div>

    <div class="mt-4 flex flex-col gap-1">
      <h3 class="text-2xl font-bold tracking-tight text-[var(--admin-fg)]">{{ value }}</h3>
      <p v-if="hint" class="text-xs text-[var(--admin-muted)]">{{ hint }}</p>
    </div>
  </div>
</template>
