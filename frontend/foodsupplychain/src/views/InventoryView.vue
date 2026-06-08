<template>
  <div class="space-y-5">
    <!-- Toolbar -->
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="relative w-full sm:max-w-xs">
        <svg class="pointer-events-none absolute left-3 top-2.5 h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke-width="1.8" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="m21 21-4.3-4.3m1.8-4.7a6.5 6.5 0 1 1-13 0 6.5 6.5 0 0 1 13 0z" /></svg>
        <input v-model="search" type="text" class="input pl-9" placeholder="Search products or SKU…" />
      </div>
      <button class="btn-primary" @click="openAdd">
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" /></svg>
        Add Item
      </button>
    </div>

    <!-- Table -->
    <div class="card !p-0 overflow-hidden">
      <div v-if="error" class="px-5 py-10 text-center">
        <p class="text-sm text-red-600">{{ error }}</p>
        <button class="btn-secondary mt-3" @click="load">Try again</button>
      </div>

      <SkeletonRows v-else-if="loading" :rows="6" :cols="6" />

      <EmptyState v-else-if="!items.length" title="No inventory items" message="Add your first item, or adjust your search." >
        <template #action><button class="btn-primary" @click="openAdd">Add Item</button></template>
      </EmptyState>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full">
          <thead class="border-b border-slate-100 dark:border-slate-800">
            <tr>
              <th class="th">Product</th>
              <th class="th">Category</th>
              <th class="th">Location</th>
              <th class="th">Quantity</th>
              <th class="th">Unit Price</th>
              <th class="th">Status</th>
              <th class="th text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-800">
            <tr v-for="item in items" :key="item.id" class="row-hover">
              <td class="td font-medium text-slate-900 dark:text-white">
                {{ item.product?.name || '—' }}
                <span class="block text-xs font-normal text-slate-400">{{ item.product?.sku }}</span>
              </td>
              <td class="td capitalize">{{ item.product?.category || '—' }}</td>
              <td class="td">{{ item.location?.name || '—' }}</td>
              <td class="td tabular-nums">{{ item.quantity }} <span class="text-xs text-slate-400">/ min {{ item.min_quantity }}</span></td>
              <td class="td tabular-nums">${{ (item.product?.unit_price ?? 0).toFixed(2) }}</td>
              <td class="td">
                <span :class="isLow(item) ? 'badge-red' : 'badge-green'">{{ isLow(item) ? 'Low Stock' : 'In Stock' }}</span>
              </td>
              <td class="td text-right">
                <button class="btn-ghost btn-sm" @click="openEdit(item)">Adjust</button>
                <button class="btn-ghost btn-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-500/10" @click="remove(item)">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
        <PaginationBar :total="total" :limit="limit" :offset="offset" @change="goTo" />
      </div>
    </div>

    <InventoryModal
      :is-open="modalOpen"
      :is-editing="!!selected"
      :item="selected"
      :products="products"
      :locations="locations"
      @close="closeModal"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { inventoryApi, productApi, locationApi } from '@/services/api'
import { useToastStore } from '@/stores/toast'
import InventoryModal from '@/components/InventoryModal.vue'
import SkeletonRows from '@/components/ui/SkeletonRows.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import PaginationBar from '@/components/ui/PaginationBar.vue'

const toast = useToastStore()
const items = ref([])
const products = ref([])
const locations = ref([])
const total = ref(0)
const limit = ref(20)
const offset = ref(0)
const search = ref('')
const loading = ref(true)
const error = ref(null)
const modalOpen = ref(false)
const selected = ref(null)

const isLow = (item) => item.quantity <= item.min_quantity

const load = async () => {
  loading.value = true
  error.value = null
  try {
    const params = { limit: limit.value, offset: offset.value }
    if (search.value.trim()) params.search = search.value.trim()
    const res = await inventoryApi.list(params)
    items.value = res.data || []
    total.value = res.total || 0
  } catch (err) {
    error.value = err.message || 'Failed to load inventory'
  } finally {
    loading.value = false
  }
}

const loadRefs = async () => {
  try {
    const [p, l] = await Promise.all([productApi.getAll(), locationApi.getAll()])
    products.value = p || []
    locations.value = l || []
  } catch {
    /* dropdown data is best-effort */
  }
}

const goTo = (newOffset) => {
  offset.value = newOffset
  load()
}

let debounce
watch(search, () => {
  clearTimeout(debounce)
  debounce = setTimeout(() => {
    offset.value = 0
    load()
  }, 300)
})

const openAdd = () => {
  selected.value = null
  modalOpen.value = true
}
const openEdit = (item) => {
  selected.value = item
  modalOpen.value = true
}
const closeModal = () => {
  selected.value = null
  modalOpen.value = false
}

const handleSubmit = async (payload) => {
  try {
    if (selected.value) {
      await inventoryApi.updateQuantity(selected.value.id, payload.quantity)
      toast.success('Quantity updated')
    } else {
      await inventoryApi.create(payload)
      toast.success('Inventory item created')
    }
    closeModal()
    await load()
  } catch (err) {
    toast.error(err.message || 'Operation failed')
  }
}

const remove = async (item) => {
  if (!confirm(`Delete inventory for "${item.product?.name || item.id}"?`)) return
  try {
    await inventoryApi.remove(item.id)
    toast.success('Inventory item deleted')
    if (items.value.length === 1 && offset.value > 0) offset.value -= limit.value
    await load()
  } catch (err) {
    toast.error(err.message || 'Failed to delete')
  }
}

onMounted(() => {
  load()
  loadRefs()
})
</script>
