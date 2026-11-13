<script setup lang="ts">
  import { ref, onMounted, onUnmounted, computed } from 'vue'
  import { RouterLink, useRoute, useRouter } from 'vue-router'
  import { useUIStore } from '../../../store/ui'
  import { logout, getMe, type UserMe } from '../../../api/auth'

  defineProps<{
    sidebarOpen?: boolean
  }>()

  defineEmits<{
    toggleSidebar: []
  }>()

  const ui = useUIStore()
  const isDark = computed(() => ui.theme === 'dark')
  const route = useRoute()
  const router = useRouter()

  const currentUser = ref<UserMe | null>(null)

  onMounted(async () => {
    try {
      currentUser.value = await getMe()
    } catch {
      currentUser.value = { id: 1, email: 'admin@system.local', roles: [] }
    }
  })

  const userEmail = computed(() => currentUser.value?.email || 'admin@system.local')
  const userInitial = computed(() => {
    const email = userEmail.value
    return email ? email.charAt(0).toUpperCase() : 'A'
  })

  const breadcrumbs = computed(() => {
    const name = route.name as string
    if (name === 'admin') {
      return ['Dashboard']
    }
    if (name === 'admin-register') {
      return ['Auth', 'Register']
    }
    if (name === 'admin-cron') {
      return ['Worker', 'Cron']
    }
    if (name === 'admin-cron-create') {
      return ['Worker', 'Cron', 'Create']
    }
    if (name === 'admin-cron-edit') {
      return ['Worker', 'Cron', 'Edit']
    }
    // Generic fallback
    const segments = route.path.split('/').filter(Boolean)
    const result: string[] = []
    for (const seg of segments) {
      if (seg === 'admin' || seg === 'work' || seg === 'auth') continue
      if (/^\d+$/.test(seg)) continue
      const capitalized = seg.charAt(0).toUpperCase() + seg.slice(1)
      result.push(capitalized)
    }
    return result.length > 0 ? result : ['Dashboard']
  })

  const toggleTheme = () => {
    ui.toggleTheme()
  }

  const handleLogout = async () => {
    const ok = await ui.confirm({
      title: 'Sign Out',
      message: 'Are you sure you want to sign out?',
      confirmText: 'Sign Out',
      cancelText: 'Cancel',
      type: 'info',
    })
    if (!ok) return

    try {
      await logout()
    } finally {
      dropdownOpen.value = false
      router.push('/login')
    }
  }

  const dropdownOpen = ref(false)
  const dropdownRef = ref<HTMLElement | null>(null)

  const handleClickOutside = (e: MouseEvent) => {
    if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
      dropdownOpen.value = false
    }
  }

  onMounted(() => document.addEventListener('click', handleClickOutside))
  onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>

<template>
  <header
    class="sticky top-0 z-30 flex h-16 w-full items-center justify-between border-b border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 shadow-xs sm:px-6"
  >
    <!-- Left Area: Toggle Button + Breadcrumb/Title -->
    <div class="flex items-center gap-3">
      <button
        type="button"
        class="inline-flex size-9 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-bg)] text-[var(--admin-muted)] shadow-xs transition-colors hover:bg-[var(--admin-border)]/50 hover:text-[var(--admin-fg)]"
        :aria-label="sidebarOpen ? 'Hide Sidebar' : 'Show Sidebar'"
        :title="sidebarOpen ? 'Hide Menu' : 'Show Menu'"
        @click="$emit('toggleSidebar')"
      >
        <svg class="size-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>

      <nav class="hidden items-center gap-2 text-sm font-medium sm:flex" aria-label="Breadcrumb">
        <template v-for="(crumb, index) in breadcrumbs" :key="index">
          <span
            :class="index === breadcrumbs.length - 1
              ? 'font-semibold text-[var(--admin-fg)]'
              : 'text-[var(--admin-muted)]'"
          >{{ crumb }}</span>
          <span v-if="index < breadcrumbs.length - 1" class="text-[var(--admin-muted)]">/</span>
        </template>
      </nav>
    </div>

    <!-- Right Area: Actions -->
    <div class="flex items-center gap-2 sm:gap-3">
      <!-- Search Input Mock -->
      <div class="relative hidden sm:block">
        <svg
          class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-[var(--admin-muted)]"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          type="text"
          placeholder="Search metrics or data..."
          class="h-9 w-48 rounded-xl border border-[var(--admin-border)] bg-[var(--admin-bg)] pr-3 pl-9 text-xs text-[var(--admin-fg)] placeholder:text-[var(--admin-muted)] focus:border-blue-500 focus:outline-hidden md:w-64"
        />
      </div>

      <!-- Notifications -->
      <button
        type="button"
        class="relative inline-flex size-9 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-bg)] text-[var(--admin-muted)] hover:text-[var(--admin-fg)]"
        aria-label="Notifications"
      >
        <svg class="size-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
          />
        </svg>
        <span class="absolute top-1.5 right-1.5 size-2 rounded-full bg-blue-500" />
      </button>

      <!-- Theme Toggle -->
      <button
        type="button"
        class="inline-flex size-9 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-bg)] text-[var(--admin-muted)] hover:text-[var(--admin-fg)]"
        :aria-label="isDark ? 'Switch to light theme' : 'Switch to dark theme'"
        @click="toggleTheme"
      >
        <svg v-if="isDark" class="size-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
          />
        </svg>
        <svg v-else class="size-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
          />
        </svg>
      </button>

      <!-- Profile Chip with Dropdown -->
      <div class="relative border-l border-[var(--admin-border)] pl-3" ref="dropdownRef">
        <button
          type="button"
          class="flex items-center gap-2 focus:outline-hidden"
          @click="dropdownOpen = !dropdownOpen"
          aria-label="User profile menu"
        >
          <div
            class="flex size-8 items-center justify-center rounded-full bg-indigo-600 text-xs font-bold text-white shadow-xs"
          >
            {{ userInitial }}
          </div>
          <svg
            class="hidden size-3.5 text-[var(--admin-muted)] md:inline-block"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
          </svg>
        </button>

        <!-- Dropdown Menu -->
        <div
          v-if="dropdownOpen"
          class="absolute right-0 z-50 mt-2 w-48 rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] py-1.5 text-xs shadow-lg"
        >
          <div class="border-b border-[var(--admin-border)] px-3 py-2">
            <p class="truncate text-xs font-semibold text-[var(--admin-fg)]">{{ userEmail }}</p>
          </div>
          <RouterLink
            to="/"
            class="flex items-center gap-2 px-3 py-2 text-[var(--admin-fg)] transition-colors hover:bg-[var(--admin-border)]/50"
            @click="dropdownOpen = false"
          >
            <svg
              class="size-4 text-[var(--admin-muted)]"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
              />
            </svg>
            <span>Go to Root (/)</span>
          </RouterLink>
          <button
            type="button"
            class="flex w-full items-center gap-2 px-3 py-2 text-left text-rose-500 transition-colors hover:bg-rose-500/10"
            @click="handleLogout"
          >
            <svg class="size-4 text-rose-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
              />
            </svg>
            <span>Logout</span>
          </button>
        </div>
      </div>
    </div>
  </header>
</template>
