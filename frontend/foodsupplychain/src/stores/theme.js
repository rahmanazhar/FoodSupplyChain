import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

const KEY = 'fsc.theme'

export const useThemeStore = defineStore('theme', () => {
  const stored = localStorage.getItem(KEY)
  const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
  const isDark = ref(stored ? stored === 'dark' : prefersDark)

  function apply() {
    document.documentElement.classList.toggle('dark', isDark.value)
  }
  apply()

  watch(isDark, (value) => {
    localStorage.setItem(KEY, value ? 'dark' : 'light')
    apply()
  })

  function toggle() {
    isDark.value = !isDark.value
  }

  return { isDark, toggle }
})
