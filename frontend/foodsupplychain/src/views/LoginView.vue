<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center px-4">
    <div class="mx-auto w-full max-w-md">
      <div class="text-center">
        <h1 class="text-2xl font-bold text-primary-600">FoodSupplyChain</h1>
        <p class="mt-2 text-sm text-gray-500">Sign in to manage inventory and shipments.</p>
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

        <p class="mt-4 text-sm text-center text-gray-500">
          No account?
          <router-link to="/register" class="text-primary-600 hover:text-primary-700 font-medium">Create one</router-link>
        </p>

        <div class="mt-4 rounded-md bg-gray-50 p-3 text-xs text-gray-500">
          Demo accounts (seeded): <code>admin / admin123</code>, <code>manager / manager123</code>,
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
