<script setup lang="ts">
  import { watch, nextTick, ref } from 'vue'
  import { useUIStore } from '../../../store/ui'

  const ui = useUIStore()
  const confirmBtnRef = ref<HTMLButtonElement | null>(null)

  const typeStyles = {
    danger: {
      iconBg: 'bg-red-500/10',
      iconText: 'text-red-500',
      btnClass: 'bg-red-600 hover:bg-red-700 text-white',
      iconPath: 'M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z',
    },
    warning: {
      iconBg: 'bg-amber-500/10',
      iconText: 'text-amber-500',
      btnClass: 'bg-amber-500 hover:bg-amber-600 text-white',
      iconPath:
        'M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z',
    },
    info: {
      iconBg: 'bg-blue-500/10',
      iconText: 'text-blue-500',
      btnClass: 'bg-blue-600 hover:bg-blue-700 text-white',
      iconPath:
        'M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z',
    },
  }

  watch(
    () => ui.confirmModal.isOpen,
    async (isOpen) => {
      if (!isOpen) return
      await nextTick()
      confirmBtnRef.value?.focus()
    }
  )

  const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      ui.confirmModal.onCancel()
    }
  }
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="ui.confirmModal.isOpen"
        class="admin-theme fixed inset-0 z-[60] flex items-center justify-center bg-black/50 backdrop-blur-md"
        role="dialog"
        aria-modal="true"
        :aria-label="ui.confirmModal.title"
        @keydown="handleKeydown"
        @click.self="ui.confirmModal.onCancel()"
      >
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
          appear
        >
          <div class="w-full max-w-sm rounded-2xl border border-[var(--admin-border)] bg-[var(--admin-surface)] p-6 shadow-2xl">
            <div class="flex items-start gap-4">
              <div
                class="flex size-10 shrink-0 items-center justify-center rounded-full"
                :class="[typeStyles[ui.confirmModal.type].iconBg, typeStyles[ui.confirmModal.type].iconText]"
              >
                <svg
                  class="size-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    :d="typeStyles[ui.confirmModal.type].iconPath"
                  />
                </svg>
              </div>

              <div class="min-w-0 flex-1">
                <h3 class="text-base font-bold tracking-tight text-[var(--admin-fg)]">
                  {{ ui.confirmModal.title }}
                </h3>
                <p class="mt-1 text-sm text-[var(--admin-fg)]">
                  {{ ui.confirmModal.message }}
                </p>
              </div>
            </div>

            <div class="mt-6 flex justify-end gap-3">
              <button
                type="button"
                class="inline-flex min-h-10 items-center justify-center rounded-xl border border-[var(--admin-border)] bg-[var(--admin-surface)] px-4 text-xs font-semibold text-[var(--admin-fg)] transition-colors hover:bg-[var(--admin-border)]/50"
                @click="ui.confirmModal.onCancel()"
              >
                {{ ui.confirmModal.cancelText }}
              </button>
              <button
                ref="confirmBtnRef"
                type="button"
                class="inline-flex min-h-10 items-center justify-center rounded-xl px-4 text-xs font-semibold shadow-sm transition-colors"
                :class="typeStyles[ui.confirmModal.type].btnClass"
                @click="ui.confirmModal.onConfirm()"
              >
                {{ ui.confirmModal.confirmText }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
