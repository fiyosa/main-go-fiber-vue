<script setup lang="ts">
  import { computed } from 'vue'

  const year = computed(() => new Date().getFullYear())
</script>

<template>
  <footer class="border-t-[3px] border-fg bg-surface">
    <div
      class="mx-auto flex w-full max-w-6xl items-center justify-between gap-3 px-6 py-6 font-mono text-xs font-bold tracking-widest uppercase"
    >
      <p>Portfolio</p>
      <p>&copy; {{ year }}</p>
    </div>
  </footer>
</template>
