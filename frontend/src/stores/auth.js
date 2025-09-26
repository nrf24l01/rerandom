import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { refreshAccessToken } from '@/axios'

function isTokenExpired(token) {
  if (!token) return true
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    if (!payload.exp) return false
    return Date.now() >= payload.exp * 1000
  } catch (e) {
    return true
  }
}

export const useAuthStore = defineStore('auth', () => {
  // initialize token from sessionStorage
  const storedToken = sessionStorage.getItem('accessToken')
  const accessToken = ref(storedToken)
  const user_id = ref(null)
  const username = ref(null)

  // decode initial token
  if (storedToken) {
    try {
      const payload = JSON.parse(atob(storedToken.split('.')[1]))
      user_id.value = payload.user_id || null
      username.value = payload.username || null
    } catch (e) {
      user_id.value = null
      username.value = null
    }
  }

  const isAuthenticated = computed(() =>
    !!accessToken.value && !isTokenExpired(accessToken.value)
  )
  const authHeader = computed(() =>
    accessToken.value ? { Authorization: 'Bearer ' + accessToken.value } : {}
  )

  function setToken(token) {
    accessToken.value = token
    if (token) {
      sessionStorage.setItem('accessToken', token)
      // extract user info
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        user_id.value = payload.user_id || null
        username.value = payload.username || null
      } catch (e) {
        user_id.value = null
        username.value = null
      }
    } else {
      sessionStorage.removeItem('accessToken')
      user_id.value = null
      username.value = null
    }
  }

  function logout() {
    setToken(null)
  }

  // Add method to refresh token
  async function refreshToken() {
    try {
      const newToken = await refreshAccessToken()
      setToken(newToken)
      return true
    } catch (error) {
      logout()
      return false
    }
  }

  return {
    accessToken,
    user_id,
    username,
    isAuthenticated,
    authHeader,
    setToken,
    logout,
    refreshToken
  }
})