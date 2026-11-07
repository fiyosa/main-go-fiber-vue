<script setup lang="ts">
  import { RouterLink } from 'vue-router'
</script>

<template>
  <main class="min-h-screen bg-slate-950 px-6 py-16 text-slate-100">
    <div class="mx-auto max-w-2xl space-y-4">
      <p class="text-xs font-semibold tracking-[0.2em] text-teal-400 uppercase">About</p>
      <h1 class="text-3xl font-bold">Frontend layout</h1>
      <p class="text-sm leading-relaxed text-slate-400">
        Folders follow the project convention: <code>api</code> for endpoints, <code>lib</code> for shared clients,
        <code>store</code> for Pinia, <code>pages</code> for routed views, and <code>config</code> for static settings.
      </p>
      <RouterLink
        to="/"
        class="inline-block rounded-md bg-teal-500 px-4 py-2 text-sm font-semibold text-slate-950 transition hover:bg-teal-400"
      >
        Back home
      </RouterLink>
    </div>
  </main>
</template>
