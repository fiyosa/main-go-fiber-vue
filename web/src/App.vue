<script setup lang="ts">
  import { RouterView } from 'vue-router'

  const showAgentation = import.meta.env.DEV
</script>

<template>
  <RouterView />

  <agentation-vue v-if="showAgentation" />
</template>
