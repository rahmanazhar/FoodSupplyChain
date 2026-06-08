<template>
  <div class="pointer-events-none fixed bottom-4 right-4 z-50 flex flex-col gap-2">
    <transition-group name="toast">
      <div
        v-for="t in toast.items"
        :key="t.id"
        class="pointer-events-auto flex min-w-[260px] max-w-sm items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium text-white shadow-card-hover"
        :class="styles[t.type] || styles.info"
      >
        <svg v-if="t.type === 'success'" class="h-5 w-5 shrink-0" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.7 5.3a1 1 0 010 1.4l-7.5 7.5a1 1 0 01-1.4 0L3.3 9.7a1 1 0 011.4-1.4l3.1 3.1 6.8-6.8a1 1 0 011.4 0z" clip-rule="evenodd"/></svg>
        <svg v-else-if="t.type === 'error'" class="h-5 w-5 shrink-0" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.7 7.3a1 1 0 00-1.4 1.4L8.6 10l-1.3 1.3a1 1 0 101.4 1.4l1.3-1.3 1.3 1.3a1 1 0 001.4-1.4L11.4 10l1.3-1.3a1 1 0 00-1.4-1.4L10 8.6 8.7 7.3z" clip-rule="evenodd"/></svg>
        <svg v-else class="h-5 w-5 shrink-0" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M18 10A8 8 0 11 2 10a8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/></svg>
        <span class="flex-1">{{ t.message }}</span>
        <button class="opacity-70 transition hover:opacity-100" @click="toast.dismiss(t.id)" aria-label="Dismiss">✕</button>
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()
const styles = {
  success: 'bg-emerald-600',
  error: 'bg-red-600',
  info: 'bg-slate-800 dark:bg-slate-700',
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(16px);
}
</style>
