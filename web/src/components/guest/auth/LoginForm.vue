<script setup lang="ts">
  import { nextTick, computed, ref } from 'vue'

  import { login } from '../../../api/auth'
  import { firstFieldError, toAuthError } from '../../../api/auth/errors'
  import { BrutalButton } from '../ui'

  const emit = defineEmits<{ success: [] }>()

  const email = ref('')
  const password = ref('')
  const showPassword = ref(false)
  const pending = ref(false)
  const formError = ref('')

  // Field errors keep the server envelope shape (name -> messages) so a 422
  // response can be rendered without a second mapping step.
  const fieldErrors = ref<Record<string, string[]>>({})

  const emailInput = ref<HTMLInputElement | null>(null)
  const passwordInput = ref<HTMLInputElement | null>(null)

  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

  // Full class strings keep every utility visible to Tailwind's scanner.
  const inputClasses =
    'min-h-12 w-full border-[3px] bg-surface px-4 font-mono text-sm text-fg shadow-brutal-sm placeholder:text-fg/50 disabled:opacity-60'

  const inputTone = (field: string) => (fieldErrors.value[field]?.length ? 'border-hot-red' : 'border-fg')

  const errorClasses = 'mt-1 inline-block bg-hot-red px-2 py-1 font-mono text-xs leading-relaxed font-bold text-black'

  const emailError = computed(() => firstFieldError(fieldErrors.value, 'email'))
  const passwordError = computed(() => firstFieldError(fieldErrors.value, 'password'))

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
    }
  }

  const submit = async () => {
    if (pending.value) return

    if (!validate()) {
      await focusFirstInvalid()
      return
    }

    pending.value = true

    try {
      await login({ email: email.value.trim(), password: password.value })
      emit('success')
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
  <form class="flex flex-col gap-5" novalidate @submit.prevent="submit">
    <div class="flex flex-col gap-1">
      <label for="login-email" class="font-mono text-xs font-bold tracking-widest uppercase">Email</label>

      <input
        id="login-email"
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
        :aria-describedby="emailError ? 'login-email-error' : undefined"
      />

      <p v-if="emailError" id="login-email-error" :class="errorClasses">{{ emailError }}</p>
    </div>

    <div class="flex flex-col gap-1">
      <label for="login-password" class="font-mono text-xs font-bold tracking-widest uppercase">Password</label>

      <div class="relative">
        <input
          id="login-password"
          ref="passwordInput"
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          name="password"
          autocomplete="current-password"
          placeholder="••••••••"
          :class="[inputClasses, inputTone('password'), 'pr-24']"
          :disabled="pending"
          :aria-invalid="passwordError ? 'true' : 'false'"
          :aria-describedby="passwordError ? 'login-password-error' : undefined"
        />

        <button
          type="button"
          class="absolute inset-y-0 right-1 my-auto flex min-h-11 min-w-16 items-center justify-center border-[3px] border-fg bg-muted-bg px-3 font-mono text-xs font-bold tracking-widest uppercase disabled:opacity-60"
          :disabled="pending"
          :aria-pressed="showPassword"
          :aria-label="showPassword ? 'Hide password' : 'Show password'"
          @click="showPassword = !showPassword"
        >
          {{ showPassword ? 'Hide' : 'Show' }}
        </button>
      </div>

      <p v-if="passwordError" id="login-password-error" :class="errorClasses">{{ passwordError }}</p>
    </div>

    <p
      v-if="formError"
      role="alert"
      class="border-[3px] border-fg bg-hot-red px-4 py-3 font-mono text-xs leading-relaxed font-bold text-black"
    >
      {{ formError }}
    </p>

    <BrutalButton type="submit" variant="primary" block :disabled="pending" :aria-busy="pending">
      {{ pending ? 'Signing in…' : 'Sign in' }}
    </BrutalButton>
  </form>
</template>
