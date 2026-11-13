<script setup lang="ts">
  import { computed, watch } from 'vue'

  const page = defineModel<number>('page', { required: true })
  const limit = defineModel<number>('limit', { required: true })

  const props = withDefaults(
    defineProps<{
      total: number
      limitOptions?: number[]
    }>(),
    { limitOptions: () => [25, 50, 100] }
  )

  const totalPages = computed(() => Math.max(1, Math.ceil(props.total / limit.value)))
  const currentPage = computed(() => Math.min(page.value, totalPages.value))

  watch(limit, () => {
    page.value = 1
  })

  watch(totalPages, (total) => {
    if (page.value > total) page.value = total
  })

  const startItem = computed(() => (props.total === 0 ? 0 : (currentPage.value - 1) * limit.value + 1))
  const endItem = computed(() => Math.min(currentPage.value * limit.value, props.total))
</script>

<template>
  <div
    class="flex flex-col gap-4 rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 py-3 shadow-sm sm:flex-row sm:items-center sm:justify-between"
  >
    <div class="flex flex-wrap items-center gap-2 text-xs text-[var(--admin-muted)]">
      <span>Rows per page:</span>
      <select
        v-model.number="limit"
        class="rounded-lg border border-[var(--admin-border)] bg-[var(--admin-surface)] px-2.5 py-1.5 text-xs text-[var(--admin-fg)] focus:border-blue-500 outline-none"
      >
        <option v-for="option in limitOptions" :key="option" :value="option">{{ option }}</option>
      </select>
      <span>
        Showing {{ startItem }} to {{ endItem }} of {{ total }} {{ total === 1 ? 'entry' : 'entries' }}
      </span>
    </div>

    <div class="flex items-center gap-1.5">
      <button
        type="button"
        class="inline-flex min-h-8 min-w-8 items-center justify-center rounded-lg border border-[var(--admin-border)] bg-[var(--admin-surface)] text-[var(--admin-fg)] transition-colors hover:bg-[var(--admin-border)]/50 disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="currentPage <= 1"
        title="First Page"
        aria-label="First Page"
        @click="page = 1"
      >
        ««
      </button>
      <button
        type="button"
        class="inline-flex min-h-8 min-w-8 items-center justify-center rounded-lg border border-[var(--admin-border)] bg-[var(--admin-surface)] text-[var(--admin-fg)] transition-colors hover:bg-[var(--admin-border)]/50 disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="currentPage <= 1"
        title="Previous Page"
        aria-label="Previous Page"
        @click="page--"
      >
        ‹
      </button>
      <span class="px-3 py-1 font-medium text-[var(--admin-fg)]">Page {{ currentPage }} of {{ totalPages }}</span>
      <button
        type="button"
        class="inline-flex min-h-8 min-w-8 items-center justify-center rounded-lg border border-[var(--admin-border)] bg-[var(--admin-surface)] text-[var(--admin-fg)] transition-colors hover:bg-[var(--admin-border)]/50 disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="currentPage >= totalPages"
        title="Next Page"
        aria-label="Next Page"
        @click="page++"
      >
        ›
      </button>
      <button
        type="button"
        class="inline-flex min-h-8 min-w-8 items-center justify-center rounded-lg border border-[var(--admin-border)] bg-[var(--admin-surface)] text-[var(--admin-fg)] transition-colors hover:bg-[var(--admin-border)]/50 disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="currentPage >= totalPages"
        title="Last Page"
        aria-label="Last Page"
        @click="page = totalPages"
      >
        »»
      </button>
    </div>
  </div>
</template>
