<script setup lang="ts">
  import { computed, nextTick, ref } from 'vue'

  import { register } from '../../../api/auth'
  import { firstFieldError, toAuthError } from '../../../api/auth/errors'

  const emit = defineEmits<{ success: [email: string] }>()

  const email = ref('')
  const password = ref('')
  const confirmPassword = ref('')
  const showPassword = ref(false)
  const showConfirm = ref(false)
  const pending = ref(false)
  const succeeded = ref(false)
  const registeredEmail = ref('')
  const formError = ref('')

  // Field errors keep the server envelope shape (name -> messages) so a 422
  // response can be rendered without a second mapping step.
  const fieldErrors = ref<Record<string, string[]>>({})

  const emailInput = ref<HTMLInputElement | null>(null)
  const passwordInput = ref<HTMLInputElement | null>(null)
  const confirmInput = ref<HTMLInputElement | null>(null)

  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

  // Full class strings keep every utility visible to Tailwind's scanner.
  const inputClasses =
    'min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3.5 text-sm text-[var(--admin-fg)] placeholder:text-[var(--admin-muted)] outline-none transition-colors focus:border-[var(--admin-accent)] disabled:opacity-60'

  const inputTone = (field: string) => (fieldErrors.value[field]?.length ? 'border-red-500' : '')

  const errorClasses = 'mt-1 text-xs font-medium text-red-500'

  const emailError = computed(() => firstFieldError(fieldErrors.value, 'email'))
  const passwordError = computed(() => firstFieldError(fieldErrors.value, 'password'))
  const confirmError = computed(() => firstFieldError(fieldErrors.value, 'confirmPassword'))

  const setFieldError = (field: string, message: string) => {
    fieldErrors.value = { ...fieldErrors.value, [field]: [message] }
  }

  // Client-side checks mirror the Go validator so the user gets instant
  // feedback. The server remains the source of truth; a 422 replaces these.
  const validate = (): boolean => {
    fieldErrors.value = {}
    formError.value = ''

    if (!email.value.trim()) {
      setFieldError('email', 'Email is required.')
    } else if (!emailPattern.test(email.value.trim())) {
      setFieldError('email', 'Enter a valid email address.')
    }

    if (!password.value) {
      setFieldError('password', 'Password is required.')
    } else if (password.value.length < 8) {
      setFieldError('password', 'Password must be at least 8 characters.')
    } else if (password.value.length > 72) {
      setFieldError('password', 'Password must be at most 72 characters.')
    }

    if (!confirmPassword.value) {
      setFieldError('confirmPassword', 'Please confirm your password.')
    } else if (confirmPassword.value !== password.value) {
      setFieldError('confirmPassword', 'Passwords do not match.')
    }

    return Object.keys(fieldErrors.value).length === 0
  }

  // Focus lands on the first invalid control after the messages render.
  const focusFirstInvalid = async () => {
    await nextTick()

    if (fieldErrors.value.email?.length) {
      emailInput.value?.focus()
    } else if (fieldErrors.value.password?.length) {
      passwordInput.value?.focus()
    } else if (fieldErrors.value.confirmPassword?.length) {
      confirmInput.value?.focus()
    }
  }

  const resetForm = () => {
    email.value = ''
    password.value = ''
    confirmPassword.value = ''
    fieldErrors.value = {}
    formError.value = ''
    succeeded.value = false
    registeredEmail.value = ''
  }

  const submit = async () => {
    if (pending.value) return

    if (!validate()) {
      await focusFirstInvalid()
      return
    }

    pending.value = true

    try {
      const result = await register({ email: email.value.trim(), password: password.value })

      registeredEmail.value = result.data?.email || email.value.trim()
      email.value = ''
      password.value = ''
      confirmPassword.value = ''
      fieldErrors.value = {}
      formError.value = ''
      succeeded.value = true
      emit('success', registeredEmail.value)
    } catch (error) {
      const authError = toAuthError(error)
      const hasFieldErrors = Object.keys(authError.fields).length > 0

      fieldErrors.value = authError.fields
      formError.value = hasFieldErrors ? '' : authError.message

      await focusFirstInvalid()
    } finally {
      pending.value = false
    }
  }
</script>

