<script setup lang="ts">
  import { RouterLink } from 'vue-router'
</script>

<template>
  <main class="flex min-h-screen items-center justify-center bg-slate-950 px-6 text-slate-100">
    <div class="space-y-4 text-center">
      <p class="text-7xl font-black text-teal-400">404</p>
      <h1 class="text-xl font-semibold">Page not found</h1>
      <p class="text-sm text-slate-400">The route you requested does not exist.</p>
      <RouterLink
        to="/"
        class="inline-block rounded-md bg-teal-500 px-4 py-2 text-sm font-semibold text-slate-950 transition hover:bg-teal-400"
      >
        Back home
      </RouterLink>
    </div>
  </main>
</template>
