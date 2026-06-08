<template>
  <AppModal :open="isOpen" title="Create Shipment" @close="$emit('close')">
    <div class="space-y-4">
      <div>
        <label class="label">Order ID</label>
        <input v-model="form.order_id" type="text" class="input mt-1" :class="{ 'ring-red-500': errors.order_id }" />
        <p v-if="errors.order_id" class="mt-1 text-sm text-red-600">{{ errors.order_id }}</p>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Origin</label>
          <input v-model="form.origin" type="text" class="input mt-1" :class="{ 'ring-red-500': errors.origin }" />
        </div>
        <div>
          <label class="label">Destination</label>
          <input v-model="form.destination" type="text" class="input mt-1" :class="{ 'ring-red-500': errors.destination }" />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Carrier</label>
          <input v-model="form.carrier_id" type="text" class="input mt-1" placeholder="e.g. dhl" />
        </div>
        <div>
          <label class="label">Tracking #</label>
          <input v-model="form.tracking_number" type="text" class="input mt-1" />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Status</label>
          <select v-model="form.status" class="input mt-1">
            <option v-for="s in statuses" :key="s" :value="s">{{ s.replace('_', ' ') }}</option>
          </select>
        </div>
        <div>
          <label class="label">Est. Arrival</label>
          <input v-model="form.estimated_arrival" type="datetime-local" class="input mt-1" />
        </div>
      </div>
      <div>
        <label class="label">Notes</label>
        <textarea v-model="form.notes" rows="2" class="input mt-1"></textarea>
      </div>
    </div>

    <template #footer>
      <button class="btn-secondary" :disabled="isLoading" @click="$emit('close')">Cancel</button>
      <button class="btn-primary" :disabled="isLoading" @click="handleSubmit">{{ isLoading ? 'Saving…' : 'Create' }}</button>
    </template>
  </AppModal>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import AppModal from '@/components/ui/AppModal.vue'

const props = defineProps({ isOpen: { type: Boolean, required: true } })
const emit = defineEmits(['close', 'submit'])

const statuses = ['pending', 'in_transit', 'delivered', 'cancelled']
const isLoading = ref(false)
const errors = reactive({})
const blank = () => ({ order_id: '', origin: '', destination: '', carrier_id: '', tracking_number: '', status: 'pending', estimated_arrival: '', notes: '' })
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

const handleSubmit = () => {
  if (!validate()) return
  const payload = { ...form }
  if (payload.estimated_arrival) payload.estimated_arrival = new Date(payload.estimated_arrival).toISOString()
  else delete payload.estimated_arrival
  emit('submit', payload)
}
</script>
