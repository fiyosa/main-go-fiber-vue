<script setup lang="ts">
  import { RouterLink } from 'vue-router'

  import { BrutalBadge, BrutalCard } from '../../components/ui'
</script>

<template>
  <div class="flex min-h-[50vh] items-center justify-center">
    <BrutalCard accent="hot-red" tone="surface" class="w-full max-w-lg -rotate-1 text-center">
      <template #header>Route error</template>

      <div class="flex flex-col items-center gap-4">
        <BrutalBadge tone="vivid-yellow" :rotation="-3">404</BrutalBadge>

        <p class="text-6xl leading-none font-bold tracking-tight sm:text-7xl">Lost?</p>

        <h1 class="text-xl font-bold">Page not found</h1>
        <p class="text-sm leading-relaxed">The route you requested does not exist or was moved.</p>

        <RouterLink
          to="/"
          class="brutal-press inline-flex min-h-12 items-center justify-center border-[3px] border-fg bg-mint px-6 font-mono text-sm font-bold tracking-widest text-black uppercase shadow-brutal-sm"
        >
          Back home
        </RouterLink>
      </div>
    </BrutalCard>
  </div>
</template>
