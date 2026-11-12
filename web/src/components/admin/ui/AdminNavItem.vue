<script setup lang="ts">
  import { RouterLink } from 'vue-router'

  defineProps<{
    to: string
    label: string
    badge?: string | number
    collapsed?: boolean
  }>()
</script>

<template>
  <RouterLink
    :to="to"
    :title="collapsed ? label : undefined"
    class="group/item flex items-center gap-3 rounded-xl px-3.5 py-2.5 text-sm font-medium text-[var(--admin-sidebar-fg)] transition-colors hover:bg-[var(--admin-sidebar-active-bg)] hover:text-[var(--admin-sidebar-fg-hover)]"
    :class="collapsed ? 'lg:justify-center lg:px-2' : 'justify-between'"
    active-class="bg-[var(--admin-sidebar-active-bg)] text-[var(--admin-sidebar-active-fg)] font-semibold shadow-sm"
  >
    <div class="flex items-center gap-3">
      <slot name="icon" />
      <span :class="collapsed ? 'lg:hidden' : ''">{{ label }}</span>
    </div>

    <span
      v-if="badge !== undefined"
      class="rounded-full bg-blue-500/10 px-2.5 py-0.5 text-xs font-semibold text-blue-400"
      :class="collapsed ? 'lg:hidden' : ''"
    >
      {{ badge }}
    </span>
  </RouterLink>
</template>
