import { defineStore } from 'pinia'
import { ref } from 'vue'

let seq = 0

// Global toast notifications, rendered by ToastHost.
export const useToastStore = defineStore('toast', () => {
  const items = ref([])

  function push(message, type = 'success', timeout = 3500) {
    const id = ++seq
    items.value.push({ id, message, type })
    if (timeout) setTimeout(() => dismiss(id), timeout)
    return id
  }

  function dismiss(id) {
    items.value = items.value.filter((t) => t.id !== id)
  }

  return {
    items,
    push,
    dismiss,
    success: (m) => push(m, 'success'),
    error: (m) => push(m, 'error'),
    info: (m) => push(m, 'info'),
  }
})
