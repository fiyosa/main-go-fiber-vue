<script setup lang="ts">
  import { useRouter } from 'vue-router'

  import { img } from '../../assets/img'
  import { BrutalBadge, BrutalCard } from '../../components/guest/ui'
  import LoginForm from '../../components/guest/auth/LoginForm.vue'

  const router = useRouter()

  // The guest header and footer stay mounted, so the page only owns the card
  // and the redirect after a successful sign-in.
  const onSuccess = () => {
    void router.push({ name: 'admin' })
  }
</script>

<template>
  <section class="mx-auto flex w-full max-w-4xl flex-col gap-6">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="font-mono text-xs font-bold tracking-[0.3em] uppercase">Members area</p>

      <BrutalBadge tone="hot-red" :rotation="-2">Members only</BrutalBadge>
    </div>

    <BrutalCard accent="vivid-yellow">
      <template #header>Sign in</template>

      <div class="grid gap-8 lg:grid-cols-[1.1fr_0.9fr]">
        <div class="flex flex-col gap-4">
          <h1 class="text-3xl leading-none font-bold tracking-tight sm:text-4xl">Welcome back</h1>

          <p class="max-w-md text-sm leading-relaxed">
            Sign in with the email on your account. Session tokens are stored in HTTP-only cookies, so nothing is kept
            in the browser.
          </p>

          <LoginForm class="mt-2" @success="onSuccess" />
        </div>

        <aside
          aria-hidden="true"
          class="hidden items-center justify-center border-[3px] border-fg bg-hot-red p-4 shadow-brutal-sm lg:flex"
        >
          <img :src="img.loginSecure" alt="" class="h-full w-full object-contain" />
        </aside>
      </div>
    </BrutalCard>
  </section>
</template>
