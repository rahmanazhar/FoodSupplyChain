<template>
  <transition name="modal">
    <div v-if="open" class="fixed inset-0 z-40 overflow-y-auto" @keydown.esc="$emit('close')">
      <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm" @click="$emit('close')"></div>
      <div class="flex min-h-full items-end justify-center p-4 sm:items-center">
        <div
          class="relative z-10 w-full max-w-lg animate-fade-in rounded-xl border border-slate-200 bg-white p-6 shadow-xl dark:border-slate-800 dark:bg-slate-900"
          role="dialog"
          aria-modal="true"
        >
          <div class="mb-4 flex items-start justify-between">
            <h3 class="text-lg font-semibold text-slate-900 dark:text-white">{{ title }}</h3>
            <button class="rounded-md p-1 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200" @click="$emit('close')" aria-label="Close">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" /></svg>
            </button>
          </div>
          <slot />
          <div v-if="$slots.footer" class="mt-6 flex justify-end gap-3">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
defineProps({
  open: { type: Boolean, required: true },
  title: { type: String, default: '' },
})
defineEmits(['close'])
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
