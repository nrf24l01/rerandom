<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded shadow-lg w-full max-w-md">
      <h3 class="text-lg font-bold mb-4">Add New Number</h3>
      <form @submit.prevent="handleSubmit">
        <div class="mb-4">
          <label class="block text-gray-700 mb-2" for="answ">Answer</label>
          <input id="answ" v-model.number="answ" type="number" class="w-full px-3 py-2 border rounded" required />
        </div>
        <div class="mb-4">
          <label class="block text-gray-700 mb-2" for="min">Min</label>
          <input id="min" v-model.number="min" type="number" class="w-full px-3 py-2 border rounded" />
        </div>
        <div class="mb-4">
          <label class="block text-gray-700 mb-2" for="max">Max</label>
          <input id="max" v-model.number="max" type="number" class="w-full px-3 py-2 border rounded" />
        </div>
        <div class="mb-4">
          <label class="block text-gray-700 mb-2" for="drop_count">Drop Count</label>
          <input id="drop_count" v-model.number="drop_count" type="number" class="w-full px-3 py-2 border rounded" />
        </div>
        <div class="flex justify-end">
          <button type="button" @click="$emit('close')" class="mr-2 px-4 py-2 border rounded">Cancel</button>
          <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Add</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '@/axios'

const emit = defineEmits(['close', 'added'])
let answ = ref(null)
let min = ref(null)
let max = ref(null)
let drop_count = ref(1)

async function handleSubmit() {
  try {
    await api.post('/predict/add', {
      answ: answ.value,
      min: min.value,
      max: max.value,
      drop_count: drop_count.value
    })
    emit('added')
  } catch (err) {
    console.error(err)
  }
}
</script>

<style scoped>
/* Optional custom styles */
</style>
