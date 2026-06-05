<template>
  <div>
    <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl">Dashboard</h2>

    <!-- Stats -->
    <div class="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
      <div class="card">
        <dt class="text-sm font-medium text-gray-500 truncate">Inventory Items</dt>
        <dd class="mt-1 text-2xl font-semibold text-gray-900">{{ inventory.length }}</dd>
      </div>
      <div class="card">
        <dt class="text-sm font-medium text-gray-500 truncate">Low Stock</dt>
        <dd class="mt-1 text-2xl font-semibold" :class="lowStock.length ? 'text-red-600' : 'text-gray-900'">{{ lowStock.length }}</dd>
      </div>
      <div class="card">
        <dt class="text-sm font-medium text-gray-500 truncate">Products</dt>
        <dd class="mt-1 text-2xl font-semibold text-gray-900">{{ productCount }}</dd>
      </div>
      <div class="card">
        <dt class="text-sm font-medium text-gray-500 truncate">Active Shipments</dt>
        <dd class="mt-1 text-2xl font-semibold text-gray-900">{{ auth.isAuthenticated ? activeShipments : '—' }}</dd>
      </div>
    </div>

    <div class="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Low stock list -->
      <div class="card">
        <h3 class="text-lg font-medium text-gray-900">Low Stock Alerts</h3>
        <ul class="mt-4 divide-y divide-gray-200">
          <li v-for="item in lowStock" :key="item.id" class="flex items-center justify-between py-3">
            <div>
              <p class="text-sm font-medium text-gray-900">{{ item.product?.name || item.product_id }}</p>
              <p class="text-xs text-gray-400">{{ item.location?.name }}</p>
            </div>
            <span class="inline-flex rounded-full bg-red-100 px-2 text-xs font-semibold leading-5 text-red-800">
              {{ item.quantity }} / min {{ item.min_quantity }}
            </span>
          </li>
          <li v-if="lowStock.length === 0" class="py-3 text-sm text-gray-400">Everything is well stocked.</li>
        </ul>
      </div>

      <!-- Shipments by status -->
      <div class="card">
        <h3 class="text-lg font-medium text-gray-900">Shipments by Status</h3>
        <div v-if="!auth.isAuthenticated" class="mt-4 text-sm text-amber-700">
          Set an API token to view shipment metrics.
        </div>
        <ul v-else class="mt-4 divide-y divide-gray-200">
          <li v-for="s in statuses" :key="s" class="flex items-center justify-between py-3">
            <span class="text-sm text-gray-700 capitalize">{{ s.replace('_', ' ') }}</span>
            <span class="text-sm font-semibold text-gray-900">{{ statusCounts[s] || 0 }}</span>
          </li>
        </ul>
      </div>
    </div>

    <p v-if="error" class="mt-6 text-sm text-red-600">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { inventoryApi, productApi, shipmentApi } from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const statuses = ['pending', 'in_transit', 'delivered', 'cancelled']

const inventory = ref([])
const productCount = ref(0)
const shipments = ref([])
const error = ref(null)

const lowStock = computed(() => inventory.value.filter((i) => i.quantity <= i.min_quantity))
const statusCounts = computed(() =>
  shipments.value.reduce((acc, s) => {
    acc[s.status] = (acc[s.status] || 0) + 1
    return acc
  }, {})
)
const activeShipments = computed(
  () => (statusCounts.value.pending || 0) + (statusCounts.value.in_transit || 0)
)

const load = async () => {
  error.value = null
  try {
    const [inv, prods] = await Promise.all([inventoryApi.getAll(), productApi.getAll()])
    inventory.value = inv || []
    productCount.value = (prods || []).length
  } catch (err) {
    error.value = err.message || 'Failed to load dashboard data'
  }
  if (auth.isAuthenticated) {
    try {
      shipments.value = (await shipmentApi.getAll()) || []
    } catch {
      shipments.value = []
    }
  } else {
    shipments.value = []
  }
}

watch(() => auth.token, load)
onMounted(load)
</script>
