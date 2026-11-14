<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue'
  import { RouterLink } from 'vue-router'
  import AdminNavItem from '../ui/AdminNavItem.vue'
  import AdminNavGroup from '../ui/AdminNavGroup.vue'
  import { getMe, type UserMe } from '../../../api/auth'

  defineProps<{
    isOpen?: boolean
  }>()

  defineEmits<{
    close: []
  }>()

  const currentUser = ref<UserMe | null>(null)

  onMounted(async () => {
    try {
      currentUser.value = await getMe()
    } catch {
      currentUser.value = { id: 1, email: 'admin@system.local', roles: ['admin'] }
    }
  })

  const userEmail = computed(() => currentUser.value?.email || 'admin@system.local')
  const userRole = computed(() => {
    const roles = currentUser.value?.roles
    if (roles && roles.length > 0) {
      return roles.map((r) => r.charAt(0).toUpperCase() + r.slice(1)).join(', ')
    }
    return 'Administrator'
  })
</script>

<template>
  <!-- Mobile Backdrop -->
  <div v-if="isOpen" class="fixed inset-0 z-40 bg-slate-950/60 backdrop-blur-xs lg:hidden" @click="$emit('close')" />

  <!-- Sidebar Container -->
  <aside
    class="group fixed inset-y-0 left-0 z-50 flex flex-col border-r border-[var(--admin-border)] bg-[var(--admin-sidebar-bg)] transition-all duration-300 lg:sticky lg:top-0 lg:h-dvh lg:self-start"
    :class="[isOpen ? 'w-64 translate-x-0' : '-translate-x-full lg:w-20 lg:translate-x-0', { collapsed: !isOpen }]"
  >
    <!-- Brand Header -->
    <div
      class="flex h-16 shrink-0 items-center border-b border-[var(--admin-border)] px-4 lg:px-6"
      :class="isOpen ? 'justify-between' : 'lg:justify-center'"
    >
      <RouterLink to="/admin" class="flex items-center gap-3">
        <div
          class="flex size-9 shrink-0 items-center justify-center rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-500 text-white shadow-md shadow-blue-500/20"
        >
          <svg class="size-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <rect x="3" y="3" width="7" height="7" rx="1.5" />
            <rect x="14" y="3" width="7" height="7" rx="1.5" />
            <rect x="14" y="14" width="7" height="7" rx="1.5" />
            <rect x="3" y="14" width="7" height="7" rx="1.5" />
          </svg>
        </div>
        <div class="flex flex-col whitespace-nowrap lg:group-[.collapsed]:hidden">
          <span class="truncate text-sm font-bold tracking-tight text-[var(--admin-fg)]">{{ userEmail }}</span>
          <span class="text-[10px] font-medium text-[var(--admin-muted)]">{{ userRole }}</span>
        </div>
      </RouterLink>

      <button
        type="button"
        class="rounded-lg p-1.5 text-[var(--admin-muted)] hover:bg-[var(--admin-border)]/50 hover:text-[var(--admin-fg)] lg:hidden"
        @click="$emit('close')"
      >
        <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Navigation List -->
    <div
      class="no-scrollbar flex-1 overflow-y-auto px-3 py-6 whitespace-nowrap lg:px-4 lg:group-[.collapsed]:overflow-visible"
    >
      <div class="space-y-6">
        <!-- Main Section -->
        <div>
          <p
            class="px-3 text-[11px] font-semibold tracking-wider text-[var(--admin-muted)] uppercase lg:group-[.collapsed]:hidden"
          >
            Overview
          </p>
          <div class="mt-2 space-y-1">
            <AdminNavItem to="/admin" label="Dashboard" :collapsed="!isOpen">
              <template #icon>
                <svg class="size-4.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                  />
                </svg>
              </template>
            </AdminNavItem>
          </div>
        </div>

        <!-- Management Section -->
        <div>
          <p
            class="px-3 text-[11px] font-semibold tracking-wider text-[var(--admin-muted)] uppercase lg:group-[.collapsed]:hidden"
          >
            Management
          </p>
          <div class="mt-2 space-y-1">
            <AdminNavGroup
              label="Auth"
              :items="[
                { key: 'register', label: 'Register', to: { name: 'admin-register' } },
                { key: 'change-password', label: 'Change Password', to: { name: 'admin-change-password' } },
              ]"
              :collapsed="!isOpen"
            >
              <template #icon>
                <svg class="size-4.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                  />
                </svg>
              </template>
            </AdminNavGroup>
          </div>
        </div>

        <!-- Compute Section -->
        <div>
          <p
            class="px-3 text-[11px] font-semibold tracking-wider text-[var(--admin-muted)] uppercase lg:group-[.collapsed]:hidden"
          >
            Compute
          </p>
          <div class="mt-2 space-y-1">
            <AdminNavGroup
              label="Worker"
              :items="[{ key: 'cron', label: 'Cron', to: { name: 'admin-cron' } }]"
              :collapsed="!isOpen"
            >
              <template #icon>
                <svg class="size-4.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </template>
            </AdminNavGroup>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>
