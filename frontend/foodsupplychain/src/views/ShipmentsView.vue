<template>
  <div>
    <div class="md:flex md:items-center md:justify-between">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
          Shipment Tracking
        </h2>
      </div>
      <div class="mt-4 flex md:ml-4 md:mt-0">
        <button type="button" class="btn-primary" :disabled="!auth.isAuthenticated" @click="openCreate">
          Create Shipment
        </button>
      </div>
    </div>

    <!-- Auth notice -->
    <div v-if="!auth.isAuthenticated" class="mt-8 rounded-md bg-amber-50 p-4 text-sm text-amber-800">
      Shipment endpoints require a JWT. Set an API token from the top bar
      (mint one with <code>go run ./cmd/token -role admin</code>) to manage shipments.
    </div>

    <template v-else>
      <!-- Filters -->
      <div class="mt-8 card">
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1">
            <label for="search" class="block text-sm font-medium text-gray-700">Search</label>
            <input id="search" type="text" class="input mt-1" placeholder="Search by id, order, origin…" v-model="searchQuery" />
          </div>
          <div class="sm:w-48">
            <label for="status" class="block text-sm font-medium text-gray-700">Status</label>
            <select id="status" class="input mt-1" v-model="selectedStatus">
              <option value="">All Status</option>
              <option v-for="s in statuses" :key="s" :value="s">{{ s.replace('_', ' ') }}</option>
            </select>
          </div>
        </div>
      </div>

      <div v-if="isLoading" class="mt-8 flex justify-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>

      <div v-else-if="error" class="mt-8 rounded-md bg-red-50 p-4">
        <h3 class="text-sm font-medium text-red-800">Error loading shipments</h3>
        <p class="mt-2 text-sm text-red-700">{{ error }}</p>
        <button type="button" class="btn-secondary mt-4" @click="load">Try Again</button>
      </div>

      <div v-else class="mt-8 space-y-4">
        <div v-for="shipment in filtered" :key="shipment.id" class="card">
          <div class="flex items-start justify-between">
            <div class="flex items-center">
              <div :class="['h-12 w-12 rounded-full flex items-center justify-center', color(shipment.status).bg]">
                <svg class="h-6 w-6" :class="color(shipment.status).text" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17a2 2 0 11-4 0 2 2 0 014 0zM9 17h6m4 0a2 2 0 11-4 0 2 2 0 014 0zM13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1m8-1a1 1 0 01-1 1H9m4-1V8a1 1 0 011-1h2.586a1 1 0 01.707.293l3.414 3.414a1 1 0 01.293.707V16a1 1 0 01-1 1h-1m-6-1a2 2 0 104 0" />
                </svg>
              </div>
              <div class="ml-4">
                <h3 class="text-lg font-medium">Order {{ shipment.order_id }}</h3>
                <div class="mt-1 flex items-center gap-2">
                  <span :class="['inline-flex rounded-full px-2 text-xs font-semibold leading-5', color(shipment.status).badge]">
                    {{ (shipment.status || '').replace('_', ' ').toUpperCase() }}
                  </span>
                  <span class="text-xs text-gray-400">#{{ shipment.id.slice(0, 8) }}</span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-3 text-sm">
              <button class="text-primary-600 hover:text-primary-900 font-medium" @click="toggleTrack(shipment)">
                {{ trackingId === shipment.id ? 'Hide' : 'Track' }}
              </button>
              <button class="text-red-600 hover:text-red-900 font-medium" @click="remove(shipment)">Delete</button>
            </div>
          </div>

          <div class="mt-4 border-t border-gray-200 pt-4 flex flex-wrap items-center justify-between gap-3 text-sm">
            <div>
              <span class="font-medium text-gray-900">From:</span>
              <span class="ml-1 text-gray-500">{{ shipment.origin }}</span>
              <span class="mx-2 text-gray-300">→</span>
              <span class="font-medium text-gray-900">To:</span>
              <span class="ml-1 text-gray-500">{{ shipment.destination }}</span>
            </div>
            <div class="flex items-center gap-2">
              <select class="input !py-1 !w-auto" :value="shipment.status" @change="changeStatus(shipment, $event.target.value)">
                <option v-for="s in statuses" :key="s" :value="s">{{ s.replace('_', ' ') }}</option>
              </select>
            </div>
          </div>

          <!-- Tracking timeline -->
          <div v-if="trackingId === shipment.id" class="mt-4 border-t border-gray-200 pt-4">
            <h4 class="text-sm font-medium text-gray-900 mb-3">Tracking history</h4>
            <div v-if="trackEvents.length === 0" class="text-sm text-gray-400">No events yet.</div>
            <ol v-else class="relative border-l border-gray-200 ml-2">
              <li v-for="ev in trackEvents" :key="ev.id" class="mb-4 ml-4">
                <div class="absolute w-2 h-2 bg-primary-500 rounded-full -left-1 mt-1.5"></div>
                <p class="text-sm font-medium text-gray-900">{{ ev.description || ev.type }}</p>
                <p class="text-xs text-gray-400">
                  {{ formatDate(ev.created_at) }}<span v-if="ev.location"> · {{ ev.location }}</span>
                </p>
              </li>
            </ol>
          </div>
        </div>

        <div v-if="filtered.length === 0" class="card text-center text-sm text-gray-500">No shipments found.</div>
      </div>
    </template>

    <div v-if="toast.show" class="fixed bottom-4 right-4 rounded-md p-4 text-white" :class="toast.type === 'success' ? 'bg-green-500' : 'bg-red-500'">
      {{ toast.message }}
    </div>

    <ShipmentModal :is-open="isModalOpen" @close="isModalOpen = false" @submit="handleCreate" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { shipmentApi } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import ShipmentModal from '@/components/ShipmentModal.vue'

