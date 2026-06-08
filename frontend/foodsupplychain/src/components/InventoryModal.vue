<template>
  <AppModal :open="isOpen" :title="isEditing ? 'Update Quantity' : 'Add Inventory Item'" @close="$emit('close')">
    <div class="space-y-4">
      <template v-if="!isEditing">
        <div>
          <label class="label">Product</label>
          <select v-model="form.product_id" class="input mt-1" :class="{ 'ring-red-500': errors.product_id }">
            <option value="">Select a product…</option>
            <option v-for="p in products" :key="p.id" :value="p.id">{{ p.name }} ({{ p.sku }})</option>
          </select>
          <p v-if="!products.length" class="mt-1 text-xs text-amber-600">No products yet — add one in Catalog.</p>
          <p v-if="errors.product_id" class="mt-1 text-sm text-red-600">{{ errors.product_id }}</p>
        </div>
        <div>
          <label class="label">Location</label>
          <select v-model="form.location_id" class="input mt-1" :class="{ 'ring-red-500': errors.location_id }">
            <option value="">Select a location…</option>
            <option v-for="l in locations" :key="l.id" :value="l.id">{{ l.name }}</option>
          </select>
          <p v-if="!locations.length" class="mt-1 text-xs text-amber-600">No locations yet — add one in Catalog.</p>
          <p v-if="errors.location_id" class="mt-1 text-sm text-red-600">{{ errors.location_id }}</p>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label">Min Quantity</label>
            <input type="number" v-model.number="form.min_quantity" min="0" class="input mt-1" />
          </div>
          <div>
            <label class="label">Max Quantity</label>
            <input type="number" v-model.number="form.max_quantity" min="0" class="input mt-1" />
          </div>
        </div>
      </template>

      <template v-else>
        <div class="rounded-lg bg-slate-50 p-3 text-sm text-slate-600 dark:bg-slate-800/60 dark:text-slate-300">
          <div><span class="font-medium text-slate-900 dark:text-white">Product:</span> {{ item?.product?.name || item?.product_id }}</div>
          <div><span class="font-medium text-slate-900 dark:text-white">Location:</span> {{ item?.location?.name || item?.location_id }}</div>
          <div><span class="font-medium text-slate-900 dark:text-white">Reorder at:</span> {{ item?.min_quantity }}</div>
        </div>
      </template>

      <div>
        <label class="label">Quantity</label>
        <input type="number" v-model.number="form.quantity" min="0" class="input mt-1" :class="{ 'ring-red-500': errors.quantity }" />
        <p v-if="errors.quantity" class="mt-1 text-sm text-red-600">{{ errors.quantity }}</p>
      </div>
    </div>

    <template #footer>
      <button class="btn-secondary" :disabled="isLoading" @click="$emit('close')">Cancel</button>
      <button class="btn-primary" :disabled="isLoading" @click="handleSubmit">
        {{ isLoading ? 'Saving…' : isEditing ? 'Update' : 'Create' }}
      </button>
    </template>
  </AppModal>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import AppModal from '@/components/ui/AppModal.vue'

const props = defineProps({
  isOpen: { type: Boolean, required: true },
  isEditing: { type: Boolean, default: false },
  item: { type: Object, default: () => ({}) },
  products: { type: Array, default: () => [] },
  locations: { type: Array, default: () => [] },
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

const handleSubmit = () => {
  if (!validate()) return
  if (props.isEditing) {
    emit('submit', { quantity: Number(form.quantity) })
  } else {
    emit('submit', {
      product_id: form.product_id,
      location_id: form.location_id,
      quantity: Number(form.quantity),
      min_quantity: Number(form.min_quantity),
      max_quantity: Number(form.max_quantity),
    })
  }
}
</script>
