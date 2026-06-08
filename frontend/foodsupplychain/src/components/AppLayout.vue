<template>
  <div class="min-h-screen bg-slate-50 dark:bg-slate-950">
    <!-- Sidebar -->
    <aside
      class="fixed inset-y-0 left-0 z-30 flex w-64 flex-col border-r border-slate-200 bg-white transition-transform dark:border-slate-800 dark:bg-slate-900 lg:translate-x-0"
      :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <div class="flex h-16 items-center gap-2 px-6">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary-600 text-white">
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
        </div>
        <span class="text-lg font-bold tracking-tight text-slate-900 dark:text-white">FoodSupplyChain</span>
      </div>

      <nav class="flex-1 space-y-1 px-3 py-4">
        <router-link
          v-for="link in visibleLinks"
          :key="link.to"
          :to="link.to"
          class="group flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors"
          :class="isActive(link.to)
            ? 'bg-primary-50 text-primary-700 dark:bg-primary-500/10 dark:text-primary-400'
            : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-slate-100'"
          @click="sidebarOpen = false"
        >
          <svg class="h-5 w-5 shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.6" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" :d="icons[link.icon]" /></svg>
          {{ link.label }}
        </router-link>
      </nav>

      <!-- Sidebar footer: theme + user -->
      <div class="border-t border-slate-200 p-3 dark:border-slate-800">
        <button
          class="mb-2 flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-800"
          @click="theme.toggle()"
        >
          <svg v-if="theme.isDark" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.6" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" /></svg>
          <svg v-else class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.6" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.72 9.72 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" /></svg>
          {{ theme.isDark ? 'Light mode' : 'Dark mode' }}
        </button>

        <div class="flex items-center gap-3 rounded-lg px-3 py-2">
          <div class="flex h-9 w-9 items-center justify-center rounded-full bg-primary-100 font-semibold uppercase text-primary-700 dark:bg-primary-500/15 dark:text-primary-400">
            {{ (auth.subject || 'u').charAt(0) }}
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-medium text-slate-900 dark:text-white">{{ auth.subject || 'user' }}</p>
            <p class="text-xs capitalize text-primary-600 dark:text-primary-400">{{ auth.role }}</p>
          </div>
          <button class="rounded-md p-1.5 text-slate-400 hover:bg-slate-100 hover:text-slate-700 dark:hover:bg-slate-800 dark:hover:text-slate-200" title="Sign out" @click="logout">
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.6" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" /></svg>
          </button>
        </div>
      </div>
    </aside>

    <!-- Mobile overlay -->
    <div v-if="sidebarOpen" class="fixed inset-0 z-20 bg-slate-900/50 lg:hidden" @click="sidebarOpen = false"></div>

    <!-- Main -->
    <div class="lg:pl-64">
      <header class="sticky top-0 z-10 flex h-16 items-center gap-4 border-b border-slate-200 bg-white/80 px-4 backdrop-blur dark:border-slate-800 dark:bg-slate-900/80 sm:px-6">
        <button class="rounded-md p-2 text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 lg:hidden" @click="sidebarOpen = true">
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.6" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5M3.75 17.25h16.5" /></svg>
        </button>
        <h1 class="text-lg font-semibold text-slate-900 dark:text-white">{{ pageTitle }}</h1>
      </header>

      <main class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'

const auth = useAuthStore()
const theme = useThemeStore()
const route = useRoute()
const router = useRouter()

const sidebarOpen = ref(false)

const links = [
  { to: '/', label: 'Dashboard', icon: 'dashboard' },
  { to: '/inventory', label: 'Inventory', icon: 'inventory' },
  { to: '/shipments', label: 'Shipments', icon: 'shipments' },
  { to: '/catalog', label: 'Catalog', icon: 'catalog' },
  { to: '/users', label: 'Users', icon: 'users', adminOnly: true },
]

const visibleLinks = computed(() => links.filter((l) => !l.adminOnly || auth.role === 'admin'))

const icons = {
  dashboard: 'M3.75 3.75h6v6h-6v-6zm10.5 0h6v6h-6v-6zm-10.5 10.5h6v6h-6v-6zm10.5 0h6v6h-6v-6z',
  inventory: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4',
  shipments: 'M9 17a2 2 0 11-4 0 2 2 0 014 0zm10 0a2 2 0 11-4 0 2 2 0 014 0zM13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1',
  catalog: 'M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3z M6 6h.008v.008H6V6z',
  users: 'M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z',
}

const titles = {
  '/': 'Dashboard',
  '/inventory': 'Inventory',
  '/shipments': 'Shipments',
  '/catalog': 'Catalog',
  '/users': 'Users',
}
const pageTitle = computed(() => titles[route.path] || 'FoodSupplyChain')

const isActive = (to) => (to === '/' ? route.path === '/' : route.path.startsWith(to))

const logout = () => {
  auth.logout()
  router.replace('/login')
}
</script>

<style scoped>
.page-enter-active,
.page-leave-active {
  transition: opacity 0.15s ease;
}
.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>
