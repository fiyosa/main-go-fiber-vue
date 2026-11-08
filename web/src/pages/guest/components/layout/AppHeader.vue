<script setup lang="ts">
  import { computed } from 'vue'
  import { RouterLink, useRoute } from 'vue-router'

  import { useUIStore } from '../../../../store/ui'

  const ui = useUIStore()
  const route = useRoute()

  const links = [
    { name: 'home', to: '/', label: 'Home' },
    { name: 'about', to: '/about', label: 'About' },
  ] as const

  // Base shape only; the background swaps so the active route is an inverted
  // block instead of a color-only hint.
  const linkClasses =
    'brutal-press inline-flex min-h-11 items-center border-[3px] border-fg px-4 font-mono text-xs font-bold tracking-widest uppercase shadow-brutal-sm'

  const linkTone = (name: string) => (route.name === name ? 'bg-fg text-bg' : 'bg-surface text-fg')

  const isDark = computed(() => ui.theme === 'dark')

  const toggleTheme = () => {
    ui.toggleTheme()
  }

  const toggleMenu = () => {
    ui.toggleSidebar()
  }

  const closeMenu = () => {
    ui.sidebarOpen = false
  }
</script>

<template>
  <header class="border-b-[3px] border-fg bg-surface">
    <div class="mx-auto flex w-full max-w-6xl items-center justify-between gap-3 px-6 py-4">
      <RouterLink
        to="/"
        class="brutal-press inline-flex min-h-11 items-center gap-3 border-[3px] border-fg bg-vivid-yellow px-4 text-black shadow-brutal-sm"
        @click="closeMenu"
      >
        <svg viewBox="0 0 24 24" aria-hidden="true" class="size-6 shrink-0">
          <rect x="1.5" y="1.5" width="21" height="21" fill="none" stroke="currentColor" stroke-width="3" />
          <path d="M13 4 6.5 13H11l-1 7 6.5-9H12z" fill="currentColor" />
        </svg>
        <span class="text-lg font-bold tracking-tight uppercase">App</span>
      </RouterLink>

      <div class="flex items-center gap-2">
        <nav class="hidden items-center gap-2 md:flex" aria-label="Main">
          <RouterLink v-for="link in links" :key="link.name" :to="link.to" :class="[linkClasses, linkTone(link.name)]">
            {{ link.label }}
          </RouterLink>
        </nav>

        <button
          type="button"
          class="brutal-press inline-flex size-11 items-center justify-center border-[3px] border-fg bg-vivid-yellow text-black shadow-brutal-sm"
          :aria-label="isDark ? 'Switch to light theme' : 'Switch to dark theme'"
          :aria-pressed="isDark"
          @click="toggleTheme"
        >
          <svg v-if="isDark" viewBox="0 0 24 24" aria-hidden="true" class="size-6">
            <circle cx="12" cy="12" r="4.5" fill="none" stroke="currentColor" stroke-width="2.5" />
            <path
              d="M12 2v2.5M12 19.5V22M2 12h2.5M19.5 12H22M4.9 4.9l1.8 1.8M17.3 17.3l1.8 1.8M19.1 4.9l-1.8 1.8M6.7 17.3l-1.8 1.8"
              stroke="currentColor"
              stroke-width="2.5"
              stroke-linecap="square"
            />
          </svg>
          <svg v-else viewBox="0 0 24 24" aria-hidden="true" class="size-6">
            <path
              d="M20.5 14.5A8.5 8.5 0 0 1 9.5 3.5a8.5 8.5 0 1 0 11 11z"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              stroke-linejoin="round"
            />
          </svg>
        </button>

        <button
          type="button"
          class="brutal-press inline-flex size-11 items-center justify-center border-[3px] border-fg bg-mint text-black shadow-brutal-sm md:hidden"
          :aria-label="ui.sidebarOpen ? 'Close menu' : 'Open menu'"
          :aria-expanded="ui.sidebarOpen"
          aria-controls="mobile-navigation"
          @click="toggleMenu"
        >
          <svg v-if="ui.sidebarOpen" viewBox="0 0 24 24" aria-hidden="true" class="size-6">
            <path d="m5 5 14 14M19 5 5 19" stroke="currentColor" stroke-width="3" />
          </svg>
          <svg v-else viewBox="0 0 24 24" aria-hidden="true" class="size-6">
            <path d="M3 6h18M3 12h18M3 18h18" stroke="currentColor" stroke-width="3" />
          </svg>
        </button>
      </div>
    </div>

    <nav
      v-if="ui.sidebarOpen"
      id="mobile-navigation"
      class="border-t-[3px] border-fg bg-muted-bg md:hidden"
      aria-label="Mobile"
      @keydown.esc="closeMenu"
    >
      <div class="mx-auto flex w-full max-w-6xl flex-col gap-2 px-6 py-4">
        <RouterLink
          v-for="link in links"
          :key="link.name"
          :to="link.to"
          :class="[linkClasses, linkTone(link.name)]"
          @click="closeMenu"
        >
          {{ link.label }}
        </RouterLink>
      </div>
    </nav>
  </header>
</template>
