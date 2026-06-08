<template>
  <div class="space-y-5">
    <!-- Toolbar -->
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex flex-1 flex-col gap-3 sm:flex-row sm:items-center">
        <div class="relative w-full sm:max-w-xs">
          <svg class="pointer-events-none absolute left-3 top-2.5 h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="m21 21-4.3-4.3m1.8-4.7a6.5 6.5 0 1 1-13 0 6.5 6.5 0 0 1 13 0z" /></svg>
          <input v-model="search" type="text" class="input pl-9" placeholder="Search order, origin, destination…" />
        </div>
        <select v-model="statusFilter" class="input sm:w-44" @change="reset">
          <option value="">All statuses</option>
          <option v-for="s in statuses" :key="s" :value="s">{{ s.replace('_', ' ') }}</option>
        </select>
      </div>
      <button class="btn-primary" @click="modalOpen = true">
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
        Create Shipment
      </button>
    </div>

    <div v-if="error" class="card text-center">
      <p class="text-sm text-red-600">{{ error }}</p>
      <button class="btn-secondary mt-3" @click="load">Try again</button>
    </div>

    <div v-else-if="loading" class="space-y-4">
      <div v-for="n in 4" :key="n" class="card"><div class="h-16 skeleton"></div></div>
    </div>

    <div v-else-if="!items.length" class="card">
      <EmptyState title="No shipments" message="Create your first shipment or change your filters.">
        <template #action><button class="btn-primary" @click="modalOpen = true">Create Shipment</button></template>
      </EmptyState>
    </div>

    <div v-else class="space-y-4">
      <div v-for="s in items" :key="s.id" class="card transition-shadow hover:shadow-card-hover">
        <div class="flex items-start justify-between gap-4">
          <div class="flex items-center gap-4">
            <div class="flex h-11 w-11 items-center justify-center rounded-lg" :class="tone(s.status)">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9 17a2 2 0 11-4 0 2 2 0 014 0zm10 0a2 2 0 11-4 0 2 2 0 014 0zM13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1" /></svg>
            </div>
            <div>
              <p class="font-semibold text-slate-900 dark:text-white">Order {{ s.order_id }}</p>
              <div class="mt-1 flex items-center gap-2">
                <span :class="badge(s.status)">{{ (s.status || '').replace('_', ' ') }}</span>
                <span class="text-xs text-slate-400">#{{ s.id.slice(0, 8) }}</span>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button class="btn-ghost btn-sm" @click="toggleTrack(s)">{{ trackingId === s.id ? 'Hide' : 'Track' }}</button>
            <button class="btn-ghost btn-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-500/10" @click="remove(s)">Delete</button>
          </div>
        </div>

        <div class="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-slate-100 pt-4 text-sm dark:border-slate-800">
          <p class="text-slate-500 dark:text-slate-400">
            <span class="font-medium text-slate-700 dark:text-slate-300">{{ s.origin }}</span>
            <span class="mx-2 text-slate-300">→</span>
            <span class="font-medium text-slate-700 dark:text-slate-300">{{ s.destination }}</span>
          </p>
          <select class="input !w-auto !py-1.5 text-xs" :value="s.status" @change="changeStatus(s, $event.target.value)">
            <option v-for="st in statuses" :key="st" :value="st">{{ st.replace('_', ' ') }}</option>
          </select>
        </div>

        <div v-if="trackingId === s.id" class="mt-4 border-t border-slate-100 pt-4 dark:border-slate-800">
          <p class="mb-3 text-sm font-medium text-slate-900 dark:text-white">Tracking history</p>
          <p v-if="!trackEvents.length" class="text-sm text-slate-400">No events yet.</p>
          <ol v-else class="relative ml-2 border-l border-slate-200 dark:border-slate-700">
            <li v-for="ev in trackEvents" :key="ev.id" class="mb-4 ml-4">
              <span class="absolute -left-1 mt-1.5 h-2 w-2 rounded-full bg-primary-500"></span>
              <p class="text-sm font-medium text-slate-900 dark:text-white">{{ ev.description || ev.type }}</p>
              <p class="text-xs text-slate-400">{{ formatDate(ev.created_at) }}<span v-if="ev.location"> · {{ ev.location }}</span></p>
            </li>
          </ol>
        </div>
      </div>

      <div class="card !p-0"><PaginationBar :total="total" :limit="limit" :offset="offset" @change="goTo" /></div>
    </div>

    <ShipmentModal :is-open="modalOpen" @close="modalOpen = false" @submit="handleCreate" />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { shipmentApi } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import ShipmentModal from '@/components/ShipmentModal.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'

