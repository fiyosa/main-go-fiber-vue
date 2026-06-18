<script setup lang="ts">
  import { ImgHero, ImgVite, ImgVue } from '../../assets/img'
</script>

<template>
  <img :src="ImgHero" alt="hero" />
  <img :src="ImgVite" alt="vite" />
  <img :src="ImgVue" alt="vue" />
</template>
