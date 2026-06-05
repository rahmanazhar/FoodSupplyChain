<template>
  <div v-if="isOpen" class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity z-20">
    <div class="fixed inset-0 z-30 overflow-y-auto">
      <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
        <div class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
          <div class="absolute right-0 top-0 hidden pr-4 pt-4 sm:block">
            <button type="button" class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none" @click="closeModal">
              <span class="sr-only">Close</span>
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <h3 class="text-base font-semibold leading-6 text-gray-900">Create Shipment</h3>

          <div class="mt-4 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Order ID</label>
              <input type="text" v-model="form.order_id" class="input mt-1" :class="{ 'ring-red-500': errors.order_id }" />
              <p v-if="errors.order_id" class="mt-1 text-sm text-red-600">{{ errors.order_id }}</p>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Origin</label>
                <input type="text" v-model="form.origin" class="input mt-1" :class="{ 'ring-red-500': errors.origin }" />
                <p v-if="errors.origin" class="mt-1 text-sm text-red-600">{{ errors.origin }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Destination</label>
                <input type="text" v-model="form.destination" class="input mt-1" :class="{ 'ring-red-500': errors.destination }" />
                <p v-if="errors.destination" class="mt-1 text-sm text-red-600">{{ errors.destination }}</p>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Carrier</label>
                <input type="text" v-model="form.carrier_id" class="input mt-1" placeholder="e.g. dhl" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Tracking #</label>
                <input type="text" v-model="form.tracking_number" class="input mt-1" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Status</label>
                <select v-model="form.status" class="input mt-1">
                  <option v-for="s in statuses" :key="s" :value="s">{{ s.replace('_', ' ') }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Est. Arrival</label>
                <input type="datetime-local" v-model="form.estimated_arrival" class="input mt-1" />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">Notes</label>
              <textarea v-model="form.notes" rows="2" class="input mt-1"></textarea>
            </div>
          </div>

          <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
            <button type="button" class="btn-primary w-full sm:ml-3 sm:w-auto" :disabled="isLoading" @click="handleSubmit">
              {{ isLoading ? 'Saving…' : 'Create' }}
            </button>
            <button type="button" class="mt-3 w-full sm:mt-0 sm:w-auto btn-secondary" :disabled="isLoading" @click="closeModal">
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'

const props = defineProps({
  isOpen: { type: Boolean, required: true }
})
const emit = defineEmits(['close', 'submit'])

const statuses = ['pending', 'in_transit', 'delivered', 'cancelled']
const isLoading = ref(false)
const errors = reactive({})

const blank = () => ({
  order_id: '',
  origin: '',
  destination: '',
  carrier_id: '',
  tracking_number: '',
  status: 'pending',
  estimated_arrival: '',
  notes: ''
})
const form = reactive(blank())

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      Object.assign(form, blank())
      Object.keys(errors).forEach((k) => (errors[k] = ''))
    }
  }
)

const validate = () => {
  errors.order_id = !form.order_id ? 'Order ID is required' : ''
  errors.origin = !form.origin ? 'Origin is required' : ''
  errors.destination = !form.destination ? 'Destination is required' : ''
  return !Object.values(errors).some((e) => e)
}

const handleSubmit = async () => {
  if (!validate()) return
  isLoading.value = true
  try {
    const payload = { ...form }
    if (payload.estimated_arrival) {
      payload.estimated_arrival = new Date(payload.estimated_arrival).toISOString()
    } else {
      delete payload.estimated_arrival
    }
    emit('submit', payload)
  } finally {
    isLoading.value = false
  }
}

const closeModal = () => emit('close')
</script>
