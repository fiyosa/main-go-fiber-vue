<script setup lang="ts">
  import { computed } from 'vue'

  const year = computed(() => new Date().getFullYear())
</script>

<template>
  <footer class="border-t-[3px] border-fg bg-surface">
    <div
      class="mx-auto flex w-full max-w-6xl flex-col items-center justify-between gap-1 px-4 py-5 font-mono text-xs font-bold tracking-widest uppercase sm:flex-row sm:gap-3 sm:px-6 sm:py-6"
    >
      <p>Portfolio</p>
      <p>&copy; {{ year }}</p>
    </div>
  </footer>
</template>
