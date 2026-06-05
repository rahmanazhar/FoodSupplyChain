<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center px-4">
    <div class="mx-auto w-full max-w-md">
      <div class="text-center">
        <h1 class="text-2xl font-bold text-primary-600">FoodSupplyChain</h1>
        <p class="mt-2 text-sm text-gray-500">Sign in to manage inventory and shipments.</p>
      </div>

      <div class="mt-8 card">
        <label class="block text-sm font-medium text-gray-700">Username <span class="text-gray-400">(optional)</span></label>
        <input v-model="username" type="text" class="input mt-1" placeholder="e.g. alice" @keyup.enter="signIn" />

        <label class="mt-4 block text-sm font-medium text-gray-700">Role</label>
        <div class="mt-2 grid grid-cols-2 gap-2">
          <button
            v-for="r in roles"
            :key="r.value"
            type="button"
            class="rounded-md border px-3 py-2 text-sm text-left transition"
            :class="role === r.value ? 'border-primary-500 ring-2 ring-primary-200 bg-primary-50' : 'border-gray-300 hover:border-gray-400'"
            @click="role = r.value"
          >
            <span class="block font-medium text-gray-900 capitalize">{{ r.value }}</span>
            <span class="block text-xs text-gray-500">{{ r.hint }}</span>
          </button>
        </div>

        <p v-if="error" class="mt-4 text-sm text-red-600">{{ error }}</p>

        <button type="button" class="btn-primary w-full mt-6" :disabled="loading" @click="signIn">
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </button>

        <p class="mt-4 text-xs text-gray-400 text-center">
          Tokens are issued by the gateway and signed with a per-run secret — nothing to paste.
        </p>
      </div>
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

const roles = [
  { value: 'admin', hint: 'Full access incl. delete' },
  { value: 'manager', hint: 'Manage & delete shipments' },
  { value: 'operator', hint: 'Create & update' },
  { value: 'viewer', hint: 'Read-only' }
]

const username = ref('')
const role = ref('admin')
const loading = ref(false)
const error = ref('')

const signIn = async () => {
  loading.value = true
  error.value = ''
  try {
    await auth.login(role.value, username.value.trim())
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.replace(redirect)
  } catch (err) {
    error.value = err.message || 'Sign in failed'
  } finally {
    loading.value = false
  }
}
</script>
