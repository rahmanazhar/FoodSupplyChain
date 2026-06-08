<template>
  <div class="flex min-h-screen flex-col justify-center bg-slate-50 px-4 dark:bg-slate-950">
    <div class="mx-auto w-full max-w-md">
      <div class="text-center">
        <div class="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-primary-600 text-white">
          <svg class="h-7 w-7" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
        </div>
        <h1 class="text-2xl font-bold text-slate-900 dark:text-white">FoodSupplyChain</h1>
        <p class="mt-2 text-sm text-slate-500 dark:text-slate-400">Sign in to manage inventory and shipments.</p>
      </div>

      <form class="mt-8 card" @submit.prevent="signIn">
        <label class="block text-sm font-medium text-gray-700">Username</label>
        <input v-model="username" type="text" autocomplete="username" class="input mt-1" placeholder="username" />

        <label class="mt-4 block text-sm font-medium text-gray-700">Password</label>
        <input v-model="password" type="password" autocomplete="current-password" class="input mt-1" placeholder="••••••••" />

        <p v-if="error" class="mt-4 text-sm text-red-600">{{ error }}</p>

        <button type="submit" class="btn-primary w-full mt-6" :disabled="loading">
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </button>

        <p class="mt-4 text-center text-sm text-slate-500 dark:text-slate-400">
          No account?
          <router-link to="/register" class="font-medium text-primary-600 hover:text-primary-700">Create one</router-link>
        </p>

        <div class="mt-4 rounded-lg bg-slate-50 p-3 text-xs text-slate-500 dark:bg-slate-800/60 dark:text-slate-400">
          Demo accounts: <code>admin / admin123</code>, <code>manager / manager123</code>,
          <code>operator / operator123</code>, <code>viewer / viewer123</code>.
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const signIn = async () => {
  if (!username.value || !password.value) {
    error.value = 'Enter your username and password'
    return
  }
  loading.value = true
  error.value = ''
  try {
    await auth.login(username.value.trim(), password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.replace(redirect)
  } catch (err) {
    error.value = err.message || 'Sign in failed'
  } finally {
    loading.value = false
  }
}
</script>
