<template>
  <div class="card !p-0 overflow-hidden">
    <div class="flex items-center justify-between border-b border-slate-100 px-5 py-4 dark:border-slate-800">
      <h3 class="text-base font-semibold text-slate-900 dark:text-white">Users</h3>
      <span class="text-sm text-slate-400">{{ users.length }} total</span>
    </div>

    <SkeletonRows v-if="loading" :rows="4" :cols="4" />

    <div v-else-if="error" class="px-5 py-10 text-center">
      <p class="text-sm text-red-600">{{ error }}</p>
      <button class="btn-secondary mt-3" @click="load">Try again</button>
    </div>

    <div v-else class="overflow-x-auto">
      <table class="min-w-full">
        <thead class="border-b border-slate-100 dark:border-slate-800">
          <tr>
            <th class="th">User</th>
            <th class="th">Email</th>
            <th class="th">Role</th>
            <th class="th">Joined</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
          <tr v-for="u in users" :key="u.id" class="row-hover">
            <td class="td">
              <div class="flex items-center gap-3">
                <div class="flex h-8 w-8 items-center justify-center rounded-full bg-primary-100 text-xs font-semibold uppercase text-primary-700 dark:bg-primary-500/15 dark:text-primary-400">
                  {{ (u.username || 'u').charAt(0) }}
                </div>
                <span class="font-medium text-slate-900 dark:text-white">{{ u.username }}</span>
              </div>
            </td>
            <td class="td">{{ u.email || '—' }}</td>
            <td class="td">
              <select
                class="input !w-36 !py-1.5 text-xs"
                :value="u.role"
                :disabled="u.username === auth.subject"
                @change="changeRole(u, $event.target.value)"
              >
                <option v-for="r in roles" :key="r" :value="r">{{ r }}</option>
              </select>
            </td>
            <td class="td text-slate-400">{{ formatDate(u.created_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { userApi } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'
import SkeletonRows from '@/components/ui/SkeletonRows.vue'

const auth = useAuthStore()
const toast = useToastStore()
const roles = ['admin', 'manager', 'operator', 'viewer']

const users = ref([])
const loading = ref(true)
const error = ref(null)

const formatDate = (d) => (d ? new Date(d).toLocaleDateString() : '—')

const load = async () => {
  loading.value = true
  error.value = null
  try {
    users.value = (await userApi.list()) || []
  } catch (err) {
    error.value = err.message || 'Failed to load users'
  } finally {
    loading.value = false
  }
}

const changeRole = async (user, role) => {
  if (role === user.role) return
  try {
    await userApi.updateRole(user.id, role)
    user.role = role
    toast.success(`${user.username} is now ${role}`)
  } catch (err) {
    toast.error(err.message || 'Failed to update role')
    await load()
  }
}

onMounted(load)
</script>
