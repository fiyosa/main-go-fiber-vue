<script setup lang="ts">
  import { BrutalBadge, BrutalCard } from './components/ui'
  import { serverContent } from './content/server'

  const secondaryLinkClasses =
    'brutal-press inline-flex min-h-12 items-center justify-center gap-2 border-[3px] border-fg bg-surface px-6 font-mono text-sm font-bold tracking-widest text-fg uppercase shadow-brutal-sm'

  const badgeTones = ['sky', 'mint', 'soft-violet', 'vivid-yellow', 'hot-red', 'surface'] as const
</script>

<template>
  <div class="flex flex-col gap-14">
    <!-- Hero -->
    <section class="grid items-start gap-8 lg:grid-cols-[1.4fr_0.6fr]" aria-labelledby="server-title">
      <div class="flex flex-col items-start gap-4">
        <BrutalBadge tone="mint" :rotation="-2">{{ serverContent.headline }}</BrutalBadge>

        <h1 id="server-title" class="text-4xl leading-none font-bold tracking-tight sm:text-5xl">
          {{ serverContent.title }}
        </h1>

        <p class="max-w-2xl text-base leading-relaxed">
          {{ serverContent.summary }}
        </p>

        <div class="flex flex-wrap items-center gap-3">
          <a href="#specs" :class="secondaryLinkClasses">View Telemetry</a>
        </div>
      </div>

      <BrutalCard accent="mint" tone="muted" class="w-full rotate-1">
        <template #header>{{ serverContent.user }}</template>

        <p class="text-sm leading-relaxed font-bold">
          {{ serverContent.host }}
        </p>
        <p class="text-muted-fg mt-3 font-mono text-xs">Uptime: {{ serverContent.uptime }}</p>
      </BrutalCard>
    </section>

    <!-- Specs Grid -->
    <section id="specs" class="flex scroll-mt-24 flex-col gap-10" aria-labelledby="specs-title">
      <div class="flex flex-wrap items-center gap-3">
        <h2 id="specs-title" class="text-2xl font-bold tracking-tight sm:text-3xl">System Telemetry & Hardware</h2>
        <BrutalBadge tone="vivid-yellow">Armbian Amlogic S905X</BrutalBadge>
      </div>

      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <!-- OS & Kernel -->
        <BrutalCard accent="vivid-yellow" tone="surface">
          <template #header>Operating System</template>
          <div class="flex flex-col gap-2">
            <p class="font-mono text-sm font-bold">{{ serverContent.os }}</p>
            <p class="text-muted-fg text-xs">Kernel: {{ serverContent.kernel }}</p>
            <div class="mt-2 flex flex-wrap gap-2">
              <BrutalBadge tone="surface">dpkg: {{ serverContent.packages }}</BrutalBadge>
              <BrutalBadge tone="surface">Shell: {{ serverContent.shell }}</BrutalBadge>
            </div>
          </div>
        </BrutalCard>

        <!-- Hardware (CPU & GPU) -->
        <BrutalCard accent="sky" tone="surface">
          <template #header>Processor & GPU</template>
          <div class="flex flex-col gap-2">
            <p class="font-mono text-sm font-bold">{{ serverContent.cpu }}</p>
            <p class="text-muted-fg text-xs">{{ serverContent.gpu }}</p>
            <div class="mt-2 flex flex-wrap gap-2">
              <BrutalBadge tone="sky">ARMv8</BrutalBadge>
              <BrutalBadge tone="sky">Mali Integrated</BrutalBadge>
            </div>
          </div>
        </BrutalCard>

        <!-- Memory & Swap -->
        <BrutalCard accent="soft-violet" tone="surface">
          <template #header>Memory & Swap</template>
          <div class="flex flex-col gap-2">
            <p class="font-mono text-sm font-bold">RAM: {{ serverContent.memory }}</p>
            <p class="text-muted-fg font-mono text-xs">Swap: {{ serverContent.swap }}</p>
            <div class="mt-2 flex flex-wrap gap-2">
              <BrutalBadge tone="soft-violet">Active Usage</BrutalBadge>
            </div>
          </div>
        </BrutalCard>

        <!-- Storage Disks -->
        <BrutalCard accent="hot-red" tone="surface">
          <template #header>Storage (Ext4)</template>
          <div class="flex flex-col gap-2">
            <p class="font-mono text-sm font-bold">Root (/): {{ serverContent.diskRoot }}</p>
            <p class="text-muted-fg font-mono text-xs">Log (/var/log): {{ serverContent.diskLog }}</p>
            <div class="mt-2 flex flex-wrap gap-2">
              <BrutalBadge tone="hot-red">ext4 file system</BrutalBadge>
            </div>
          </div>
        </BrutalCard>

        <!-- Network & Display -->
        <BrutalCard accent="mint" tone="surface">
          <template #header>Network & Environment</template>
          <div class="flex flex-col gap-2">
            <p class="font-mono text-sm font-bold">IP: {{ serverContent.localIp }}</p>
            <p class="text-muted-fg font-mono text-xs">Locale: {{ serverContent.locale }}</p>
            <div class="mt-2 flex flex-wrap gap-2">
              <BrutalBadge tone="mint">Display: {{ serverContent.display }}</BrutalBadge>
            </div>
          </div>
        </BrutalCard>

        <!-- Cloudflare Tunnel Deployment -->
        <BrutalCard accent="vivid-yellow" tone="surface">
          <template #header>Public Deployment</template>
          <div class="flex flex-col gap-2">
            <p class="font-mono text-sm font-bold">{{ serverContent.deployment.name }}</p>
            <p class="text-xs leading-relaxed">{{ serverContent.deployment.description }}</p>
            <div class="mt-2 flex flex-wrap gap-2">
              <BrutalBadge tone="vivid-yellow">Secure Edge Tunnel</BrutalBadge>
            </div>
          </div>
        </BrutalCard>
      </div>
    </section>

    <!-- Cloudflare Tunnel Details Section -->
    <section aria-labelledby="tunnel-title">
      <BrutalCard accent="sky" tone="muted">
        <template #header>Cloudflare Tunnel Architecture</template>

        <div class="flex flex-col gap-4">
          <h2 id="tunnel-title" class="text-xl font-bold tracking-tight">How public access is routed to Armbian</h2>
          <ul class="flex flex-col gap-2">
            <li
              v-for="detail in serverContent.deployment.details"
              :key="detail"
              class="flex gap-3 text-sm leading-relaxed"
            >
              <span aria-hidden="true" class="mt-2 size-2 shrink-0 border-2 border-fg bg-sky"></span>
              <span>{{ detail }}</span>
            </li>
          </ul>
        </div>
      </BrutalCard>
    </section>
  </div>
</template>