const auth = useAuthStore()
const statuses = ['pending', 'in_transit', 'delivered', 'cancelled']

const shipments = ref([])
const isLoading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const selectedStatus = ref('')
const isModalOpen = ref(false)
const trackingId = ref(null)
const trackEvents = ref([])
const toast = ref({ show: false, message: '', type: 'success' })

const showToast = (message, type = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => (toast.value.show = false), 3000)
}

const palette = {
  pending: { bg: 'bg-yellow-100', text: 'text-yellow-600', badge: 'bg-yellow-100 text-yellow-800' },
  in_transit: { bg: 'bg-blue-100', text: 'text-blue-600', badge: 'bg-blue-100 text-blue-800' },
  delivered: { bg: 'bg-green-100', text: 'text-green-600', badge: 'bg-green-100 text-green-800' },
  cancelled: { bg: 'bg-gray-100', text: 'text-gray-500', badge: 'bg-gray-100 text-gray-700' }
}
const color = (status) => palette[status] || palette.pending

const formatDate = (d) => (d ? new Date(d).toLocaleString() : '')

const filtered = computed(() =>
  shipments.value.filter((s) => {
    const q = searchQuery.value.toLowerCase()
    const matchesSearch =
      s.id.toLowerCase().includes(q) ||
      (s.order_id || '').toLowerCase().includes(q) ||
      (s.origin || '').toLowerCase().includes(q) ||
      (s.destination || '').toLowerCase().includes(q)
    const matchesStatus = !selectedStatus.value || s.status === selectedStatus.value
    return matchesSearch && matchesStatus
  })
)

const load = async () => {
  if (!auth.isAuthenticated) return
  isLoading.value = true
  error.value = null
  try {
    shipments.value = (await shipmentApi.getAll()) || []
  } catch (err) {
    error.value = err.message || 'Failed to load shipments'
  } finally {
    isLoading.value = false
  }
}

const openCreate = () => (isModalOpen.value = true)

const handleCreate = async (payload) => {
  try {
    await shipmentApi.create(payload)
    isModalOpen.value = false
    showToast('Shipment created')
    await load()
  } catch (err) {
    showToast(err.message || 'Failed to create shipment', 'error')
  }
}

const changeStatus = async (shipment, status) => {
  if (status === shipment.status) return
  try {
    await shipmentApi.updateStatus(shipment.id, status, '')
    showToast(`Status updated to ${status.replace('_', ' ')}`)
    await load()
    if (trackingId.value === shipment.id) await loadTrack(shipment.id)
  } catch (err) {
    showToast(err.message || 'Failed to update status', 'error')
    await load()
  }
}

const loadTrack = async (id) => {
  const data = await shipmentApi.track(id)
  trackEvents.value = data.events || []
}

const toggleTrack = async (shipment) => {
  if (trackingId.value === shipment.id) {
    trackingId.value = null
    trackEvents.value = []
    return
  }
  trackingId.value = shipment.id
  try {
    await loadTrack(shipment.id)
  } catch (err) {
    showToast(err.message || 'Failed to load tracking', 'error')
  }
}

const remove = async (shipment) => {
  if (!confirm(`Delete shipment for order "${shipment.order_id}"?`)) return
  try {
    await shipmentApi.remove(shipment.id)
    showToast('Shipment deleted')
    await load()
  } catch (err) {
    showToast(err.message || 'Failed to delete (admin/manager role required)', 'error')
  }
}

watch(() => auth.token, load)
onMounted(load)
</script>
