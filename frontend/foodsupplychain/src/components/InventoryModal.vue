<template>
  <div v-if="isOpen" class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity">
    <div class="fixed inset-0 z-10 overflow-y-auto">
      <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
        <div class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
          <div class="absolute right-0 top-0 hidden pr-4 pt-4 sm:block">
            <button
              type="button"
              class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none"
              @click="closeModal"
            >
              <span class="sr-only">Close</span>
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="sm:flex sm:items-start">
            <div class="mt-3 text-center sm:mt-0 sm:text-left w-full">
              <h3 class="text-base font-semibold leading-6 text-gray-900">
                {{ isEditing ? 'Edit Item' : 'Add New Item' }}
              </h3>
              <div class="mt-4 space-y-4">
                <div>
                  <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
                  <input
                    type="text"
                    id="name"
                    v-model="form.name"
                    class="input mt-1 block w-full"
                    :class="{ 'border-red-500': errors.name }"
                  >
                  <p v-if="errors.name" class="mt-1 text-sm text-red-600">{{ errors.name }}</p>
                </div>

                <div>
                  <label for="category" class="block text-sm font-medium text-gray-700">Category</label>
                  <select
                    id="category"
                    v-model="form.category"
                    class="input mt-1 block w-full"
                    :class="{ 'border-red-500': errors.category }"
                  >
                    <option value="">Select Category</option>
                    <option value="fruits">Fruits</option>
                    <option value="vegetables">Vegetables</option>
                    <option value="meat">Meat</option>
                    <option value="dairy">Dairy</option>
                  </select>
                  <p v-if="errors.category" class="mt-1 text-sm text-red-600">{{ errors.category }}</p>
                </div>

                <div>
                  <label for="quantity" class="block text-sm font-medium text-gray-700">Quantity</label>
                  <input
                    type="number"
                    id="quantity"
                    v-model="form.quantity"
                    class="input mt-1 block w-full"
                    :class="{ 'border-red-500': errors.quantity }"
                  >
                  <p v-if="errors.quantity" class="mt-1 text-sm text-red-600">{{ errors.quantity }}</p>
                </div>

                <div>
                  <label for="price" class="block text-sm font-medium text-gray-700">Unit Price</label>
                  <input
                    type="number"
                    id="price"
                    v-model="form.price"
                    step="0.01"
                    class="input mt-1 block w-full"
                    :class="{ 'border-red-500': errors.price }"
                  >
                  <p v-if="errors.price" class="mt-1 text-sm text-red-600">{{ errors.price }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
            <button
              type="button"
              class="btn-primary w-full sm:ml-3 sm:w-auto"
              :disabled="isLoading"
              @click="handleSubmit"
            >
              {{ isLoading ? 'Saving...' : (isEditing ? 'Update' : 'Create') }}
            </button>
            <button
              type="button"
              class="mt-3 w-full sm:mt-0 sm:w-auto btn-secondary"
              @click="closeModal"
              :disabled="isLoading"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true
  },
  isEditing: {
    type: Boolean,
    default: false
  },
  item: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['close', 'submit'])

const isLoading = ref(false)
const errors = reactive({})

const form = reactive({
  name: '',
  category: '',
  quantity: 0,
  price: 0
})

// Watch for item changes when editing
watch(() => props.item, (newItem) => {
  if (newItem && Object.keys(newItem).length) {
    form.name = newItem.name
    form.category = newItem.category
    form.quantity = newItem.quantity
    form.price = newItem.price
  }
}, { immediate: true })

const validateForm = () => {
  errors.name = !form.name ? 'Name is required' : ''
  errors.category = !form.category ? 'Category is required' : ''
  errors.quantity = form.quantity <= 0 ? 'Quantity must be greater than 0' : ''
  errors.price = form.price <= 0 ? 'Price must be greater than 0' : ''

  return !Object.values(errors).some(error => error)
}

const handleSubmit = async () => {
  if (!validateForm()) return

  isLoading.value = true
  try {
    emit('submit', { ...form })
    closeModal()
  } catch (error) {
    console.error('Error submitting form:', error)
  } finally {
    isLoading.value = false
  }
}

const closeModal = () => {
  // Reset form
  form.name = ''
  form.category = ''
  form.quantity = 0
  form.price = 0
  Object.keys(errors).forEach(key => errors[key] = '')
  emit('close')
}
</script>
