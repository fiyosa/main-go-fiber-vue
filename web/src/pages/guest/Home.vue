<script setup lang="ts">
  import { computed, onMounted, ref } from 'vue'
  import { RouterLink } from 'vue-router'

  import { getHealth } from '../../api/health'
  import type { HealthData } from '../../api/health'
  import { BrutalBadge, BrutalButton, BrutalCard } from './components/ui'
  import { useCountStore, useToggleStore } from '../../store/data'

  const count = useCountStore()
  const toggle = useToggleStore()

  const health = ref<HealthData | null>(null)
  const healthError = ref<string | null>(null)
  const healthLoading = ref(false)

  // Single source of truth for the health card so the badge and the message
  // never disagree.
  const healthState = computed<'checking' | 'error' | 'ok'>(() => {
    if (healthLoading.value) return 'checking'
    if (healthError.value) return 'error'
    return 'ok'
  })

  const healthLabel = computed(() => {
    if (healthState.value === 'checking') return 'Checking'
    if (healthState.value === 'error') return 'Offline'
    return health.value?.status ?? 'Unknown'
  })

  const healthTone = computed(() => {
    if (healthState.value === 'error') return 'hot-red'
    if (healthState.value === 'checking') return 'soft-violet'
    return 'mint'
  })

  const stats = [
    { label: 'Router', value: 'Vue 3', tone: 'mint' },
    { label: 'State', value: 'Pinia', tone: 'soft-violet' },
    { label: 'Styling', value: 'Tailwind v4', tone: 'sky' },
    { label: 'HTTP', value: 'Axios', tone: 'vivid-yellow' },
  ] as const

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
  <div class="flex flex-col gap-8">
    <section class="grid items-center gap-6 lg:grid-cols-[1.2fr_0.8fr]">
      <div class="flex flex-col items-start gap-4">
        <BrutalBadge tone="hot-red" :rotation="-2">Neo-Brutalism UI</BrutalBadge>

        <h1 class="text-4xl leading-none font-bold tracking-tight sm:text-5xl">
          Frontend scaffold with
          <span class="bg-vivid-yellow px-2 text-black">hard edges</span>
        </h1>

        <p class="max-w-xl text-base leading-relaxed">
          Router, Pinia stores, an axios client and a brutalist component set. Every theme, font and asset is served
          from this project — no CDN, no network fetch.
        </p>

        <div class="flex flex-wrap items-center gap-3">
          <RouterLink
            to="/about"
            class="brutal-press inline-flex min-h-12 items-center justify-center border-[3px] border-fg bg-vivid-yellow px-6 font-mono text-sm font-bold tracking-widest text-black uppercase shadow-brutal-sm"
          >
            Explore layout
          </RouterLink>
          <a
            href="#demos"
            class="brutal-press inline-flex min-h-12 items-center justify-center border-[3px] border-fg bg-surface px-6 font-mono text-sm font-bold tracking-widest text-fg uppercase shadow-brutal-sm"
          >
            See demos
          </a>
        </div>
      </div>

      <BrutalCard accent="vivid-yellow" tone="muted" class="rotate-1">
        <template #header>Quick start</template>

        <ol class="flex flex-col gap-3 font-mono text-sm">
          <li class="flex gap-3">
            <span class="border-[3px] border-fg bg-mint px-2 font-bold text-black">1</span>
            <span>Install dependencies</span>
          </li>
          <li class="flex gap-3">
            <span class="border-[3px] border-fg bg-soft-violet px-2 font-bold text-black">2</span>
            <span>Run the Go server</span>
          </li>
          <li class="flex gap-3">
            <span class="border-[3px] border-fg bg-sky px-2 font-bold text-black">3</span>
            <span>Open the SPA</span>
          </li>
        </ol>
      </BrutalCard>
    </section>

    <section class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4" aria-label="Stack">
      <BrutalCard v-for="stat in stats" :key="stat.label" :accent="stat.tone">
        <template #header>{{ stat.label }}</template>
        <p class="text-2xl font-bold">{{ stat.value }}</p>
      </BrutalCard>
    </section>

    <section id="demos" class="grid gap-6 lg:grid-cols-3" aria-label="Interactive demos">
      <BrutalCard accent="mint" class="lg:col-span-1">
        <template #header>Pinia counter</template>

        <div class="flex items-center gap-3">
          <BrutalButton variant="ghost" size="sm" aria-label="Decrease count" @click="count.decrement()">
            <svg viewBox="0 0 24 24" aria-hidden="true" class="size-5">
              <path d="M4 12h16" stroke="currentColor" stroke-width="3" />
            </svg>
          </BrutalButton>

          <span class="min-w-16 border-[3px] border-fg bg-surface px-3 py-2 text-center text-2xl font-bold">
            {{ count.count }}
          </span>

          <BrutalButton size="sm" aria-label="Increase count" @click="count.increment()">
            <svg viewBox="0 0 24 24" aria-hidden="true" class="size-5">
              <path d="M12 4v16M4 12h16" stroke="currentColor" stroke-width="3" />
            </svg>
          </BrutalButton>
        </div>

        <p class="mt-4 font-mono text-xs tracking-widest uppercase">Double: {{ count.double }}</p>
      </BrutalCard>

      <BrutalCard accent="soft-violet">
        <template #header>Toggle</template>

        <div class="flex items-center justify-between gap-3">
          <span class="text-sm">Store state</span>
          <BrutalButton
            size="sm"
            :variant="toggle.active ? 'primary' : 'ghost'"
            :aria-pressed="toggle.active"
            @click="toggle.toggle()"
          >
            {{ toggle.active ? 'Active' : 'Inactive' }}
          </BrutalButton>
        </div>
      </BrutalCard>

      <BrutalCard :accent="healthTone">
        <template #header>API health</template>

        <div class="flex flex-col gap-3">
          <BrutalBadge :tone="healthTone">{{ healthLabel }}</BrutalBadge>

          <p class="font-mono text-xs break-all">
            <span v-if="healthState === 'checking'">Waiting for response...</span>
            <span v-else-if="healthState === 'error'">{{ healthError }}</span>
            <span v-else>Endpoint responded successfully.</span>
          </p>
        </div>
      </BrutalCard>
    </section>
  </div>
</template>
