<script setup lang="ts">
  import { computed } from 'vue'

  type CardTone = 'surface' | 'muted'
  type CardAccent = 'none' | 'hot-red' | 'vivid-yellow' | 'soft-violet' | 'mint' | 'sky'

  const props = withDefaults(
    defineProps<{
      tone?: CardTone
      accent?: CardAccent
    }>(),
    {
      tone: 'surface',
      accent: 'none',
    }
  )

  // Full class strings keep every utility visible to Tailwind's scanner.
  const toneClasses: Record<CardTone, string> = {
    surface: 'bg-surface',
    muted: 'bg-muted-bg',
  }

  const accentClasses: Record<CardAccent, string> = {
    none: '',
    'hot-red': 'bg-hot-red text-black',
    'vivid-yellow': 'bg-vivid-yellow text-black',
    'soft-violet': 'bg-soft-violet text-black',
    mint: 'bg-mint text-black',
    sky: 'bg-sky text-black',
  }

  const classes = computed(() => ['border-[3px] border-fg text-fg shadow-brutal', toneClasses[props.tone]])

  const headerClasses = computed(() => [
    'border-b-[3px] border-fg px-6 py-4 font-mono text-xs font-bold tracking-[0.2em] uppercase',
    accentClasses[props.accent],
  ])
</script>

<template>
  <div :class="classes">
    <div v-if="$slots.header" :class="headerClasses">
      <slot name="header" />
    </div>

    <div class="p-6">
      <slot />
    </div>

    <div v-if="$slots.footer" class="border-t-[3px] border-fg px-6 py-4">
      <slot name="footer" />
    </div>
  </div>
</template>
