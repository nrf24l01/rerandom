<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6 text-center">Login</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block text-gray-700 mb-2" for="username">Username</label>
          <input id="username" v-model="username" type="text" class="w-full px-3 py-2 border rounded" required />
        </div>
        <div class="mb-4">
          <label class="block text-gray-700 mb-2" for="password">Password</label>
          <input id="password" v-model="password" type="password" class="w-full px-3 py-2 border rounded" required />
        </div>
        <div v-if="error" class="mb-4 text-red-500 text-sm">{{ error }}</div>
        <button type="submit" class="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600" :disabled="loading">
          <span v-if="loading">Logging in...</span>
          <span v-else>Login</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const username = ref('')
const password = ref('')
const error = ref(null)
const loading = ref(false)

const router = useRouter()
const auth = useAuthStore()

async function handleLogin() {
  loading.value = true
  error.value = null
  try {
    const response = await axios.post(
      import.meta.env.VITE_BACKEND_URL + '/auth/login',
      { username: username.value, password: password.value },
      { withCredentials: true }
    )
    const token = response.data.access_token
    auth.setToken(token)
    router.push('/')
  } catch (err) {
    if (err.response && err.response.data && err.response.data.message) {
      error.value = err.response.data.message
    } else {
      error.value = 'An unexpected error occurred'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* Optional custom styles */
</style>