<template>
  <!-- Success state: registration completed but the API does not auto-login. -->
  <div
    v-if="succeeded"
    role="status"
    class="flex flex-col items-center gap-4 rounded-2xl border border-emerald-500/30 bg-emerald-500/10 p-8 text-center"
  >
    <div class="flex size-12 items-center justify-center rounded-full bg-emerald-500 text-white">
      <svg class="size-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
      </svg>
    </div>

    <div class="flex flex-col gap-1">
      <h2 class="text-lg font-bold tracking-tight text-[var(--admin-fg)]">Account created</h2>
      <p class="text-sm text-[var(--admin-muted)]">
        {{ registeredEmail || 'The new account' }} was registered successfully. Sign in separately to continue.
      </p>
    </div>

    <button
      type="button"
      class="inline-flex min-h-11 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 text-sm font-semibold text-[var(--admin-fg)] transition-colors hover:bg-[var(--admin-border)]/50"
      @click="resetForm"
    >
      Register another account
    </button>
  </div>

  <form v-else class="flex flex-col gap-5" novalidate @submit.prevent="submit">
    <div class="flex flex-col gap-1.5">
      <label for="register-email" class="text-xs font-semibold text-[var(--admin-fg)]">Email</label>

      <input
        id="register-email"
        ref="emailInput"
        v-model="email"
        type="email"
        name="email"
        autocomplete="email"
        inputmode="email"
        placeholder="you@example.com"
        :class="[inputClasses, inputTone('email')]"
        :disabled="pending"
        :aria-invalid="emailError ? 'true' : 'false'"
        :aria-describedby="emailError ? 'register-email-error' : undefined"
      />

      <p v-if="emailError" id="register-email-error" :class="errorClasses">{{ emailError }}</p>
    </div>

    <div class="flex flex-col gap-1.5">
      <label for="register-password" class="text-xs font-semibold text-[var(--admin-fg)]">Password</label>

      <div class="relative">
        <input
          id="register-password"
          ref="passwordInput"
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          name="password"
          autocomplete="new-password"
          placeholder="Minimum 8 characters"
          :class="[inputClasses, inputTone('password'), 'pr-24']"
          :disabled="pending"
          :aria-invalid="passwordError ? 'true' : 'false'"
          :aria-describedby="passwordError ? 'register-password-error' : undefined"
        />

        <button
          type="button"
          class="absolute inset-y-0 right-1 my-auto flex min-h-10 min-w-16 items-center justify-center rounded-lg border border-[var(--admin-border)] bg-[var(--admin-bg)] px-3 text-xs font-semibold text-[var(--admin-muted)] transition-colors hover:text-[var(--admin-fg)] disabled:opacity-60"
          :disabled="pending"
          :aria-pressed="showPassword"
          :aria-label="showPassword ? 'Hide password' : 'Show password'"
          @click="showPassword = !showPassword"
        >
          {{ showPassword ? 'Hide' : 'Show' }}
        </button>
      </div>

      <p v-if="passwordError" id="register-password-error" :class="errorClasses">{{ passwordError }}</p>
    </div>

    <div class="flex flex-col gap-1.5">
      <label for="register-confirm" class="text-xs font-semibold text-[var(--admin-fg)]">Confirm Password</label>

      <div class="relative">
        <input
          id="register-confirm"
          ref="confirmInput"
          v-model="confirmPassword"
          :type="showConfirm ? 'text' : 'password'"
          name="confirmPassword"
          autocomplete="new-password"
          placeholder="Repeat your password"
          :class="[inputClasses, inputTone('confirmPassword'), 'pr-24']"
          :disabled="pending"
          :aria-invalid="confirmError ? 'true' : 'false'"
          :aria-describedby="confirmError ? 'register-confirm-error' : undefined"
        />

        <button
          type="button"
          class="absolute inset-y-0 right-1 my-auto flex min-h-10 min-w-16 items-center justify-center rounded-lg border border-[var(--admin-border)] bg-[var(--admin-bg)] px-3 text-xs font-semibold text-[var(--admin-muted)] transition-colors hover:text-[var(--admin-fg)] disabled:opacity-60"
          :disabled="pending"
          :aria-pressed="showConfirm"
          :aria-label="showConfirm ? 'Hide password confirmation' : 'Show password confirmation'"
          @click="showConfirm = !showConfirm"
        >
          {{ showConfirm ? 'Hide' : 'Show' }}
        </button>
      </div>

      <p v-if="confirmError" id="register-confirm-error" :class="errorClasses">{{ confirmError }}</p>
    </div>

    <p
      v-if="formError"
      role="alert"
      class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-xs font-medium text-red-500"
    >
      {{ formError }}
    </p>

    <button
      type="submit"
      class="inline-flex min-h-11 items-center justify-center rounded-xl bg-[var(--admin-accent)] px-4 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-[var(--admin-accent-hover)] disabled:opacity-60"
      :disabled="pending"
      :aria-busy="pending"
    >
      {{ pending ? 'Creating account…' : 'Create account' }}
    </button>
  </form>
</template>
