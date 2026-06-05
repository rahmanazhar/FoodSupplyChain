<template>
  <div>
    <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl">Catalog</h2>
    <p class="mt-1 text-sm text-gray-500">Manage the products and locations that inventory records reference.</p>

    <div class="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Products -->
      <div class="card">
        <h3 class="text-lg font-medium text-gray-900">Products</h3>
        <form class="mt-4 grid grid-cols-2 gap-3" @submit.prevent="addProduct">
          <input v-model="newProduct.name" class="input" placeholder="Name" />
          <input v-model="newProduct.sku" class="input" placeholder="SKU" />
          <input v-model="newProduct.category" class="input" placeholder="Category" />
          <input v-model.number="newProduct.unit_price" type="number" step="0.01" min="0" class="input" placeholder="Unit price" />
          <div class="col-span-2 flex justify-end">
            <button type="submit" class="btn-primary" :disabled="savingProduct">Add Product</button>
          </div>
        </form>
        <ul class="mt-4 divide-y divide-gray-200">
          <li v-for="p in products" :key="p.id" class="flex items-center justify-between py-3">
            <div>
              <p class="text-sm font-medium text-gray-900">{{ p.name }}</p>
              <p class="text-xs text-gray-400">{{ p.sku }} · {{ p.category || 'uncategorised' }} · ${{ (p.unit_price ?? 0).toFixed(2) }}</p>
            </div>
            <button class="text-red-600 hover:text-red-900 text-sm" @click="removeProduct(p)">Delete</button>
          </li>
          <li v-if="products.length === 0" class="py-3 text-sm text-gray-400">No products yet.</li>
        </ul>
      </div>

      <!-- Locations -->
      <div class="card">
        <h3 class="text-lg font-medium text-gray-900">Locations</h3>
        <form class="mt-4 grid grid-cols-2 gap-3" @submit.prevent="addLocation">
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
        <ul class="mt-4 divide-y divide-gray-200">
          <li v-for="l in locations" :key="l.id" class="flex items-center justify-between py-3">
            <div>
              <p class="text-sm font-medium text-gray-900">{{ l.name }}</p>
              <p class="text-xs text-gray-400">{{ l.type || '—' }}<span v-if="l.city"> · {{ l.city }}</span><span v-if="l.country">, {{ l.country }}</span></p>
            </div>
            <button class="text-red-600 hover:text-red-900 text-sm" @click="removeLocation(l)">Delete</button>
          </li>
          <li v-if="locations.length === 0" class="py-3 text-sm text-gray-400">No locations yet.</li>
        </ul>
      </div>
    </div>

    <div v-if="toast.show" class="fixed bottom-4 right-4 rounded-md p-4 text-white" :class="toast.type === 'success' ? 'bg-green-500' : 'bg-red-500'">
      {{ toast.message }}
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { productApi, locationApi } from '@/services/api'

const products = ref([])
const locations = ref([])
const savingProduct = ref(false)
const savingLocation = ref(false)
const toast = ref({ show: false, message: '', type: 'success' })

const newProduct = reactive({ name: '', sku: '', category: '', unit_price: 0 })
const newLocation = reactive({ name: '', type: '', city: '', country: '' })

const showToast = (message, type = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => (toast.value.show = false), 3000)
}

const load = async () => {
  try {
    const [prods, locs] = await Promise.all([productApi.getAll(), locationApi.getAll()])
    products.value = prods || []
    locations.value = locs || []
  } catch (err) {
    showToast(err.message || 'Failed to load catalog', 'error')
  }
}

const addProduct = async () => {
  if (!newProduct.name || !newProduct.sku) {
    showToast('Name and SKU are required', 'error')
    return
  }
  savingProduct.value = true
  try {
    await productApi.create({ ...newProduct, unit_price: Number(newProduct.unit_price) })
    Object.assign(newProduct, { name: '', sku: '', category: '', unit_price: 0 })
    showToast('Product added')
    await load()
  } catch (err) {
    showToast(err.message || 'Failed to add product', 'error')
  } finally {
    savingProduct.value = false
  }
}

const removeProduct = async (p) => {
  if (!confirm(`Delete product "${p.name}"?`)) return
  try {
    await productApi.remove(p.id)
    showToast('Product deleted')
    await load()
  } catch (err) {
    showToast(err.message || 'Failed to delete (still referenced by inventory?)', 'error')
  }
}

const addLocation = async () => {
  if (!newLocation.name || !newLocation.type) {
    showToast('Name and type are required', 'error')
    return
  }
  savingLocation.value = true
  try {
    await locationApi.create({ ...newLocation })
    Object.assign(newLocation, { name: '', type: '', city: '', country: '' })
    showToast('Location added')
    await load()
  } catch (err) {
    showToast(err.message || 'Failed to add location', 'error')
  } finally {
    savingLocation.value = false
  }
}

const removeLocation = async (l) => {
  if (!confirm(`Delete location "${l.name}"?`)) return
  try {
    await locationApi.remove(l.id)
    showToast('Location deleted')
    await load()
  } catch (err) {
    showToast(err.message || 'Failed to delete (still referenced by inventory?)', 'error')
  }
}

onMounted(load)
</script>
