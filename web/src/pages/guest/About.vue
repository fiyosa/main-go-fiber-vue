<script setup lang="ts">
  import { RouterLink } from 'vue-router'

  import { BrutalBadge, BrutalButton, BrutalCard } from '../../components/ui'

  const folders = [
    { name: 'api', detail: 'Endpoint wrappers, one folder per feature.' },
    { name: 'lib', detail: 'Shared axios client and helpers.' },
    { name: 'store', detail: 'Pinia stores for UI and data state.' },
    { name: 'pages', detail: 'Routed views and route table.' },
    { name: 'config', detail: 'Static settings resolved at build time.' },
  ] as const

  const palette = [
    { name: 'bg', light: '#FFFDF5', dark: '#0F172A', className: 'bg-bg', text: 'text-fg' },
    { name: 'surface', light: '#FFFFFF', dark: '#1B263B', className: 'bg-surface', text: 'text-fg' },
    { name: 'fg', light: '#000000', dark: '#F8FAFC', className: 'bg-fg', text: 'text-bg' },
    { name: 'muted-bg', light: '#F2EEDF', dark: '#27334A', className: 'bg-muted-bg', text: 'text-fg' },
    { name: 'hot-red', light: '#FF6B6B', dark: '#FF6B6B', className: 'bg-hot-red', text: 'text-black' },
    { name: 'vivid-yellow', light: '#FFD93D', dark: '#FFD93D', className: 'bg-vivid-yellow', text: 'text-black' },
    { name: 'soft-violet', light: '#C4B5FD', dark: '#C4B5FD', className: 'bg-soft-violet', text: 'text-black' },
    { name: 'mint', light: '#4ECDC4', dark: '#4ECDC4', className: 'bg-mint', text: 'text-black' },
    { name: 'sky', light: '#2196F3', dark: '#2196F3', className: 'bg-sky', text: 'text-black' },
  ] as const

  const buttonVariants = [
    { label: 'Primary', variant: 'primary' },
    { label: 'Secondary', variant: 'secondary' },
    { label: 'Ghost', variant: 'ghost' },
    { label: 'Danger', variant: 'danger' },
  ] as const

  const badgeTones = ['vivid-yellow', 'hot-red', 'soft-violet', 'mint', 'sky'] as const
</script>

<template>
  <div class="flex flex-col gap-10">
    <header class="flex flex-col items-start gap-3">
      <BrutalBadge tone="soft-violet" :rotation="2">About</BrutalBadge>
      <h1 class="text-4xl leading-none font-bold tracking-tight sm:text-5xl">Design system</h1>
      <p class="max-w-2xl text-base leading-relaxed">
        Hard 3px borders, offset shadows and flat accents. Tokens flip with the theme, accents stay fixed, and every
        font is served from this repository.
      </p>
    </header>

    <section class="flex flex-col gap-4" aria-labelledby="typography-heading">
      <h2 id="typography-heading" class="font-mono text-xs font-bold tracking-[0.2em] uppercase">Typography</h2>

      <BrutalCard accent="sky">
        <template #header>Space Grotesk / Display</template>

        <div class="flex flex-col gap-3">
          <p class="text-4xl leading-none font-bold">Aa Bb Cc 0123</p>
          <p class="text-lg leading-relaxed">The quick brown fox jumps over the lazy dog.</p>
          <p class="font-mono text-xs tracking-widest uppercase">Space Mono / Labels &amp; metadata</p>
        </div>
      </BrutalCard>
    </section>

    <section class="flex flex-col gap-4" aria-labelledby="palette-heading">
      <h2 id="palette-heading" class="font-mono text-xs font-bold tracking-[0.2em] uppercase">Palette</h2>

      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <div
          v-for="color in palette"
          :key="color.name"
          class="border-[3px] border-fg shadow-brutal-sm"
          :class="[color.className, color.text]"
        >
          <p class="border-b-[3px] border-current px-3 py-2 font-mono text-xs font-bold tracking-widest uppercase">
            {{ color.name }}
          </p>
          <dl class="px-3 py-2 font-mono text-[11px] leading-relaxed">
            <div class="flex justify-between gap-2">
              <dt>Light</dt>
              <dd>{{ color.light }}</dd>
            </div>
            <div class="flex justify-between gap-2">
              <dt>Dark</dt>
              <dd>{{ color.dark }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </section>

    <section class="flex flex-col gap-4" aria-labelledby="gallery-heading">
      <h2 id="gallery-heading" class="font-mono text-xs font-bold tracking-[0.2em] uppercase">Primitives</h2>

      <div class="grid gap-4 lg:grid-cols-2">
        <BrutalCard accent="mint">
          <template #header>Buttons</template>

          <div class="flex flex-wrap items-center gap-3">
            <BrutalButton v-for="button in buttonVariants" :key="button.label" :variant="button.variant" size="sm">
              {{ button.label }}
            </BrutalButton>
          </div>

          <div class="mt-4 flex flex-wrap items-center gap-3">
            <BrutalButton disabled>Disabled</BrutalButton>
            <BrutalButton variant="ghost" size="sm">Small</BrutalButton>
          </div>
        </BrutalCard>

        <BrutalCard accent="vivid-yellow">
          <template #header>Badges</template>

          <div class="flex flex-wrap items-center gap-3">
            <BrutalBadge
              v-for="(tone, index) in badgeTones"
              :key="tone"
              :tone="tone"
              :rotation="index % 2 === 0 ? -2 : 2"
            >
              {{ tone }}
            </BrutalBadge>
          </div>
        </BrutalCard>
      </div>
    </section>

    <section class="flex flex-col gap-4" aria-labelledby="structure-heading">
      <h2 id="structure-heading" class="font-mono text-xs font-bold tracking-[0.2em] uppercase">Structure</h2>

      <div class="grid gap-4 md:grid-cols-2">
        <BrutalCard v-for="folder in folders" :key="folder.name" accent="hot-red">
          <template #header>{{ folder.name }}</template>
          <p class="text-sm leading-relaxed">{{ folder.detail }}</p>
        </BrutalCard>
      </div>
    </section>

    <RouterLink
      to="/"
      class="brutal-press inline-flex min-h-12 w-fit items-center justify-center border-[3px] border-fg bg-vivid-yellow px-6 font-mono text-sm font-bold tracking-widest text-black uppercase shadow-brutal-sm"
    >
      Back home
    </RouterLink>
  </div>
</template>
