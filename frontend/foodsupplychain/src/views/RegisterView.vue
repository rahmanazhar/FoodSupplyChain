<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center px-4">
    <div class="mx-auto w-full max-w-md">
      <div class="text-center">
        <h1 class="text-2xl font-bold text-primary-600">FoodSupplyChain</h1>
        <p class="mt-2 text-sm text-gray-500">Create an account.</p>
      </div>

      <form class="mt-8 card" @submit.prevent="submit">
        <label class="block text-sm font-medium text-gray-700">Username</label>
        <input v-model="form.username" type="text" autocomplete="username" class="input mt-1" placeholder="at least 3 characters" />

        <label class="mt-4 block text-sm font-medium text-gray-700">Email <span class="text-gray-400">(optional)</span></label>
        <input v-model="form.email" type="email" autocomplete="email" class="input mt-1" placeholder="you@example.com" />

        <label class="mt-4 block text-sm font-medium text-gray-700">Password</label>
        <input v-model="form.password" type="password" autocomplete="new-password" class="input mt-1" placeholder="at least 6 characters" />

        <label class="mt-4 block text-sm font-medium text-gray-700">Confirm password</label>
        <input v-model="form.confirm" type="password" autocomplete="new-password" class="input mt-1" placeholder="repeat password" />

        <p v-if="error" class="mt-4 text-sm text-red-600">{{ error }}</p>

        <button type="submit" class="btn-primary w-full mt-6" :disabled="loading">
          {{ loading ? 'Creating…' : 'Create account' }}
        </button>

        <p class="mt-4 text-sm text-center text-gray-500">
          Already have an account?
          <router-link to="/login" class="text-primary-600 hover:text-primary-700 font-medium">Sign in</router-link>
        </p>

        <p class="mt-3 text-xs text-center text-gray-400">New accounts start with the read-only <strong>viewer</strong> role.</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

const form = reactive({ username: '', email: '', password: '', confirm: '' })
const loading = ref(false)
const error = ref('')

const submit = async () => {
  error.value = ''
  if (form.username.trim().length < 3) {
    error.value = 'Username must be at least 3 characters'
    return
  }
  if (form.password.length < 6) {
    error.value = 'Password must be at least 6 characters'
    return
  }
  if (form.password !== form.confirm) {
    error.value = 'Passwords do not match'
    return
  }
  loading.value = true
  try {
    await auth.register({
      username: form.username.trim(),
      email: form.email.trim(),
      password: form.password
    })
    router.replace('/')
  } catch (err) {
    error.value = err.message || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>
