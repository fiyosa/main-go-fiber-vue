<script setup lang="ts">
  import { onBeforeUnmount, ref } from 'vue'
  import { RouterLink, useRouter } from 'vue-router'

  import ChangePasswordForm from '../../components/admin/auth/ChangePasswordForm.vue'

  const router = useRouter()

  const succeeded = ref(false)
  const successMessage = ref('')
  let redirectTimer: ReturnType<typeof setTimeout> | undefined

  // The API revokes every session and clears the cookies on success, so the
  // only correct next step is the login page.
  const handleSuccess = (message: string) => {
    successMessage.value = message
    succeeded.value = true
    redirectTimer = setTimeout(() => router.push({ name: 'login' }), 2000)
  }

  onBeforeUnmount(() => {
    if (redirectTimer) clearTimeout(redirectTimer)
  })
</script>

<template>
  <section class="mx-auto flex w-full max-w-lg flex-col gap-6">
    <div class="flex flex-col gap-1">
      <p class="text-xs font-semibold tracking-widest text-[var(--admin-muted)] uppercase">Authentication</p>
      <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Change password</h1>
      <p class="text-sm text-[var(--admin-muted)]">
        Choose a new password for your account. Changing it signs you out of every active session.
      </p>
    </div>

    <div class="rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm">
      <div v-if="succeeded" role="status" class="flex flex-col items-center gap-4 text-center" aria-live="polite">
        <div class="flex size-12 items-center justify-center rounded-full bg-emerald-500 text-white">
          <svg class="size-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
        </div>

        <div class="flex flex-col gap-1">
          <h2 class="text-lg font-bold tracking-tight text-[var(--admin-fg)]">Password changed</h2>
          <p class="text-sm text-[var(--admin-muted)]">
            {{ successMessage || 'Your password was updated.' }} You will be redirected to the login page.
          </p>
        </div>

        <RouterLink
          :to="{ name: 'login' }"
          class="inline-flex min-h-11 items-center justify-center rounded-xl bg-[var(--admin-accent)] px-4 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-[var(--admin-accent-hover)]"
        >
          Sign in again
        </RouterLink>
      </div>

      <ChangePasswordForm v-else @success="handleSuccess" />
    </div>
  </section>
</template>
