<script setup lang="ts">
  import RegisterForm from '../../components/admin/auth/RegisterForm.vue'
</script>

<template>
  <section class="mx-auto flex w-full max-w-lg flex-col gap-6">
    <div class="flex flex-col gap-1">
      <p class="text-xs font-semibold tracking-widest text-[var(--admin-muted)] uppercase">Authentication</p>
      <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Create an account</h1>
      <p class="text-sm text-[var(--admin-muted)]">
        Register a new user account. The account is created without a session, so sign in separately afterwards.
      </p>
    </div>

    <div class="rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-sm">
      <RegisterForm />
    </div>
  </section>
</template>
