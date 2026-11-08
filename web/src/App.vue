<script setup lang="ts">
  import { RouterView } from 'vue-router'

  import secret from './config/secret'

  const isDev = import.meta.env.DEV && secret.APP_ENV === 'development'
</script>

<template>
  <RouterView />
  <agentation-vue v-if="isDev" />
</template>
