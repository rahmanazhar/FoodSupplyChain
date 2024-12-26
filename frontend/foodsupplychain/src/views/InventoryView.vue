<template>
  <div>
    <div class="md:flex md:items-center md:justify-between">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
          Inventory Management
        </h2>
      </div>
      <div class="mt-4 flex md:ml-4 md:mt-0">
        <button type="button" class="btn-primary" @click="openAddModal">
          Add New Item
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="mt-8 card">
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <label for="search" class="block text-sm font-medium text-gray-700">Search</label>
          <div class="mt-1">
            <input
              type="text"
              name="search"
              id="search"
              class="input"
              placeholder="Search inventory..."
              v-model="searchQuery"
            >
          </div>
        </div>
        <div class="sm:w-48">
          <label for="category" class="block text-sm font-medium text-gray-700">Category</label>
          <select
            id="category"
            name="category"
            class="input"
            v-model="selectedCategory"
          >
            <option value="">All Categories</option>
            <option value="fruits">Fruits</option>
            <option value="vegetables">Vegetables</option>
            <option value="meat">Meat</option>
            <option value="dairy">Dairy</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="mt-8 flex justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="mt-8">
      <div class="rounded-md bg-red-50 p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error loading inventory</h3>
            <div class="mt-2 text-sm text-red-700">
              {{ error }}
            </div>
            <div class="mt-4">
              <button
                type="button"
                class="btn-secondary"
                @click="fetchInventory"
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Inventory Table -->
    <div v-else class="mt-8 card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-300">
          <thead>
            <tr>
              <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900">Item Name</th>
              <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Category</th>
              <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Quantity</th>
              <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Unit Price</th>
              <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
              <th scope="col" class="relative py-3.5 pl-3 pr-4">
                <span class="sr-only">Actions</span>
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="item in filteredInventory" :key="item.id">
              <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900">{{ item.name }}</td>
              <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ item.category }}</td>
              <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ item.quantity }}</td>
              <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">${{ item.price.toFixed(2) }}</td>
              <td class="whitespace-nowrap px-3 py-4 text-sm">
                <span :class="[
                  'inline-flex rounded-full px-2 text-xs font-semibold leading-5',
                  item.quantity > 10 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                ]">
                  {{ item.quantity > 10 ? 'In Stock' : 'Low Stock' }}
                </span>
              </td>
              <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium">
                <button 
                  class="text-primary-600 hover:text-primary-900 mr-2"
                  @click="openEditModal(item)"
                >
                  Edit
                </button>
                <button 
                  class="text-red-600 hover:text-red-900"
                  @click="confirmDelete(item)"
                >
                  Delete
                </button>
              </td>
            </tr>
            <tr v-if="filteredInventory.length === 0">
              <td colspan="6" class="px-3 py-4 text-sm text-gray-500 text-center">
                No items found
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Toast Notification -->
    <div
      v-if="toast.show"
      class="fixed bottom-4 right-4 rounded-md p-4 text-white"
      :class="{
        'bg-green-500': toast.type === 'success',
        'bg-red-500': toast.type === 'error'
      }"
    >
      {{ toast.message }}
    </div>

    <!-- Inventory Modal -->
    <InventoryModal
      :is-open="isModalOpen"
      :is-editing="!!selectedItem"
      :item="selectedItem"
      @close="closeModal"
      @submit="handleSubmit"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { inventoryApi } from '@/services/api'
import InventoryModal from '@/components/InventoryModal.vue'

const inventory = ref([])
const isLoading = ref(false)
const error = ref(null)
const searchQuery = ref('')
const selectedCategory = ref('')
const isModalOpen = ref(false)
const selectedItem = ref(null)
const toast = ref({
  show: false,
  message: '',
  type: 'success'
})

const showToast = (message, type = 'success') => {
  toast.value = {
    show: true,
    message,
    type
  }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

const filteredInventory = computed(() => {
  return inventory.value.filter(item => {
    const matchesSearch = item.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesCategory = !selectedCategory.value || item.category === selectedCategory.value
    return matchesSearch && matchesCategory
  })
})

const fetchInventory = async () => {
  isLoading.value = true
  error.value = null
  try {
    const data = await inventoryApi.getAll()
    inventory.value = data
  } catch (err) {
    error.value = err.message || 'Failed to load inventory'
    showToast('Failed to load inventory', 'error')
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

const handleSubmit = async (formData) => {
  try {
    if (selectedItem.value) {
      await inventoryApi.update(selectedItem.value.id, formData)
      showToast('Item updated successfully')
    } else {
      await inventoryApi.create(formData)
      showToast('Item created successfully')
    }
    await fetchInventory()
  } catch (err) {
    showToast(err.message || 'Operation failed', 'error')
  }
}

const confirmDelete = async (item) => {
  if (confirm('Are you sure you want to delete this item?')) {
    try {
      await inventoryApi.delete(item.id)
      showToast('Item deleted successfully')
      await fetchInventory()
    } catch (err) {
      showToast(err.message || 'Failed to delete item', 'error')
    }
  }
}

onMounted(() => {
  fetchInventory()
})
</script>
