import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

// Mock the API module so the store doesn't make real HTTP calls.
vi.mock('@/services/api', () => ({
  authApi: {
    login: vi.fn(),
    register: vi.fn(),
  },
}))

import { authApi } from '@/services/api'
import { useAuthStore } from '../auth'

// Build an unsigned JWT-shaped token with the given payload (base64url).
function fakeToken(payload) {
  const b64 = (o) => btoa(JSON.stringify(o)).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
  return `${b64({ alg: 'HS256' })}.${b64(payload)}.sig`
}

describe('auth store', () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('starts unauthenticated', () => {
    const auth = useAuthStore()
    expect(auth.isAuthenticated).toBe(false)
    expect(auth.role).toBe(null)
  })

  it('decodes role and subject from the token claims', () => {
    const auth = useAuthStore()
    auth.setToken(fakeToken({ sub: 'alice', role: 'admin' }))
    expect(auth.isAuthenticated).toBe(true)
    expect(auth.subject).toBe('alice')
    expect(auth.role).toBe('admin')
    expect(localStorage.getItem('fsc.token')).toBeTruthy()
  })

  it('login stores the returned token', async () => {
    authApi.login.mockResolvedValue({ token: fakeToken({ sub: 'bob', role: 'viewer' }) })
    const auth = useAuthStore()
    await auth.login('bob', 'secret')
    expect(authApi.login).toHaveBeenCalledWith('bob', 'secret')
    expect(auth.role).toBe('viewer')
  })

  it('logout clears the session', () => {
    const auth = useAuthStore()
    auth.setToken(fakeToken({ sub: 'x', role: 'viewer' }))
    auth.logout()
    expect(auth.isAuthenticated).toBe(false)
    expect(localStorage.getItem('fsc.token')).toBe(null)
  })
})
