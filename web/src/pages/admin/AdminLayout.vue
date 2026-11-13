<script setup lang="ts">
  import { ref } from 'vue'
  import { RouterView } from 'vue-router'
  import AdminSidebar from '../../components/admin/layout/AdminSidebar.vue'
  import AdminHeader from '../../components/admin/layout/AdminHeader.vue'
  import ConfirmModal from '../../components/admin/ui/ConfirmModal.vue'
  import { useUIStore } from '../../store/ui'

  const ui = useUIStore()

  const sidebarOpen = ref(window.innerWidth >= 1024)

  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
  }

  const closeSidebar = () => {
    if (window.innerWidth < 1024) {
      sidebarOpen.value = false
    }
  }
</script>

<template>
  <div class="admin-theme flex min-h-dvh bg-[var(--admin-bg)] text-[var(--admin-fg)] antialiased">
    <!-- Sidebar -->
    <AdminSidebar :is-open="sidebarOpen" @close="closeSidebar" />

    <!-- Main Wrapper -->
    <div class="flex min-w-0 flex-1 flex-col">
      <AdminHeader :sidebar-open="sidebarOpen" @toggle-sidebar="toggleSidebar" />

      <main class="flex-1 p-6 lg:p-8">
        <div class="mx-auto max-w-7xl">
          <RouterView />
        </div>
      </main>
    </div>

    <!-- Toast Notification (Top Right) -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="transform translate-y-2 opacity-0"
      enter-to-class="transform translate-y-0 opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="transform translate-y-0 opacity-100"
      leave-to-class="transform translate-y-2 opacity-0"
    >
      <div
        v-if="ui.toast"
        class="fixed top-4 right-4 z-50 flex items-center gap-3 rounded-xl border px-4 py-3 text-xs font-semibold shadow-xl backdrop-blur-md"
        :class="ui.toast.type === 'error'
          ? 'border-red-500/30 bg-red-500/90 text-white'
          : 'border-emerald-500/30 bg-emerald-500/90 text-white'"
      >
        <span>{{ ui.toast.message }}</span>
        <button type="button" class="ml-2 opacity-80 hover:opacity-100" @click="ui.toast = null">✕</button>
      </div>
    </Transition>

    <!-- Global Confirm Modal -->
    <ConfirmModal />
  </div>
</template>
