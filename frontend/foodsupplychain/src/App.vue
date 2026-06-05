<template>
  <router-view v-if="$route.meta.blank" />
  <AppLayout v-else />
</template>

<script setup>
import AppLayout from './components/AppLayout.vue'
</script>

<style>
@import './assets/main.css';
</style>
