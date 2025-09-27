<template>
    <div class="p-6">
        <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-bold">Predicts List</h2>
            <div class="flex space-x-2">
                <button @click="showSmartModal = true" 
                                class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors flex items-center">
                    <span class="material-icons mr-1">smart_toy</span> Smart Add
                </button>
                <button @click="showModal = true" 
                                class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 transition-colors flex items-center">
                    <span class="material-icons mr-1">add</span> Add Number
                </button>
            </div>
        </div>

        <!-- Loading and error states -->
        <div v-if="loading" class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            <p>{{ error }}</p>
            <button @click="fetchPredicts" class="underline">Try again</button>
        </div>

        <div v-else-if="predicts.length === 0" class="text-center py-8 bg-gray-50 rounded">
            <p class="text-gray-500">No predicts found. Add your first one!</p>
        </div>

        <!-- Search and filter -->
        <div v-else class="mb-4">
            <input v-model="search" placeholder="Search..." 
                        class="px-4 py-2 border rounded w-full md:w-64 mb-4"/>
            
            <div class="overflow-x-auto bg-white rounded-lg shadow">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th v-for="column in columns" :key="column.key" 
                                    @click="sortBy(column.key)"
                                    class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100">
                                {{ column.label }}
                                <span v-if="sortKey === column.key">{{ sortOrder === 'asc' ? '▲' : '▼' }}</span>
                            </th>
                            <th class="px-4 py-3 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        <tr v-for="item in filteredAndSortedPredicts" :key="item.uuid" class="hover:bg-gray-50">
                            <td class="px-4 py-3 text-sm font-mono">{{ item.uuid }}</td>
                            <td class="px-4 py-3 text-sm font-medium">{{ item.answ }}</td>
                            <td class="px-4 py-3 text-sm">{{ item.min }}</td>
                            <td class="px-4 py-3 text-sm">{{ item.max }}</td>
                            <td class="px-4 py-3 text-sm">{{ item.dropped }}</td>
                            <td class="px-4 py-3 text-sm">{{ item.max_drops }}</td>
                            <td class="px-4 py-3 text-sm">
                                <span :class="item.finished ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'" 
                                            class="px-2 py-1 rounded-full text-xs font-medium">
                                    {{ item.finished ? 'Completed' : 'In Progress' }}
                                </span>
                            </td>
                            <td class="px-4 py-3 text-sm">{{ formatDate(item.added) }}</td>
                            <td class="px-4 py-3 text-sm text-right">
                                <button @click="confirmDelete(item)" class="text-red-500 hover:text-red-700 ml-2">
                                    <span class="material-icons text-sm">delete</span>
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Modals -->
        <AddPredictModal v-if="showModal" @close="showModal = false" @added="handleAdded" />
        
        <SmartGeneration v-if="showSmartModal" @close="showSmartModal = false" />
        
        <div v-if="showDeleteConfirm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div class="bg-white p-6 rounded-lg shadow-lg max-w-md w-full">
                <h3 class="text-lg font-bold mb-4">Confirm Deletion</h3>
                <p>Are you sure you want to delete this predict?</p>
                <div class="flex justify-end mt-6 space-x-2">
                    <button @click="showDeleteConfirm = false" class="px-4 py-2 border rounded">Cancel</button>
                    <button @click="deletePredict" class="px-4 py-2 bg-red-500 text-white rounded">Delete</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/axios'
import AddPredictModal from '@/view/AddPredictModal.vue'
import SmartGeneration from '@/view/SmartGeneration.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const predicts = ref([])
const showModal = ref(false)
const loading = ref(true)
const error = ref(null)
const search = ref('')
const sortKey = ref('added')
const sortOrder = ref('desc')
const showDeleteConfirm = ref(false)
const itemToDelete = ref(null)
const showSmartModal = ref(false)

const columns = [
    { key: 'uuid', label: 'UUID' },
    { key: 'answ', label: 'Answer' },
    { key: 'min', label: 'Min' },
    { key: 'max', label: 'Max' },
    { key: 'dropped', label: 'Dropped' },
    { key: 'max_drops', label: 'Max Drops' },
    { key: 'finished', label: 'Status' },
    { key: 'added', label: 'Added' }
]

const filteredAndSortedPredicts = computed(() => {
    let result = [...predicts.value]
    
    // Filter by search term
    if (search.value) {
        const searchLower = search.value.toLowerCase()
        result = result.filter(item => 
            item.uuid.toLowerCase().includes(searchLower) || 
            String(item.answ).includes(searchLower)
        )
    }
    
    // Sort
    result.sort((a, b) => {
        const aValue = sortKey.value === 'added' ? a[sortKey.value] * 1000 : a[sortKey.value]
        const bValue = sortKey.value === 'added' ? b[sortKey.value] * 1000 : b[sortKey.value]
        
        if (aValue < bValue) return sortOrder.value === 'asc' ? -1 : 1
        if (aValue > bValue) return sortOrder.value === 'asc' ? 1 : -1
        return 0
    })
    
    return result
})

function formatDate(ts) {
    return new Date(ts * 1000).toLocaleString(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

function sortBy(key) {
    if (sortKey.value === key) {
        sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    } else {
        sortKey.value = key
        sortOrder.value = 'asc'
    }
}

function confirmDelete(item) {
    itemToDelete.value = item
    showDeleteConfirm.value = true
}

async function deletePredict() {
    try {
        await api.delete(`/predict/${itemToDelete.value.uuid}`)
        await fetchPredicts()
    } catch (err) {
        error.value = "Failed to delete predict: " + err.message
    } finally {
        showDeleteConfirm.value = false
        itemToDelete.value = null
    }
}

async function fetchPredicts() {
    loading.value = true
    error.value = null
    
    try {
        const res = await api.get('/predict/list')
        predicts.value = res.data
    } catch (err) {
        if (err.response && err.response.status === 401) {
            router.push('/login')
            return
        }
        console.error(err)
        error.value = "Failed to load predicts: " + err.message
    } finally {
        loading.value = false
    }
}

async function handleAdded() {
    showModal.value = false
    await fetchPredicts()
}

onMounted(() => {
    fetchPredicts()
})
</script>