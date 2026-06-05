<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-primary-600">FoodSupplyChain</h1>
            </div>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
              <router-link
                v-for="link in links"
                :key="link.to"
                :to="link.to"
                class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                :class="$route.path === link.to ? 'border-primary-500 text-gray-900' : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'"
              >
                {{ link.label }}
              </router-link>
            </div>
          </div>

          <!-- API token control -->
          <div class="hidden sm:flex sm:items-center sm:gap-2">
            <template v-if="auth.isAuthenticated">
              <span class="inline-flex items-center rounded-full bg-primary-100 px-2.5 py-0.5 text-xs font-medium text-primary-800">
                {{ auth.role || 'token' }}
              </span>
              <button class="text-sm text-gray-500 hover:text-gray-700" @click="auth.clear()">Clear token</button>
            </template>
            <template v-else>
              <input
                v-model="tokenInput"
                type="password"
                placeholder="Paste API token…"
                class="input !w-56 !py-1 text-xs"
                @keyup.enter="saveToken"
              />
              <button class="btn-primary !py-1" @click="saveToken">Set</button>
            </template>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const tokenInput = ref('')

const links = [
  { to: '/', label: 'Dashboard' },
  { to: '/inventory', label: 'Inventory' },
  { to: '/shipments', label: 'Shipments' },
  { to: '/catalog', label: 'Catalog' }
]

const saveToken = () => {
  if (!tokenInput.value.trim()) return
  auth.setToken(tokenInput.value)
  tokenInput.value = ''
}
</script>
