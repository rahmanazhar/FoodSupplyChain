import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:3000',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor for API calls
api.interceptors.request.use(
  (config) => {
    // You can add auth token here later
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor for API calls
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    if (error.response) {
      if (error.response.status === 401 && !originalRequest._retry) {
        // Handle token refresh or redirect to login here
      }
      return Promise.reject(error.response.data)
    }
    return Promise.reject(error)
  }
)

export const inventoryApi = {
  // Get all inventory items
  getAll: async (params = {}) => {
    const response = await api.get('/inventory', { params })
    return response.data
  },

  // Get a single inventory item
  getById: async (id) => {
    const response = await api.get(`/inventory/${id}`)
    return response.data
  },

  // Create a new inventory item
  create: async (data) => {
    const response = await api.post('/inventory', data)
    return response.data
  },

  // Update an inventory item
  update: async (id, data) => {
    const response = await api.put(`/inventory/${id}`, data)
    return response.data
  },

  // Delete an inventory item
  delete: async (id) => {
    const response = await api.delete(`/inventory/${id}`)
    return response.data
  }
}

export default api
