<template>
  <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
    <!-- Products -->
    <div class="card !p-0">
      <h3 class="border-b border-slate-100 px-5 py-4 text-base font-semibold text-slate-900 dark:border-slate-800 dark:text-white">Products</h3>
      <form class="grid grid-cols-2 gap-3 p-5" @submit.prevent="addProduct">
        <input v-model="newProduct.name" class="input" placeholder="Name" />
        <input v-model="newProduct.sku" class="input" placeholder="SKU" />
        <input v-model="newProduct.category" class="input" placeholder="Category" />
        <input v-model.number="newProduct.unit_price" type="number" step="0.01" min="0" class="input" placeholder="Unit price" />
        <div class="col-span-2 flex justify-end">
          <button type="submit" class="btn-primary" :disabled="savingProduct">Add Product</button>
        </div>
      </form>
      <ul class="divide-y divide-slate-100 dark:divide-slate-800">
        <li v-for="p in products" :key="p.id" class="row-hover flex items-center justify-between px-5 py-3">
          <div>
            <p class="text-sm font-medium text-slate-900 dark:text-white">{{ p.name }}</p>
            <p class="text-xs text-slate-400">{{ p.sku }} · {{ p.category || 'uncategorised' }} · ${{ (p.unit_price ?? 0).toFixed(2) }}</p>
          </div>
          <button class="btn-ghost btn-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-500/10" @click="removeProduct(p)">Delete</button>
        </li>
        <li v-if="!products.length" class="px-5 py-8 text-center text-sm text-slate-400">No products yet.</li>
      </ul>
    </div>

    <!-- Locations -->
    <div class="card !p-0">
      <h3 class="border-b border-slate-100 px-5 py-4 text-base font-semibold text-slate-900 dark:border-slate-800 dark:text-white">Locations</h3>
      <form class="grid grid-cols-2 gap-3 p-5" @submit.prevent="addLocation">
        <input v-model="newLocation.name" class="input" placeholder="Name" />
        <select v-model="newLocation.type" class="input">
          <option value="">Type…</option>
          <option value="warehouse">Warehouse</option>
          <option value="store">Store</option>
          <option value="distribution_center">Distribution Center</option>
        </select>
        <input v-model="newLocation.city" class="input" placeholder="City" />
        <input v-model="newLocation.country" class="input" placeholder="Country" />
        <div class="col-span-2 flex justify-end">
          <button type="submit" class="btn-primary" :disabled="savingLocation">Add Location</button>
        </div>
      </form>
      <ul class="divide-y divide-slate-100 dark:divide-slate-800">
        <li v-for="l in locations" :key="l.id" class="row-hover flex items-center justify-between px-5 py-3">
          <div>
            <p class="text-sm font-medium text-slate-900 dark:text-white">{{ l.name }}</p>
            <p class="text-xs text-slate-400">{{ l.type || '—' }}<span v-if="l.city"> · {{ l.city }}</span><span v-if="l.country">, {{ l.country }}</span></p>
          </div>
          <button class="btn-ghost btn-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-500/10" @click="removeLocation(l)">Delete</button>
        </li>
        <li v-if="!locations.length" class="px-5 py-8 text-center text-sm text-slate-400">No locations yet.</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { productApi, locationApi } from '@/services/api'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()
const products = ref([])
const locations = ref([])
const savingProduct = ref(false)
const savingLocation = ref(false)

const newProduct = reactive({ name: '', sku: '', category: '', unit_price: 0 })
const newLocation = reactive({ name: '', type: '', city: '', country: '' })

const load = async () => {
  try {
    const [p, l] = await Promise.all([productApi.getAll(), locationApi.getAll()])
    products.value = p || []
    locations.value = l || []
  } catch (err) {
    toast.error(err.message || 'Failed to load catalog')
  }
}

const addProduct = async () => {
  if (!newProduct.name || !newProduct.sku) return toast.error('Name and SKU are required')
  savingProduct.value = true
  try {
    await productApi.create({ ...newProduct, unit_price: Number(newProduct.unit_price) })
    Object.assign(newProduct, { name: '', sku: '', category: '', unit_price: 0 })
    toast.success('Product added')
    await load()
  } catch (err) {
    toast.error(err.message || 'Failed to add product')
  } finally {
    savingProduct.value = false
  }
}

const removeProduct = async (p) => {
  if (!confirm(`Delete product "${p.name}"?`)) return
  try {
    await productApi.remove(p.id)
    toast.success('Product deleted')
    await load()
  } catch (err) {
    toast.error(err.message || 'Failed to delete (still referenced by inventory?)')
  }
}

const addLocation = async () => {
  if (!newLocation.name || !newLocation.type) return toast.error('Name and type are required')
  savingLocation.value = true
  try {
    await locationApi.create({ ...newLocation })
    Object.assign(newLocation, { name: '', type: '', city: '', country: '' })
    toast.success('Location added')
    await load()
  } catch (err) {
    toast.error(err.message || 'Failed to add location')
  } finally {
    savingLocation.value = false
  }
}

const removeLocation = async (l) => {
  if (!confirm(`Delete location "${l.name}"?`)) return
  try {
    await locationApi.remove(l.id)
    toast.success('Location deleted')
    await load()
  } catch (err) {
    toast.error(err.message || 'Failed to delete (still referenced by inventory?)')
  }
}

onMounted(load)
</script>
