<script setup lang="ts">
  import AdminStatCard from '../../components/admin/ui/AdminStatCard.vue'
</script>

<template>
  <div class="flex flex-col gap-8">
    <!-- Welcome Banner -->
    <div
      class="flex flex-col justify-between gap-4 rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm sm:flex-row sm:items-center"
    >
      <div class="flex flex-col gap-1">
        <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Dashboard Overview</h1>
        <p class="text-sm text-[var(--admin-muted)]">
          Welcome back, Administrator. Here's what's happening across your platform today.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          type="button"
          class="inline-flex min-h-10 items-center justify-center rounded-xl bg-blue-600 px-4 text-xs font-semibold text-white shadow-sm transition-colors hover:bg-blue-700"
        >
          Download Report
        </button>
      </div>
    </div>

    <!-- Stat Cards Grid -->
    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
      <AdminStatCard title="Total Users" value="12,845" change="+12.3%" trend="up">
        <template #icon>
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
            />
          </svg>
        </template>
      </AdminStatCard>

      <AdminStatCard title="Active Sessions" value="1,429" change="+4.1%" trend="up">
        <template #icon>
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
            />
          </svg>
        </template>
      </AdminStatCard>

      <AdminStatCard title="System Load" value="28.4%" change="-2.5%" trend="down">
        <template #icon>
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </template>
      </AdminStatCard>

      <AdminStatCard title="Error Rate" value="0.02%" change="0.0%" trend="neutral">
        <template #icon>
          <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            />
          </svg>
        </template>
      </AdminStatCard>
    </div>

    <!-- Main Content Section / Recent Activity Table Mock -->
    <div class="rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm">
      <div class="flex items-center justify-between border-b border-[var(--admin-border)] pb-6">
        <div>
          <h2 class="text-base font-bold tracking-tight">Recent Activity</h2>
          <p class="text-xs text-[var(--admin-muted)]">Latest system access events and user actions</p>
        </div>

        <button type="button" class="text-xs font-semibold text-blue-500 hover:text-blue-600">View all →</button>
      </div>

      <div class="mt-4 divide-y divide-[var(--admin-border)]">
        <div class="flex items-center justify-between py-3.5">
          <div class="flex items-center gap-3">
            <div
              class="flex size-9 items-center justify-center rounded-xl bg-emerald-500/10 text-xs font-bold text-emerald-500"
            >
              IN
            </div>
            <div>
              <p class="text-xs font-semibold">User signed in successfully</p>
              <p class="text-[11px] text-[var(--admin-muted)]">admin@example.com • 2 mins ago</p>
            </div>
          </div>
          <span
            class="rounded-full bg-emerald-500/10 px-2.5 py-1 text-[10px] font-semibold text-emerald-600 dark:text-emerald-400"
            >Success</span
          >
        </div>

        <div class="flex items-center justify-between py-3.5">
          <div class="flex items-center gap-3">
            <div
              class="flex size-9 items-center justify-center rounded-xl bg-blue-500/10 text-xs font-bold text-blue-500"
            >
              SEC
            </div>
            <div>
              <p class="text-xs font-semibold">Security token refreshed</p>
              <p class="text-[11px] text-[var(--admin-muted)]">System automated task • 15 mins ago</p>
            </div>
          </div>
          <span
            class="rounded-full bg-blue-500/10 px-2.5 py-1 text-[10px] font-semibold text-blue-600 dark:text-blue-400"
            >Info</span
          >
        </div>

        <div class="flex items-center justify-between py-3.5">
          <div class="flex items-center gap-3">
            <div
              class="flex size-9 items-center justify-center rounded-xl bg-amber-500/10 text-xs font-bold text-amber-500"
            >
              CFG
            </div>
            <div>
              <p class="text-xs font-semibold">Environment configuration updated</p>
              <p class="text-[11px] text-[var(--admin-muted)]">admin@example.com • 1 hour ago</p>
            </div>
          </div>
          <span
            class="rounded-full bg-amber-500/10 px-2.5 py-1 text-[10px] font-semibold text-amber-600 dark:text-amber-400"
            >Warning</span
          >
        </div>
      </div>
    </div>
  </div>
</template>
