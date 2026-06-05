import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const STORAGE_KEY = 'fsc.token'

// Decode a base64url segment to a UTF-8 string.
function b64urlDecode(segment) {
  let s = segment.replace(/-/g, '+').replace(/_/g, '/')
  while (s.length % 4) s += '='
  return atob(s)
}

// useAuthStore holds the JWT used for the shipment service's protected routes.
// The token is minted out-of-band (`go run ./cmd/token -role admin`) and pasted
// into the UI; it is persisted to localStorage so it survives reloads.
export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(STORAGE_KEY) || '')

  const isAuthenticated = computed(() => !!token.value)

  // Decode the JWT payload for display only (no signature verification).
  const claims = computed(() => {
    if (!token.value) return null
    try {
      return JSON.parse(b64urlDecode(token.value.split('.')[1]))
    } catch {
      return null
    }
  })

  const role = computed(() => claims.value?.role || null)
  const subject = computed(() => claims.value?.sub || null)

  function setToken(value) {
    token.value = (value || '').trim()
    if (token.value) {
      localStorage.setItem(STORAGE_KEY, token.value)
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  function clear() {
    setToken('')
  }

  return { token, isAuthenticated, claims, role, subject, setToken, clear }
})
