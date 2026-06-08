<template>
  <div class="flex flex-col items-center justify-center px-6 py-14 text-center">
    <div class="flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-500">
      <slot name="icon">
        <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
      </slot>
    </div>
    <h3 class="mt-4 text-sm font-semibold text-slate-900 dark:text-white">{{ title }}</h3>
    <p v-if="message" class="mt-1 max-w-sm text-sm text-slate-500 dark:text-slate-400">{{ message }}</p>
    <div v-if="$slots.action" class="mt-5">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: 'Nothing here yet' },
  message: { type: String, default: '' },
})
</script>
