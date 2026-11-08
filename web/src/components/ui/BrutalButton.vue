<script setup lang="ts">
  import { computed } from 'vue'

  type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'danger'
  type ButtonSize = 'sm' | 'md'

  const props = withDefaults(
    defineProps<{
      variant?: ButtonVariant
      size?: ButtonSize
      type?: 'button' | 'submit' | 'reset'
      disabled?: boolean
      block?: boolean
    }>(),
    {
      variant: 'primary',
      size: 'md',
      type: 'button',
      disabled: false,
      block: false,
    }
  )

  // Full class strings keep every utility visible to Tailwind's scanner.
  const variantClasses: Record<ButtonVariant, string> = {
    primary: 'bg-vivid-yellow text-black',
    secondary: 'bg-soft-violet text-black',
    ghost: 'bg-surface text-fg',
    danger: 'bg-hot-red text-black',
  }

  // Both sizes stay above the 44px touch target minimum.
  const sizeClasses: Record<ButtonSize, string> = {
    sm: 'min-h-11 px-4 text-xs',
    md: 'min-h-12 px-6 text-sm',
  }

  const classes = computed(() => [
    'brutal-press inline-flex items-center justify-center gap-2 border-[3px] border-fg font-mono font-bold tracking-widest uppercase select-none shadow-brutal-sm',
    'disabled:pointer-events-none disabled:translate-x-0 disabled:translate-y-0 disabled:opacity-60 disabled:shadow-none',
    variantClasses[props.variant],
    sizeClasses[props.size],
    props.block ? 'flex w-full' : '',
  ])
</script>

<template>
  <button :type="type" :disabled="disabled" :class="classes">
    <slot />
  </button>
</template>
