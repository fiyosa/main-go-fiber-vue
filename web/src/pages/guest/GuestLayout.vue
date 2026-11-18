<script setup lang="ts">
  import { RouterView } from 'vue-router'

  import { AppFooter, AppHeader } from '../../components/guest/layout'
</script>

<template>
  <div class="guest-brutal flex min-h-screen flex-col overflow-x-hidden text-fg">
    <a
      href="#main-content"
      class="sr-only focus:not-sr-only focus:absolute focus:top-3 focus:left-3 focus:z-50 focus:border-[3px] focus:border-fg focus:bg-vivid-yellow focus:px-4 focus:py-2 focus:font-mono focus:text-xs focus:font-bold focus:tracking-widest focus:text-black focus:uppercase"
    >
      Skip to content
    </a>

    <AppHeader />

    <main id="main-content" class="mx-auto w-full max-w-6xl flex-1 px-4 py-8 sm:px-6 sm:py-10">
      <RouterView />
    </main>

    <AppFooter />
  </div>
</template>
