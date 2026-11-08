<script setup lang="ts">
  import { computed } from 'vue'

  type BadgeTone = 'surface' | 'hot-red' | 'vivid-yellow' | 'soft-violet' | 'mint' | 'sky'

  const props = withDefaults(
    defineProps<{
      tone?: BadgeTone
      rotation?: number
    }>(),
    {
      tone: 'vivid-yellow',
      rotation: 0,
    }
  )

  // Full class strings keep every utility visible to Tailwind's scanner.
  const toneClasses: Record<BadgeTone, string> = {
    surface: 'bg-surface text-fg',
    'hot-red': 'bg-hot-red text-black',
    'vivid-yellow': 'bg-vivid-yellow text-black',
    'soft-violet': 'bg-soft-violet text-black',
    mint: 'bg-mint text-black',
    sky: 'bg-sky text-black',
  }

  const classes = computed(() => [
    'inline-flex items-center gap-1 border-[3px] border-fg px-3 py-1 font-mono text-xs font-bold tracking-widest uppercase shadow-brutal-sm',
    toneClasses[props.tone],
  ])

  // Rotation stays inline so callers can pass any sticker angle.
  const style = computed(() => (props.rotation ? { transform: `rotate(${props.rotation}deg)` } : undefined))
</script>

<template>
  <span :class="classes" :style="style">
    <slot />
  </span>
</template>
