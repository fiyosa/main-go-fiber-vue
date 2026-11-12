<script setup lang="ts">
  defineProps<{
    title: string
    value: string | number
    change?: string
    trend?: 'up' | 'down' | 'neutral'
  }>()
</script>

<template>
  <div
    class="flex flex-col justify-between rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm transition-all hover:shadow-md"
  >
    <div class="flex items-center justify-between">
      <span class="text-xs font-medium tracking-wider text-[var(--admin-muted)] uppercase">{{ title }}</span>
      <div v-if="$slots.icon" class="flex size-10 items-center justify-center rounded-xl bg-blue-500/10 text-blue-500">
        <slot name="icon" />
      </div>
    </div>

    <div class="mt-4 flex items-baseline justify-between">
      <h3 class="text-2xl font-bold tracking-tight text-[var(--admin-fg)]">{{ value }}</h3>

      <span
        v-if="change"
        class="inline-flex items-center gap-1 rounded-lg px-2 py-1 text-xs font-semibold"
        :class="{
          'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400': trend === 'up',
          'bg-rose-500/10 text-rose-600 dark:text-rose-400': trend === 'down',
          'bg-slate-500/10 text-[var(--admin-muted)]': trend === 'neutral' || !trend,
        }"
      >
        {{ change }}
      </span>
    </div>
  </div>
</template>
