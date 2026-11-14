<script setup lang="ts">
  import { computed, nextTick, ref } from 'vue'

  import { changePassword } from '../../../api/auth'
  import { firstFieldError, toAuthError } from '../../../api/auth/errors'

  const emit = defineEmits<{ success: [message: string] }>()

  const newPassword = ref('')
  const confirmPassword = ref('')
  const showNew = ref(false)
  const showConfirm = ref(false)
  const pending = ref(false)
  const formError = ref('')

  // Field errors keep the server envelope shape (name -> messages) so a 422
  // response can be rendered without a second mapping step. The API uses the
  // JSON field names new-password and confirm-password.
  const fieldErrors = ref<Record<string, string[]>>({})

  const newInput = ref<HTMLInputElement | null>(null)
  const confirmInput = ref<HTMLInputElement | null>(null)

  // Full class strings keep every utility visible to Tailwind's scanner.
  const inputClasses =
    'min-h-11 w-full rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-3.5 text-sm text-[var(--admin-fg)] placeholder:text-[var(--admin-muted)] outline-none transition-colors focus:border-[var(--admin-accent)] disabled:opacity-60'

  const inputTone = (field: string) => (fieldErrors.value[field]?.length ? 'border-red-500' : '')

  const errorClasses = 'mt-1 text-xs font-medium text-red-500'

  const newError = computed(() => firstFieldError(fieldErrors.value, 'new-password'))
  const confirmError = computed(() => firstFieldError(fieldErrors.value, 'confirm-password'))

  const setFieldError = (field: string, message: string) => {
    fieldErrors.value = { ...fieldErrors.value, [field]: [message] }
  }

  // Client-side checks mirror the Go validator so the user gets instant
  // feedback. The server remains the source of truth; a 422 replaces these.
  const validate = (): boolean => {
    fieldErrors.value = {}
    formError.value = ''

    if (!newPassword.value) {
      setFieldError('new-password', 'New password is required.')
    } else if (newPassword.value.length < 8) {
      setFieldError('new-password', 'Password must be at least 8 characters.')
    } else if (newPassword.value.length > 72) {
      setFieldError('new-password', 'Password must be at most 72 characters.')
    }

    if (!confirmPassword.value) {
      setFieldError('confirm-password', 'Please confirm the new password.')
    } else if (confirmPassword.value !== newPassword.value) {
      setFieldError('confirm-password', 'Passwords do not match.')
    }

    return Object.keys(fieldErrors.value).length === 0
  }

  // Focus lands on the first invalid control after the messages render.
  const focusFirstInvalid = async () => {
    await nextTick()

    if (fieldErrors.value['new-password']?.length) {
      newInput.value?.focus()
    } else if (fieldErrors.value['confirm-password']?.length) {
      confirmInput.value?.focus()
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
      const result = await changePassword({
        'new-password': newPassword.value,
        'confirm-password': confirmPassword.value,
      })

      newPassword.value = ''
      confirmPassword.value = ''
      fieldErrors.value = {}
      formError.value = ''
      emit('success', result.message)
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
    <div class="flex flex-col gap-1.5">
      <label for="change-new-password" class="text-xs font-semibold text-[var(--admin-fg)]">New Password</label>

      <div class="relative">
        <input
          id="change-new-password"
          ref="newInput"
          v-model="newPassword"
          :type="showNew ? 'text' : 'password'"
          name="new-password"
          autocomplete="new-password"
          placeholder="Minimum 8 characters"
          :class="[inputClasses, inputTone('new-password'), 'pr-24']"
          :disabled="pending"
          :aria-invalid="newError ? 'true' : 'false'"
          :aria-describedby="newError ? 'change-new-password-error' : undefined"
        />

        <button
          type="button"
          class="absolute inset-y-0 right-1 my-auto flex min-h-10 min-w-16 items-center justify-center rounded-lg border border-[var(--admin-border)] bg-[var(--admin-bg)] px-3 text-xs font-semibold text-[var(--admin-muted)] transition-colors hover:text-[var(--admin-fg)] disabled:opacity-60"
          :disabled="pending"
          :aria-pressed="showNew"
          :aria-label="showNew ? 'Hide new password' : 'Show new password'"
          @click="showNew = !showNew"
        >
          {{ showNew ? 'Hide' : 'Show' }}
        </button>
      </div>

      <p v-if="newError" id="change-new-password-error" :class="errorClasses">{{ newError }}</p>
    </div>

    <div class="flex flex-col gap-1.5">
      <label for="change-confirm-password" class="text-xs font-semibold text-[var(--admin-fg)]">
        Confirm Password
      </label>

      <div class="relative">
        <input
          id="change-confirm-password"
          ref="confirmInput"
          v-model="confirmPassword"
          :type="showConfirm ? 'text' : 'password'"
          name="confirm-password"
          autocomplete="new-password"
          placeholder="Repeat the new password"
          :class="[inputClasses, inputTone('confirm-password'), 'pr-24']"
          :disabled="pending"
          :aria-invalid="confirmError ? 'true' : 'false'"
          :aria-describedby="confirmError ? 'change-confirm-password-error' : undefined"
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

      <p v-if="confirmError" id="change-confirm-password-error" :class="errorClasses">{{ confirmError }}</p>
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
      {{ pending ? 'Changing password…' : 'Change password' }}
    </button>
  </form>
</template>