const toast = useToastStore()
const statuses = ['pending', 'in_transit', 'delivered', 'cancelled']

const items = ref([])
const total = ref(0)
const limit = ref(20)
const offset = ref(0)
const search = ref('')
const statusFilter = ref('')
const loading = ref(true)
const error = ref(null)
const modalOpen = ref(false)
const trackingId = ref(null)
const trackEvents = ref([])

const formatDate = (d) => (d ? new Date(d).toLocaleString() : '')
const tones = {
  pending: 'bg-amber-100 text-amber-600 dark:bg-amber-500/15 dark:text-amber-400',
  in_transit: 'bg-primary-100 text-primary-600 dark:bg-primary-500/15 dark:text-primary-400',
  delivered: 'bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-400',
  cancelled: 'bg-slate-100 text-slate-500 dark:bg-slate-700/50 dark:text-slate-300',
}
const tone = (s) => tones[s] || tones.pending
const badge = (s) => ({ pending: 'badge-yellow', in_transit: 'badge-blue', delivered: 'badge-green', cancelled: 'badge-gray' }[s] || 'badge-gray')

const load = async () => {
  loading.value = true
  error.value = null
  try {
    const params = { limit: limit.value, offset: offset.value }
    if (search.value.trim()) params.search = search.value.trim()
    if (statusFilter.value) params.status = statusFilter.value
    const res = await shipmentApi.list(params)
    items.value = res.data || []
    total.value = res.total || 0
  } catch (err) {
    error.value = err.message || 'Failed to load shipments'
  } finally {
    loading.value = false
  }
}

const reset = () => {
  offset.value = 0
  load()
}
const goTo = (newOffset) => {
  offset.value = newOffset
  load()
}

let debounce
watch(search, () => {
  clearTimeout(debounce)
  debounce = setTimeout(reset, 300)
})

const handleCreate = async (payload) => {
  try {
    await shipmentApi.create(payload)
    modalOpen.value = false
    toast.success('Shipment created')
    await load()
  } catch (err) {
    toast.error(err.message || 'Failed to create shipment')
  }
}

const changeStatus = async (s, status) => {
  if (status === s.status) return
  try {
    await shipmentApi.updateStatus(s.id, status, '')
    toast.success(`Status updated to ${status.replace('_', ' ')}`)
    await load()
    if (trackingId.value === s.id) await loadTrack(s.id)
  } catch (err) {
    toast.error(err.message || 'Failed to update status')
    await load()
  }
}

const loadTrack = async (id) => {
  const data = await shipmentApi.track(id)
  trackEvents.value = data.events || []
}
const toggleTrack = async (s) => {
  if (trackingId.value === s.id) {
    trackingId.value = null
    trackEvents.value = []
    return
  }
  trackingId.value = s.id
  try {
    await loadTrack(s.id)
  } catch (err) {
    toast.error(err.message || 'Failed to load tracking')
  }
}

const remove = async (s) => {
  if (!confirm(`Delete shipment for order "${s.order_id}"?`)) return
  try {
    await shipmentApi.remove(s.id)
    toast.success('Shipment deleted')
    if (items.value.length === 1 && offset.value > 0) offset.value -= limit.value
    await load()
  } catch (err) {
    toast.error(err.message || 'Failed to delete (admin/manager role required)')
  }
}

onMounted(load)
</script>
