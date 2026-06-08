<template>
  <div class="space-y-6">
    <!-- Stat cards -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <StatCard label="Inventory Items" :value="loading ? '—' : inventoryTotal" tone="primary">
        <template #icon><svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg></template>
      </StatCard>
      <StatCard label="Low Stock" :value="loading ? '—' : lowStock.length" :tone="lowStock.length ? 'red' : 'emerald'">
        <template #icon><svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" /></svg></template>
      </StatCard>
      <StatCard label="Products" :value="loading ? '—' : productCount" tone="slate">
        <template #icon><svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3z" /></svg></template>
      </StatCard>
      <StatCard label="Active Shipments" :value="loading ? '—' : activeShipments" tone="amber">
        <template #icon><svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M9 17a2 2 0 11-4 0 2 2 0 014 0zm10 0a2 2 0 11-4 0 2 2 0 014 0zM13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h1" /></svg></template>
      </StatCard>
    </div>

    <!-- Charts -->
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <div class="card">
        <h3 class="mb-4 text-base font-semibold text-slate-900 dark:text-white">Shipments by Status</h3>
        <DonutChart v-if="!loading" :segments="statusSegments" center-label="shipments" />
        <div v-else class="h-28 skeleton"></div>
      </div>
      <div class="card">
        <h3 class="mb-4 text-base font-semibold text-slate-900 dark:text-white">Stock by Category</h3>
        <BarChart v-if="!loading" :items="categoryBars" />
        <div v-else class="space-y-3"><div v-for="n in 4" :key="n" class="h-6 skeleton"></div></div>
      </div>
    </div>

    <!-- Lists -->
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <div class="card !p-0">
        <h3 class="border-b border-slate-100 px-5 py-4 text-base font-semibold text-slate-900 dark:border-slate-800 dark:text-white">Low Stock Alerts</h3>
        <ul class="divide-y divide-slate-100 dark:divide-slate-800">
          <li v-for="item in lowStock" :key="item.id" class="flex items-center justify-between px-5 py-3">
            <div>
              <p class="text-sm font-medium text-slate-900 dark:text-white">{{ item.product?.name || item.product_id }}</p>
              <p class="text-xs text-slate-400">{{ item.location?.name }}</p>
            </div>
            <span class="badge-red">{{ item.quantity }} / min {{ item.min_quantity }}</span>
          </li>
          <li v-if="!loading && !lowStock.length" class="px-5 py-8 text-center text-sm text-slate-400">Everything is well stocked.</li>
        </ul>
      </div>

      <div class="card !p-0">
        <h3 class="border-b border-slate-100 px-5 py-4 text-base font-semibold text-slate-900 dark:border-slate-800 dark:text-white">Recent Shipments</h3>
        <ul class="divide-y divide-slate-100 dark:divide-slate-800">
          <li v-for="s in recentShipments" :key="s.id" class="flex items-center justify-between px-5 py-3">
            <div>
              <p class="text-sm font-medium text-slate-900 dark:text-white">Order {{ s.order_id }}</p>
              <p class="text-xs text-slate-400">{{ s.origin }} → {{ s.destination }}</p>
            </div>
            <span :class="statusBadge(s.status)">{{ (s.status || '').replace('_', ' ') }}</span>
          </li>
          <li v-if="!loading && !recentShipments.length" class="px-5 py-8 text-center text-sm text-slate-400">No shipments yet.</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { inventoryApi, productApi, shipmentApi } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import StatCard from '@/components/ui/StatCard.vue'
import DonutChart from '@/components/ui/DonutChart.vue'
import BarChart from '@/components/ui/BarChart.vue'

const toast = useToastStore()
const loading = ref(true)
const inventory = ref([])
const inventoryTotal = ref(0)
const productCount = ref(0)
const shipments = ref([])

const statusColors = { pending: '#f59e0b', in_transit: '#0ea5e9', delivered: '#10b981', cancelled: '#94a3b8' }

const lowStock = computed(() => inventory.value.filter((i) => i.quantity <= i.min_quantity))
const statusCounts = computed(() =>
  shipments.value.reduce((acc, s) => ((acc[s.status] = (acc[s.status] || 0) + 1), acc), {})
)
const activeShipments = computed(() => (statusCounts.value.pending || 0) + (statusCounts.value.in_transit || 0))
const statusSegments = computed(() =>
  ['pending', 'in_transit', 'delivered', 'cancelled']
    .map((s) => ({ label: s.replace('_', ' '), value: statusCounts.value[s] || 0, color: statusColors[s] }))
    .filter((s) => s.value > 0)
)
const categoryBars = computed(() => {
  const byCat = {}
  for (const i of inventory.value) {
    const cat = i.product?.category || 'uncategorised'
    byCat[cat] = (byCat[cat] || 0) + i.quantity
  }
  return Object.entries(byCat)
    .map(([label, value]) => ({ label, value }))
    .sort((a, b) => b.value - a.value)
    .slice(0, 6)
})
const recentShipments = computed(() => shipments.value.slice(0, 5))

const statusBadge = (s) =>
  ({ pending: 'badge-yellow', in_transit: 'badge-blue', delivered: 'badge-green', cancelled: 'badge-gray' }[s] || 'badge-gray')

const load = async () => {
  loading.value = true
  try {
    const [inv, prods, ships] = await Promise.all([
      inventoryApi.list({ limit: 100 }),
      productApi.getAll(),
      shipmentApi.list({ limit: 100 })
    ])
    inventory.value = inv.data || []
    inventoryTotal.value = inv.total || 0
    productCount.value = (prods || []).length
    shipments.value = ships.data || []
  } catch (err) {
    toast.error(err.message || 'Failed to load dashboard')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
