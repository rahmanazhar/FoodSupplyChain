import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

import './assets/main.css'

// Apply the persisted (or system) theme before mount to avoid a flash.
const storedTheme = localStorage.getItem('fsc.theme')
const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
if (storedTheme === 'dark' || (!storedTheme && prefersDark)) {
  document.documentElement.classList.add('dark')
}

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
