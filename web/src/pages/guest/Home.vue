<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { RouterLink } from 'vue-router'

  import { getHealth } from '../../api/health'
  import type { HealthData } from '../../api/health'
  import { img } from '../../assets/img'
  import { useCountStore, useToggleStore } from '../../store/data'
  import { useUIStore } from '../../store/ui'

  const count = useCountStore()
  const toggle = useToggleStore()
  const ui = useUIStore()

  const health = ref<HealthData | null>(null)
  const healthError = ref<string | null>(null)
  const healthLoading = ref(false)
  const healthTone = computed(() => (healthError.value ? 'text-rose-400' : 'text-slate-300'))

  onMounted(async () => {
    healthLoading.value = true
    try {
      health.value = await getHealth()
    } catch (error) {
      healthError.value = error instanceof Error ? error.message : 'unknown error'
    } finally {
      healthLoading.value = false
    }
  })
</script>

<template>
  <main class="min-h-screen bg-slate-950 text-slate-100">
    <div class="mx-auto flex max-w-5xl flex-col gap-10 px-6 py-12">
      <header class="flex flex-wrap items-center justify-between gap-4">
        <RouterLink to="/" class="flex items-center gap-3 text-lg font-semibold">
          <img :src="img.vite" alt="" class="h-8 w-8" />
          <span>App</span>
        </RouterLink>
        <nav class="flex items-center gap-4 text-sm">
          <RouterLink to="/about" class="text-slate-300 transition hover:text-white"> About </RouterLink>
          <button
            type="button"
            class="rounded-md border border-slate-700 px-3 py-1.5 text-slate-200 transition hover:border-teal-400 hover:text-white"
            @click="ui.toggleSidebar()"
          >
            Sidebar: {{ ui.sidebarOpen ? 'open' : 'closed' }}
          </button>
        </nav>
      </header>

      <section
        class="overflow-hidden rounded-2xl border border-slate-800 bg-slate-900/60 shadow-2xl shadow-teal-500/10"
      >
        <img :src="img.hero" alt="" class="h-48 w-full object-cover" />
        <div class="grid gap-6 p-8 md:grid-cols-2">
          <div class="space-y-3">
            <p class="text-xs font-semibold tracking-[0.2em] text-teal-400 uppercase">Vue + Tailwind + Pinia + Axios</p>
            <h1 class="text-3xl font-bold tracking-tight">Frontend scaffolding</h1>
            <p class="text-sm leading-relaxed text-slate-400">
              Router, Pinia stores, axios client, and TanStack Query are wired up. Replace this page with the real
              application shell.
            </p>
          </div>

          <div class="space-y-4">
            <div class="rounded-xl border border-slate-800 bg-slate-950/60 p-4">
              <p class="text-xs tracking-wider text-slate-500 uppercase">Pinia counter</p>
              <div class="mt-2 flex items-center gap-3">
                <button
                  type="button"
                  class="h-8 w-8 rounded-md bg-slate-800 text-lg leading-none transition hover:bg-slate-700"
                  @click="count.decrement()"
                >
                  -
                </button>
                <span class="min-w-10 text-center text-xl font-semibold">{{ count.count }}</span>
                <button
                  type="button"
                  class="h-8 w-8 rounded-md bg-teal-500 text-lg leading-none text-slate-950 transition hover:bg-teal-400"
                  @click="count.increment()"
                >
                  +
                </button>
                <span class="ml-auto text-xs text-slate-500">double: {{ count.double }}</span>
              </div>
            </div>

            <div class="flex items-center justify-between rounded-xl border border-slate-800 bg-slate-950/60 p-4">
              <span class="text-sm text-slate-400">Toggle</span>
              <button
                type="button"
                class="rounded-full px-3 py-1 text-xs font-semibold transition"
                :class="toggle.active ? 'bg-teal-500 text-slate-950' : 'bg-slate-800 text-slate-300'"
                @click="toggle.toggle()"
              >
                {{ toggle.active ? 'active' : 'inactive' }}
              </button>
            </div>

            <div class="rounded-xl border border-slate-800 bg-slate-950/60 p-4 text-sm">
              <p class="text-xs tracking-wider text-slate-500 uppercase">API health</p>
              <p class="mt-2" :class="healthTone">
                <span v-if="healthLoading">checking...</span>
                <span v-else-if="healthError">{{ healthError }}</span>
                <span v-else>{{ health?.status ?? 'unknown' }}</span>
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  </main>
</template>
