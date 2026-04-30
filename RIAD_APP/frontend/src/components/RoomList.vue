<script setup>
import { ref, onMounted, computed } from 'vue'
import { riadService } from '../services/riadService'

const rooms = ref([])
const loading = ref(true)
const searchQuery = ref('')
const filterStatus = ref('all')

const fetchRooms = async () => {
  loading.value = true
  try {
    rooms.value = await riadService.getRooms()
  } catch (e) {
    console.error("Error fetching rooms:", e)
  } finally {
    loading.value = false
  }
}

const filteredRooms = computed(() => {
  return rooms.value.filter(room => {
    const matchesSearch = room.numero.toString().includes(searchQuery.value) || 
                          room.description.toLowerCase().includes(searchQuery.value.toLowerCase());
    const matchesStatus = filterStatus.value === 'all' || room.statut === filterStatus.value;
    return matchesSearch && matchesStatus;
  })
})

onMounted(fetchRooms)
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col md:flex-row gap-4 justify-between items-center bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
      <div class="relative w-full md:w-96">
        <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </span>
        <input v-model="searchQuery" type="text" placeholder="Rechercher une chambre..." 
               class="pl-10 pr-4 py-2 w-full border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
      </div>
      
      <div class="flex gap-2">
        <button v-for="status in ['all', 'libre', 'occupee']" :key="status"
                @click="filterStatus = status"
                :class="['px-4 py-2 text-sm font-medium rounded-lg transition-all', 
                         filterStatus === status ? 'bg-emerald-600 text-white shadow-sm' : 'bg-gray-100 text-gray-600 hover:bg-gray-200']">
          {{ status === 'all' ? 'Tous' : status.charAt(0).toUpperCase() + status.slice(1) }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-500"></div>
      <span class="ml-3 text-gray-600">Chargement...</span>
    </div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="room in filteredRooms" :key="room.id" 
           class="bg-white rounded-xl shadow-sm border border-gray-100 p-5 hover:shadow-md transition-all group">
        <div class="flex justify-between items-start mb-4">
          <div class="flex items-center gap-2">
            <div class="w-10 h-10 bg-emerald-100 text-emerald-700 rounded-lg flex items-center justify-center font-bold">
              {{ room.numero }}
            </div>
            <h3 class="text-lg font-semibold text-gray-900">Chambre {{ room.numero }}</h3>
          </div>
          <span :class="[
            'px-2 py-1 text-xs font-medium rounded-full',
            room.statut === 'libre' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
          ]">
            {{ room.statut }}
          </span>
        </div>
        
        <div class="space-y-2 mb-4">
          <div class="flex justify-between text-sm">
            <span class="text-gray-500">Type</span> 
            <span class="font-medium text-gray-900">{{ room.type }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-500">Prix / nuit</span> 
            <span class="font-bold text-emerald-600">{{ room.prix }} DH</span>
          </div>
        </div>
        
        <p class="text-sm text-gray-500 italic line-clamp-2 border-t pt-3">
          {{ room.description }}
        </p>
      </div>
    </div>
    
    <div v-if="filteredRooms.length === 0 && !loading" class="text-center py-12">
      <p class="text-gray-500">Aucune chambre trouvée correspondant à vos critères.</p>
    </div>
  </div>
</template>
