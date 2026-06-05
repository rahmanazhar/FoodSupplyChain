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

          <!-- Signed-in user -->
          <div class="hidden sm:flex sm:items-center sm:gap-3">
            <div class="text-right leading-tight">
              <p class="text-sm font-medium text-gray-900">{{ auth.subject || 'user' }}</p>
              <p class="text-xs text-primary-600 capitalize">{{ auth.role }}</p>
            </div>
            <div class="h-8 w-8 rounded-full bg-primary-100 flex items-center justify-center">
              <span class="text-primary-800 font-medium uppercase">{{ (auth.subject || 'u').charAt(0) }}</span>
            </div>
            <button class="text-sm text-gray-500 hover:text-gray-700" @click="logout">Sign out</button>
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
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

const links = [
  { to: '/', label: 'Dashboard' },
  { to: '/inventory', label: 'Inventory' },
  { to: '/shipments', label: 'Shipments' },
  { to: '/catalog', label: 'Catalog' }
]

const logout = () => {
  auth.logout()
  router.replace('/login')
}
</script>
