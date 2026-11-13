<script setup lang="ts">
  import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
  import { RouterLink, useRoute, type RouteLocationRaw } from 'vue-router'

  export interface NavSubItem {
    key: string
    label: string
    to: RouteLocationRaw
  }

  const props = defineProps<{
    label: string
    items: NavSubItem[]
    collapsed?: boolean
    badge?: string | number
  }>()

  const route = useRoute()
  const isOpen = ref(false)
  const flyoutOpen = ref(false)
  const groupRef = ref<HTMLElement | null>(null)

  // RouterLink active state ignores query strings, so the active child is
  // resolved manually; placeholder tabs share the same route path.
  const isItemActive = (item: NavSubItem): boolean => {
    const target = item.to
    if (typeof target === 'object' && target !== null) {
      const name = 'name' in target ? (target.name as string | undefined) : undefined
      if (name !== undefined && route.name !== name) return false

      const query = 'query' in target ? (target.query as Record<string, unknown> | undefined) : undefined
      if (query) {
        return Object.entries(query).every(([key, value]) => route.query[key] === String(value))
      }

      const path = 'path' in target ? (target.path as string | undefined) : undefined
      return path === undefined || route.path === path
    }

    return route.path === String(target)
  }

  const childClass = (item: NavSubItem) => [
    'flex items-center justify-between rounded-lg px-3 py-2 text-xs font-medium transition-colors',
    isItemActive(item)
      ? 'bg-[var(--admin-sidebar-active-bg)] text-[var(--admin-sidebar-active-fg)] font-semibold shadow-xs'
      : 'text-[var(--admin-sidebar-fg)] hover:bg-[var(--admin-sidebar-active-bg)] hover:text-[var(--admin-sidebar-fg-hover)]',
  ]

  const hasActiveChild = computed(() => props.items.some((item) => isItemActive(item)))

  // Auto-open the accordion when one of the children is the active route.
  watch(
    hasActiveChild,
    (active) => {
      if (active) isOpen.value = true
    },
    { immediate: true }
  )

  const toggle = () => {
    if (props.collapsed) {
      flyoutOpen.value = !flyoutOpen.value
      return
    }
    isOpen.value = !isOpen.value
  }

  const handleClickOutside = (event: MouseEvent) => {
    if (groupRef.value && !groupRef.value.contains(event.target as Node)) {
      flyoutOpen.value = false
    }
  }

  onMounted(() => document.addEventListener('click', handleClickOutside))
  onUnmounted(() => document.removeEventListener('click', handleClickOutside))
</script>

<template>
  <div ref="groupRef" class="relative">
    <!-- Parent toggle: dropdown on wide sidebar, flyout on the collapsed rail. -->
    <button
      type="button"
      :aria-expanded="collapsed ? flyoutOpen : isOpen"
      :title="collapsed ? label : undefined"
      class="group/item flex w-full items-center gap-3 rounded-xl px-3.5 py-2.5 text-sm font-medium transition-colors hover:bg-[var(--admin-sidebar-active-bg)] hover:text-[var(--admin-sidebar-fg-hover)]"
      :class="[
        collapsed ? 'lg:justify-center lg:px-2' : 'justify-between',
        hasActiveChild
          ? 'bg-[var(--admin-sidebar-active-bg)]/50 font-semibold text-[var(--admin-sidebar-active-fg)]'
          : 'text-[var(--admin-sidebar-fg)]',
      ]"
      @click="toggle"
    >
      <div class="flex min-w-0 items-center gap-3">
        <slot name="icon" />
        <span class="truncate" :class="collapsed ? 'lg:hidden' : ''">{{ label }}</span>
      </div>

      <div
        class="flex shrink-0 items-center gap-1.5"
        :class="collapsed ? 'lg:hidden' : ''"
      >
        <span
          v-if="badge !== undefined"
          class="rounded-full bg-blue-500/10 px-2.5 py-0.5 text-xs font-semibold text-[var(--admin-sidebar-active-fg)] lg:group-[.collapsed]:hidden"
          :class="collapsed ? 'lg:hidden' : ''"
        >
          {{ badge }}
        </span>
        <svg
          class="size-4 text-[var(--admin-sidebar-fg)] transition-transform duration-200"
          :class="[isOpen && !collapsed ? 'rotate-180' : '', collapsed ? 'lg:hidden' : '']"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </button>

    <!-- Accordion for the wide sidebar and the mobile drawer. -->
    <div
      class="grid transition-all duration-200 ease-in-out"
      :class="[isOpen && !collapsed ? 'mt-1 grid-rows-[1fr]' : 'grid-rows-[0fr]', collapsed ? 'lg:hidden' : '']"
    >
      <div class="overflow-hidden">
        <div class="ml-4 space-y-1 border-l border-[var(--admin-border)] py-1 pl-3">
          <RouterLink v-for="sub in items" :key="sub.key" :to="sub.to" :class="childClass(sub)">
            <span class="truncate">{{ sub.label }}</span>
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Flyout popover for the collapsed mini-rail. -->
    <div
      v-if="collapsed && flyoutOpen"
      class="absolute top-0 left-full z-50 ml-2 hidden w-48 rounded-xl border border-[var(--admin-border)] bg-[var(--admin-sidebar-bg)] p-1.5 shadow-xl lg:block"
    >
      <div class="mb-1 border-b border-[var(--admin-border)] px-3 py-2">
        <p class="truncate text-xs font-bold tracking-wide text-[var(--admin-fg)]">{{ label }}</p>
      </div>
      <div class="space-y-0.5">
        <RouterLink
          v-for="sub in items"
          :key="sub.key"
          :to="sub.to"
          :class="childClass(sub)"
          @click="flyoutOpen = false"
        >
          {{ sub.label }}
        </RouterLink>
      </div>
    </div>
  </div>
</template>
