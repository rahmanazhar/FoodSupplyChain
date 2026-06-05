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

          <div class="w-full">
            <h3 class="text-base font-semibold leading-6 text-gray-900">
              {{ isEditing ? 'Update Quantity' : 'Add Inventory Item' }}
            </h3>

            <div class="mt-4 space-y-4">
              <!-- Create-only fields -->
              <template v-if="!isEditing">
                <div>
                  <label class="block text-sm font-medium text-gray-700">Product</label>
                  <select v-model="form.product_id" class="input mt-1" :class="{ 'ring-red-500': errors.product_id }">
                    <option value="">Select a product…</option>
                    <option v-for="p in products" :key="p.id" :value="p.id">{{ p.name }} ({{ p.sku }})</option>
                  </select>
                  <p v-if="!products.length" class="mt-1 text-sm text-amber-600">No products yet — add one in Catalog.</p>
                  <p v-if="errors.product_id" class="mt-1 text-sm text-red-600">{{ errors.product_id }}</p>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700">Location</label>
                  <select v-model="form.location_id" class="input mt-1" :class="{ 'ring-red-500': errors.location_id }">
                    <option value="">Select a location…</option>
                    <option v-for="l in locations" :key="l.id" :value="l.id">{{ l.name }}</option>
                  </select>
                  <p v-if="!locations.length" class="mt-1 text-sm text-amber-600">No locations yet — add one in Catalog.</p>
                  <p v-if="errors.location_id" class="mt-1 text-sm text-red-600">{{ errors.location_id }}</p>
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700">Min Quantity</label>
                    <input type="number" v-model.number="form.min_quantity" class="input mt-1" min="0" />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700">Max Quantity</label>
                    <input type="number" v-model.number="form.max_quantity" class="input mt-1" min="0" />
                  </div>
                </div>
              </template>

              <!-- Edit-only: read-only context -->
              <template v-else>
                <div class="rounded-md bg-gray-50 p-3 text-sm text-gray-600">
                  <div><span class="font-medium text-gray-900">Product:</span> {{ item?.product?.name || item?.product_id }}</div>
                  <div><span class="font-medium text-gray-900">Location:</span> {{ item?.location?.name || item?.location_id }}</div>
                  <div><span class="font-medium text-gray-900">Reorder at:</span> {{ item?.min_quantity }}</div>
                </div>
              </template>

              <div>
                <label class="block text-sm font-medium text-gray-700">Quantity</label>
                <input type="number" v-model.number="form.quantity" class="input mt-1" min="0" :class="{ 'ring-red-500': errors.quantity }" />
                <p v-if="errors.quantity" class="mt-1 text-sm text-red-600">{{ errors.quantity }}</p>
              </div>
            </div>
          </div>

          <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
            <button type="button" class="btn-primary w-full sm:ml-3 sm:w-auto" :disabled="isLoading" @click="handleSubmit">
              {{ isLoading ? 'Saving…' : (isEditing ? 'Update' : 'Create') }}
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
  isOpen: { type: Boolean, required: true },
  isEditing: { type: Boolean, default: false },
  item: { type: Object, default: () => ({}) },
  products: { type: Array, default: () => [] },
  locations: { type: Array, default: () => [] }
})

const emit = defineEmits(['close', 'submit'])

const isLoading = ref(false)
const errors = reactive({})

const blank = () => ({ product_id: '', location_id: '', quantity: 0, min_quantity: 0, max_quantity: 0 })
const form = reactive(blank())

watch(
  () => [props.isOpen, props.item],
  () => {
    Object.keys(errors).forEach((k) => (errors[k] = ''))
    if (props.isEditing && props.item) {
      form.product_id = props.item.product_id
      form.location_id = props.item.location_id
      form.quantity = props.item.quantity
      form.min_quantity = props.item.min_quantity
      form.max_quantity = props.item.max_quantity
    } else {
      Object.assign(form, blank())
    }
  },
  { immediate: true }
)

const validate = () => {
  errors.quantity = form.quantity < 0 ? 'Quantity cannot be negative' : ''
  if (!props.isEditing) {
    errors.product_id = !form.product_id ? 'Product is required' : ''
    errors.location_id = !form.location_id ? 'Location is required' : ''
  }
  return !Object.values(errors).some((e) => e)
}

const handleSubmit = async () => {
  if (!validate()) return
  isLoading.value = true
  try {
    if (props.isEditing) {
      emit('submit', { quantity: Number(form.quantity) })
    } else {
      emit('submit', {
        product_id: form.product_id,
        location_id: form.location_id,
        quantity: Number(form.quantity),
        min_quantity: Number(form.min_quantity),
        max_quantity: Number(form.max_quantity)
      })
    }
  } finally {
    isLoading.value = false
  }
}

const closeModal = () => emit('close')
</script>
