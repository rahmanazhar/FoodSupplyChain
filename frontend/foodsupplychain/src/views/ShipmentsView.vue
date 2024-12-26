<template>
  <div>
    <div class="md:flex md:items-center md:justify-between">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
          Shipment Tracking
        </h2>
      </div>
      <div class="mt-4 flex md:ml-4 md:mt-0">
        <button type="button" class="btn-primary">
          Create Shipment
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
              placeholder="Search shipments..."
              v-model="searchQuery"
            >
          </div>
        </div>
        <div class="sm:w-48">
          <label for="status" class="block text-sm font-medium text-gray-700">Status</label>
          <select
            id="status"
            name="status"
            class="input"
            v-model="selectedStatus"
          >
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="in_transit">In Transit</option>
            <option value="delivered">Delivered</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Shipments List -->
    <div class="mt-8">
      <div class="space-y-4">
        <div v-for="shipment in filteredShipments" :key="shipment.id" class="card">
          <div class="flex items-center justify-between">
            <div class="flex-1">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <div :class="[
                    'h-12 w-12 rounded-full flex items-center justify-center',
                    statusColors[shipment.status].bg
                  ]">
                    <svg class="h-6 w-6" :class="statusColors[shipment.status].text" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                </div>
                <div class="ml-4">
                  <h3 class="text-lg font-medium">Shipment #{{ shipment.id }}</h3>
                  <div class="mt-1 flex items-center">
                    <span :class="[
                      'inline-flex rounded-full px-2 text-xs font-semibold leading-5',
                      statusColors[shipment.status].badge
                    ]">
                      {{ shipment.status.replace('_', ' ').toUpperCase() }}
                    </span>
                    <span class="ml-2 text-sm text-gray-500">{{ shipment.date }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="mt-4 sm:mt-0">
              <div class="flex -space-x-1 overflow-hidden">
                <div v-for="(item, index) in shipment.items" :key="index" class="inline-block">
                  <span class="relative inline-flex h-8 w-8 items-center justify-center rounded-full bg-gray-100 text-xs font-medium">
                    {{ item.quantity }}
                  </span>
                </div>
              </div>
            </div>
          </div>
          <div class="mt-4 border-t border-gray-200 pt-4">
            <div class="flex justify-between text-sm">
              <div>
                <span class="font-medium text-gray-900">From:</span>
                <span class="ml-1 text-gray-500">{{ shipment.origin }}</span>
              </div>
              <div>
                <span class="font-medium text-gray-900">To:</span>
                <span class="ml-1 text-gray-500">{{ shipment.destination }}</span>
              </div>
              <button class="text-primary-600 hover:text-primary-900 font-medium">
                View Details
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const searchQuery = ref('')
const selectedStatus = ref('')

const statusColors = {
  pending: {
    bg: 'bg-yellow-100',
    text: 'text-yellow-600',
    badge: 'bg-yellow-100 text-yellow-800'
  },
  in_transit: {
    bg: 'bg-blue-100',
    text: 'text-blue-600',
    badge: 'bg-blue-100 text-blue-800'
  },
  delivered: {
    bg: 'bg-green-100',
    text: 'text-green-600',
    badge: 'bg-green-100 text-green-800'
  }
}

const shipments = ref([
  {
    id: '1234',
    status: 'in_transit',
    date: '2024-01-15',
    origin: 'Warehouse A',
    destination: 'Store B',
    items: [
      { quantity: 5 },
      { quantity: 3 },
      { quantity: 2 }
    ]
  },
  {
    id: '1235',
    status: 'pending',
    date: '2024-01-16',
    origin: 'Supplier X',
    destination: 'Warehouse A',
    items: [
      { quantity: 10 },
      { quantity: 8 }
    ]
  },
  {
    id: '1236',
    status: 'delivered',
    date: '2024-01-14',
    origin: 'Warehouse B',
    destination: 'Store C',
    items: [
      { quantity: 4 },
      { quantity: 6 },
      { quantity: 2 }
    ]
  }
])

const filteredShipments = computed(() => {
  return shipments.value.filter(shipment => {
    const matchesSearch = shipment.id.includes(searchQuery.value) ||
      shipment.origin.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      shipment.destination.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = !selectedStatus.value || shipment.status === selectedStatus.value
    return matchesSearch && matchesStatus
  })
})
</script>
