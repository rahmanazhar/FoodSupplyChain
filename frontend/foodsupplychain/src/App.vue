<template>
  <router-view v-if="$route.meta.blank" />
  <AppLayout v-else />
  <ToastHost />
</template>

<script setup>
import AppLayout from './components/AppLayout.vue'
import ToastHost from './components/ui/ToastHost.vue'
</script>

<style>
@import './assets/main.css';
</style>
