<script setup lang="ts">
  import { computed } from 'vue'

  const year = computed(() => new Date().getFullYear())
</script>

<template>
  <footer class="border-t-[3px] border-fg bg-surface">
    <div
      class="mx-auto flex w-full max-w-6xl flex-col gap-2 px-6 py-6 font-mono text-xs tracking-widest uppercase sm:flex-row sm:items-center sm:justify-between"
    >
      <p>Vue 3 / Pinia / Tailwind v4 / Axios</p>
      <p>{{ year }} / Neo-Brutalism UI</p>
    </div>
  </footer>
</template>
