<script setup lang="ts">
  import { ref } from 'vue'
  import { RouterView } from 'vue-router'
  import AdminSidebar from '../../components/admin/layout/AdminSidebar.vue'
  import AdminHeader from '../../components/admin/layout/AdminHeader.vue'

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
  <div class="admin-theme flex min-h-screen bg-[var(--admin-bg)] text-[var(--admin-fg)] antialiased">
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
  </div>
</template>
