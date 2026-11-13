<script setup lang="ts">
  import { RouterLink, useRoute } from 'vue-router'

  const route = useRoute()
</script>

<template>
  <div
    class="admin-theme relative flex min-h-dvh flex-col items-center justify-center overflow-hidden bg-[var(--admin-bg)] px-6 py-16 text-[var(--admin-fg)] antialiased"
  >
    <!-- Decorative background -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
      <div class="absolute -top-32 left-1/2 h-96 w-96 -translate-x-1/2 rounded-full bg-blue-500/10 blur-3xl" />
      <div
        class="absolute right-0 bottom-0 h-80 w-80 translate-x-1/3 translate-y-1/3 rounded-full bg-indigo-500/10 blur-3xl"
      />
      <div
        class="absolute inset-0 bg-[linear-gradient(to_right,var(--admin-border)_1px,transparent_1px),linear-gradient(to_bottom,var(--admin-border)_1px,transparent_1px)] [mask-image:radial-gradient(ellipse_at_center,black_30%,transparent_70%)] bg-[size:3rem_3rem] opacity-20"
      />
    </div>

    <div class="relative z-10 flex w-full max-w-xl flex-col items-center text-center">
      <span
        class="inline-flex items-center gap-2 rounded-full border border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 py-1.5 text-xs font-semibold tracking-widest text-[var(--admin-muted)] uppercase shadow-xs"
      >
        Error 404
      </span>

      <h1 class="mt-6 text-8xl leading-none font-extrabold tracking-tighter sm:text-9xl">
        <span class="bg-gradient-to-br from-blue-500 to-indigo-500 bg-clip-text text-transparent">404</span>
      </h1>

      <h2 class="mt-6 text-xl font-bold tracking-tight sm:text-2xl">Page not found</h2>

      <p class="mt-3 max-w-md text-sm leading-relaxed text-[var(--admin-muted)]">
        The page you are looking for does not exist or has been moved.
      </p>

      <p
        class="mt-4 max-w-full truncate rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 py-2 font-mono text-xs text-[var(--admin-muted)] shadow-xs"
        :title="route.fullPath"
      >
        {{ route.fullPath }}
      </p>

      <div class="mt-8 flex flex-col items-center gap-3 sm:flex-row">
        <RouterLink
          to="/admin"
          class="inline-flex min-h-11 items-center justify-center gap-2 rounded-xl bg-blue-600 px-6 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-blue-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-500"
        >
          <svg class="size-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7M8 12h13" />
          </svg>
          Back to dashboard
        </RouterLink>

        <RouterLink
          to="/"
          class="inline-flex min-h-11 items-center justify-center gap-2 rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-6 text-sm font-semibold text-[var(--admin-fg)] shadow-xs transition-colors hover:bg-[var(--admin-border)]/50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-500"
        >
          Go to homepage
        </RouterLink>
      </div>
    </div>
  </div>
</template>
