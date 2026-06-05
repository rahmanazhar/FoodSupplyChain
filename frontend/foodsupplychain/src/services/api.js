import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

// All requests go through the API gateway, which proxies to the inventory and
// shipment services.
const api = axios.create({
  baseURL: 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Attach the bearer token (required by the shipment service's /api/v1 routes).
api.interceptors.request.use(
  (config) => {
    const auth = useAuthStore()
    if (auth.token) {
      config.headers.Authorization = `Bearer ${auth.token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Normalise error responses to a single readable Error.
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const data = error.response?.data
    const message = data?.error || data?.message || error.message || 'Request failed'
    const err = new Error(message)
    err.status = error.response?.status
    return Promise.reject(err)
  }
)

export const inventoryApi = {
  getAll: () => api.get('/inventory').then((r) => r.data),
  get: (id) => api.get(`/inventory/${id}`).then((r) => r.data),
  create: (data) => api.post('/inventory', data).then((r) => r.data),
  updateQuantity: (id, quantity) => api.put(`/inventory/${id}`, { quantity }).then((r) => r.data),
  remove: (id) => api.delete(`/inventory/${id}`).then((r) => r.data)
}

export const productApi = {
  getAll: () => api.get('/products').then((r) => r.data),
  create: (data) => api.post('/products', data).then((r) => r.data),
  remove: (id) => api.delete(`/products/${id}`).then((r) => r.data)
}

export const locationApi = {
  getAll: () => api.get('/locations').then((r) => r.data),
  create: (data) => api.post('/locations', data).then((r) => r.data),
  remove: (id) => api.delete(`/locations/${id}`).then((r) => r.data)
}

export const authApi = {
  login: (username, password) => api.post('/auth/login', { username, password }).then((r) => r.data),
  register: (payload) => api.post('/auth/register', payload).then((r) => r.data),
  me: () => api.get('/auth/me').then((r) => r.data)
}

export const shipmentApi = {
  getAll: () => api.get('/shipments').then((r) => r.data),
  get: (id) => api.get(`/shipments/${id}`).then((r) => r.data),
  create: (data) => api.post('/shipments', data).then((r) => r.data),
  update: (id, data) => api.put(`/shipments/${id}`, data).then((r) => r.data),
  updateStatus: (id, status, location) =>
    api.put(`/shipments/${id}/status`, { status, location }).then((r) => r.data),
  track: (id) => api.get(`/shipments/${id}/track`).then((r) => r.data),
  remove: (id) => api.delete(`/shipments/${id}`).then((r) => r.data)
}

export default api
