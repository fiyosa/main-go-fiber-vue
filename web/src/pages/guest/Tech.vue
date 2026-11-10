<script setup lang="ts">
  import { BrutalBadge, BrutalCard } from './components/ui'
  import { techContent } from './content/tech'

  const secondaryLinkClasses =
    'brutal-press inline-flex min-h-12 items-center justify-center gap-2 border-[3px] border-fg bg-surface px-6 font-mono text-sm font-bold tracking-widest text-fg uppercase shadow-brutal-sm'

  const categoryAccents = ['vivid-yellow', 'sky', 'mint', 'soft-violet', 'hot-red'] as const
  const badgeTones = ['sky', 'mint', 'soft-violet', 'vivid-yellow', 'hot-red', 'surface'] as const
</script>

<template>
  <div class="flex flex-col gap-14">
    <!-- Hero -->
    <section class="grid items-start gap-8 lg:grid-cols-[1.4fr_0.6fr]" aria-labelledby="tech-title">
      <div class="flex flex-col items-start gap-4">
        <BrutalBadge tone="sky" :rotation="-2">{{ techContent.headline }}</BrutalBadge>

        <h1 id="tech-title" class="text-4xl leading-none font-bold tracking-tight sm:text-5xl">
          {{ techContent.title }}
        </h1>

        <p class="max-w-2xl text-base leading-relaxed">
          {{ techContent.summary }}
        </p>

        <div class="flex flex-wrap items-center gap-3">
          <a href="#categories" :class="secondaryLinkClasses">Explore Stack</a>
        </div>
      </div>

      <BrutalCard accent="sky" tone="muted" class="w-full rotate-1">
        <template #header>Project Overview</template>

        <p class="text-sm leading-relaxed font-bold">
          Fullstack Go & Vue Monorepo Architecture with Neo-Brutalist Design, PostgreSQL, Fiber v3, and MCP Integration.
        </p>
        <p class="text-muted-fg mt-3 font-mono text-xs">Tested, type-safe, and production-ready.</p>
      </BrutalCard>
    </section>

    <!-- Categories -->
    <section id="categories" class="flex scroll-mt-24 flex-col gap-10" aria-labelledby="categories-title">
      <div v-for="(group, groupIndex) in techContent.categories" :key="group.category" class="flex flex-col gap-4">
        <div class="flex flex-wrap items-center gap-3">
          <h2 class="text-2xl font-bold tracking-tight sm:text-3xl">{{ group.category }}</h2>
          <BrutalBadge :tone="badgeTones[groupIndex % badgeTones.length]"
            >{{ group.items.length }} technologies</BrutalBadge
          >
        </div>
        <p class="text-muted-fg max-w-3xl text-sm leading-relaxed">{{ group.description }}</p>

        <div class="mt-2 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <BrutalCard
            v-for="(item, itemIndex) in group.items"
            :key="item.name"
            :accent="categoryAccents[(groupIndex + itemIndex) % categoryAccents.length]"
            tone="surface"
          >
            <template #header>
              <div class="flex items-center justify-between gap-2">
                <span>{{ item.name }}</span>
                <span
                  v-if="item.version"
                  class="border-[2px] border-fg bg-surface px-2 py-0.5 font-mono text-xs font-bold text-fg"
                >
                  {{ item.version }}
                </span>
              </div>
            </template>

            <p class="text-sm leading-relaxed">{{ item.description }}</p>

            <ul class="mt-4 flex flex-wrap gap-2">
              <li v-for="badge in item.badges" :key="badge">
                <BrutalBadge :tone="badgeTones[(groupIndex + itemIndex + badge.length) % badgeTones.length]">
                  {{ badge }}
                </BrutalBadge>
              </li>
            </ul>
          </BrutalCard>
        </div>
      </div>
    </section>

    <!-- Footer CTA / Info -->
    <section aria-labelledby="tech-footer-title">
      <BrutalCard accent="vivid-yellow" tone="muted">
        <template #header>Architecture Standard</template>

        <div class="flex flex-col gap-3">
          <h2 id="tech-footer-title" class="text-xl font-bold tracking-tight">
            Clean separation, high performance, strict types.
          </h2>
          <p class="text-sm leading-relaxed">
            The backend leverages Go's concurrency and Fiber framework for high throughput, while the frontend provides
            an interactive, reactive user experience built with Vue 3 and Tailwind CSS.
          </p>
        </div>
      </BrutalCard>
    </section>
  </div>
</template>
