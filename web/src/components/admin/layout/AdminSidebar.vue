<script setup lang="ts">
  import { RouterLink } from 'vue-router'
  import AdminNavItem from '../ui/AdminNavItem.vue'
  import AdminNavGroup from '../ui/AdminNavGroup.vue'

  defineProps<{
    isOpen?: boolean
  }>()

  defineEmits<{
    close: []
  }>()
</script>

<template>
  <!-- Mobile Backdrop -->
  <div
    v-if="isOpen"
    class="fixed inset-0 z-40 bg-slate-950/60 backdrop-blur-xs lg:hidden"
    @click="$emit('close')"
  />

  <!-- Sidebar Container -->
  <aside
    class="group fixed inset-y-0 left-0 z-50 flex flex-col border-r border-[var(--admin-border)] bg-[var(--admin-sidebar-bg)] transition-all duration-300 lg:sticky lg:top-0 lg:h-dvh lg:self-start"
    :class="[
      isOpen ? 'translate-x-0 w-64' : '-translate-x-full lg:translate-x-0 lg:w-20',
      { 'collapsed': !isOpen }
    ]"
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
          <span class="text-sm font-bold tracking-tight text-[var(--admin-fg)]">Admin Control</span>
          <span class="text-[10px] font-medium text-[var(--admin-muted)]">Enterprise Suite</span>
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
    <div class="flex-1 overflow-y-auto px-3 py-6 lg:px-4 whitespace-nowrap lg:group-[.collapsed]:overflow-visible">
      <div class="space-y-6">
        <!-- Main Section -->
        <div>
          <p class="px-3 text-[11px] font-semibold tracking-wider text-[var(--admin-muted)] uppercase lg:group-[.collapsed]:hidden">
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
          <p class="px-3 text-[11px] font-semibold tracking-wider text-[var(--admin-muted)] uppercase lg:group-[.collapsed]:hidden">
            Management
          </p>
          <div class="mt-2 space-y-1">
            <AdminNavItem to="/admin" label="Analytics" badge="Live" :collapsed="!isOpen">
              <template #icon>
                <svg class="size-4.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                  />
                </svg>
              </template>
            </AdminNavItem>

            <AdminNavGroup
              label="System Users"
              :items="[
                { key: 'users', label: 'Users', to: { name: 'admin', query: { tab: 'users' } } },
                { key: 'roles', label: 'Roles', to: { name: 'admin', query: { tab: 'roles' } } },
                { key: 'permissions', label: 'Permissions', to: { name: 'admin', query: { tab: 'permissions' } } },
              ]"
              :collapsed="!isOpen"
            >
              <template #icon>
                <svg class="size-4.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
                  />
                </svg>
              </template>
            </AdminNavGroup>
          </div>
        </div>
      </div>
    </div>

    <!-- User Profile Footer -->
    <div class="shrink-0 border-t border-[var(--admin-border)] p-3 lg:p-4">
      <div
        class="flex items-center gap-3 rounded-xl bg-[var(--admin-border)]/30 p-2.5 whitespace-nowrap"
        :class="!isOpen ? 'lg:justify-center' : ''"
      >
        <div class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-blue-600 text-xs font-bold text-white">
          AD
        </div>
        <div class="flex flex-1 flex-col truncate lg:group-[.collapsed]:hidden">
          <span class="truncate text-xs font-semibold text-[var(--admin-fg)]">Administrator</span>
          <span class="truncate text-[10px] text-[var(--admin-muted)]">admin@system.local</span>
        </div>
      </div>
    </div>
  </aside>
</template>
