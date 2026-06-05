<template>
  <div>
    <div class="md:flex md:items-center md:justify-between">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
          Inventory Management
        </h2>
      </div>
      <div class="mt-4 flex md:ml-4 md:mt-0">
        <button type="button" class="btn-primary" @click="openAddModal">Add Inventory Item</button>
      </div>
    </div>

    <!-- Filters -->
    <div class="mt-8 card">
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <label for="search" class="block text-sm font-medium text-gray-700">Search</label>
          <input id="search" type="text" class="input mt-1" placeholder="Search by product or SKU…" v-model="searchQuery" />
        </div>
        <div class="sm:w-48">
          <label for="category" class="block text-sm font-medium text-gray-700">Category</label>
          <select id="category" class="input mt-1" v-model="selectedCategory">
            <option value="">All Categories</option>
            <option v-for="c in categories" :key="c" :value="c">{{ c }}</option>
          </select>
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="mt-8 flex justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>

    <div v-else-if="error" class="mt-8 rounded-md bg-red-50 p-4">
      <h3 class="text-sm font-medium text-red-800">Error loading inventory</h3>
      <p class="mt-2 text-sm text-red-700">{{ error }}</p>
      <button type="button" class="btn-secondary mt-4" @click="load">Try Again</button>
    </div>

    <div v-else class="mt-8 card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-300">
          <thead>
            <tr>
              <th class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900">Product</th>
              <th class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Category</th>
              <th class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Location</th>
              <th class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Quantity</th>
              <th class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Unit Price</th>
              <th class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
              <th class="relative py-3.5 pl-3 pr-4"><span class="sr-only">Actions</span></th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="item in filtered" :key="item.id">
              <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900">
                {{ item.product?.name || '—' }}
                <span class="block text-xs text-gray-400">{{ item.product?.sku }}</span>
              </td>
              <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ item.product?.category || '—' }}</td>
              <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ item.location?.name || '—' }}</td>
              <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                {{ item.quantity }}
                <span class="text-xs text-gray-400">/ min {{ item.min_quantity }}</span>
              </td>
              <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">${{ (item.product?.unit_price ?? 0).toFixed(2) }}</td>
              <td class="whitespace-nowrap px-3 py-4 text-sm">
                <span :class="['inline-flex rounded-full px-2 text-xs font-semibold leading-5', isLow(item) ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800']">
                  {{ isLow(item) ? 'Low Stock' : 'In Stock' }}
                </span>
              </td>
              <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium">
                <button class="text-primary-600 hover:text-primary-900 mr-2" @click="openEditModal(item)">Adjust</button>
                <button class="text-red-600 hover:text-red-900" @click="remove(item)">Delete</button>
              </td>
            </tr>
            <tr v-if="filtered.length === 0">
              <td colspan="7" class="px-3 py-6 text-sm text-gray-500 text-center">No inventory items found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="toast.show" class="fixed bottom-4 right-4 rounded-md p-4 text-white" :class="toast.type === 'success' ? 'bg-green-500' : 'bg-red-500'">
      {{ toast.message }}
    </div>

    <InventoryModal
      :is-open="isModalOpen"
      :is-editing="!!selectedItem"
      :item="selectedItem"
      :products="products"
      :locations="locations"
      @close="closeModal"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { inventoryApi, productApi, locationApi } from '@/services/api'
import InventoryModal from '@/components/InventoryModal.vue'

const inventory = ref([])
const products = ref([])
const locations = ref([])
const isLoading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const selectedCategory = ref('')
const isModalOpen = ref(false)
const selectedItem = ref(null)
const toast = ref({ show: false, message: '', type: 'success' })

const showToast = (message, type = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => (toast.value.show = false), 3000)
}

const isLow = (item) => item.quantity <= item.min_quantity

const categories = computed(() =>
  [...new Set(products.value.map((p) => p.category).filter(Boolean))].sort()
)

const filtered = computed(() =>
  inventory.value.filter((item) => {
    const q = searchQuery.value.toLowerCase()
    const name = (item.product?.name || '').toLowerCase()
    const sku = (item.product?.sku || '').toLowerCase()
    const matchesSearch = name.includes(q) || sku.includes(q)
    const matchesCategory = !selectedCategory.value || item.product?.category === selectedCategory.value
    return matchesSearch && matchesCategory
  })
)

const load = async () => {
  isLoading.value = true
  error.value = null
  try {
    const [inv, prods, locs] = await Promise.all([
      inventoryApi.getAll(),
      productApi.getAll(),
      locationApi.getAll()
    ])
    inventory.value = inv || []
    products.value = prods || []
    locations.value = locs || []
  } catch (err) {
    error.value = err.message || 'Failed to load inventory'
  } finally {
    isLoading.value = false
  }
}

const openAddModal = () => {
  selectedItem.value = null
  isModalOpen.value = true
}

const openEditModal = (item) => {
  selectedItem.value = item
  isModalOpen.value = true
}

const closeModal = () => {
  selectedItem.value = null
  isModalOpen.value = false
}

const handleSubmit = async (payload) => {
  try {
    if (selectedItem.value) {
      await inventoryApi.updateQuantity(selectedItem.value.id, payload.quantity)
      showToast('Quantity updated')
    } else {
      await inventoryApi.create(payload)
      showToast('Inventory item created')
    }
    closeModal()
    await load()
  } catch (err) {
    showToast(err.message || 'Operation failed', 'error')
  }
}

const remove = async (item) => {
  if (!confirm(`Delete inventory for "${item.product?.name || item.id}"?`)) return
  try {
    await inventoryApi.remove(item.id)
    showToast('Inventory item deleted')
    await load()
  } catch (err) {
    showToast(err.message || 'Failed to delete', 'error')
  }
}

onMounted(load)
</script>
