import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi } from '@/lib/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(email: string, password: string) {
    const response = await authApi.login({ email, password })
    const data = response.data.data!
    
    token.value = data.access_token
    refreshToken.value = data.refresh_token
    user.value = data.user

    localStorage.setItem('access_token', data.access_token)
    localStorage.setItem('refresh_token', data.refresh_token)

    return data
  }

  async function register(email: string, password: string, username: string) {
    const response = await authApi.register({ email, password, username })
    return response.data.data
  }

  async function fetchUser() {
    if (!token.value) return null

    try {
      const response = await authApi.me()
      user.value = response.data.data!
      return user.value
    } catch {
      logout()
      return null
    }
  }

  function logout() {
    user.value = null
    token.value = null
    refreshToken.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  async function forgotPassword(email: string) {
    const response = await authApi.forgotPassword(email)
    return response.data
  }

  async function resetPassword(token: string, password: string) {
    const response = await authApi.resetPassword(token, password)
    return response.data
  }

  return {
    user,
    token,
    refreshToken,
    isAuthenticated,
    isAdmin,
    login,
    register,
    fetchUser,
    logout,
    forgotPassword,
    resetPassword,
  }
})
