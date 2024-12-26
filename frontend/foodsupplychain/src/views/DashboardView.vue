<template>
  <div>
    <div class="md:flex md:items-center md:justify-between">
      <div class="min-w-0 flex-1">
        <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
          Dashboard
        </h2>
      </div>
    </div>

    <!-- Stats -->
    <div class="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
      <!-- Total Inventory -->
      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg class="h-6 w-6 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
            </svg>
          </div>
          <div class="ml-5 w-0 flex-1">
            <dl>
              <dt class="text-sm font-medium text-gray-500 truncate">Total Inventory</dt>
              <dd class="flex items-baseline">
                <div class="text-2xl font-semibold text-gray-900">2,450</div>
                <div class="ml-2 flex items-baseline text-sm font-semibold text-green-600">
                  <svg class="self-center flex-shrink-0 h-5 w-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span class="sr-only">Increased by</span>
                  12%
                </div>
              </dd>
            </dl>
          </div>
        </div>
      </div>

      <!-- Active Shipments -->
      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg class="h-6 w-6 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-5 w-0 flex-1">
            <dl>
              <dt class="text-sm font-medium text-gray-500 truncate">Active Shipments</dt>
              <dd class="flex items-baseline">
                <div class="text-2xl font-semibold text-gray-900">35</div>
                <div class="ml-2 flex items-baseline text-sm font-semibold text-green-600">
                  <svg class="self-center flex-shrink-0 h-5 w-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span class="sr-only">Increased by</span>
                  5%
                </div>
              </dd>
            </dl>
          </div>
        </div>
      </div>

      <!-- Completed Deliveries -->
      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg class="h-6 w-6 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-5 w-0 flex-1">
            <dl>
              <dt class="text-sm font-medium text-gray-500 truncate">Completed Deliveries</dt>
              <dd class="flex items-baseline">
                <div class="text-2xl font-semibold text-gray-900">156</div>
                <div class="ml-2 flex items-baseline text-sm font-semibold text-green-600">
                  <svg class="self-center flex-shrink-0 h-5 w-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span class="sr-only">Increased by</span>
                  8%
                </div>
              </dd>
            </dl>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="mt-8">
      <h3 class="text-lg font-medium leading-6 text-gray-900">Recent Activity</h3>
      <div class="mt-4 card">
        <div class="divide-y divide-gray-200">
          <div v-for="(activity, index) in recentActivities" :key="index" class="py-4 first:pt-0 last:pb-0">
            <div class="flex space-x-3">
              <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                  <h3 class="text-sm font-medium">{{ activity.title }}</h3>
                  <p class="text-sm text-gray-500">{{ activity.time }}</p>
                </div>
                <p class="text-sm text-gray-500">{{ activity.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const recentActivities = ref([
  {
    title: 'New Shipment Created',
    time: '5 minutes ago',
    description: 'Shipment #1234 has been created and is ready for processing'
  },
  {
    title: 'Inventory Updated',
    time: '2 hours ago',
    description: 'Added 500 units of Product A to warehouse inventory'
  },
  {
    title: 'Delivery Completed',
    time: '1 day ago',
    description: 'Shipment #1230 has been successfully delivered to the customer'
  }
])
</script>
