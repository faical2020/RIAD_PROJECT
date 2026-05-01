<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRiadStore } from '../stores/riad'

const riad = useRiadStore()
const searchQuery = ref('')
const filterStatus = ref('all')

onMounted(() => {
  riad.fetchRooms()
})

const filteredRooms = computed(() => {
  return riad.rooms.filter(room => {
    const matchesSearch = room.numero.toString().includes(searchQuery.value) ||
      (room.type && room.type.toLowerCase().includes(searchQuery.value.toLowerCase()))
    const matchesStatus = filterStatus.value === 'all' || room.statut === filterStatus.value
    return matchesSearch && matchesStatus
  })
})

function statutClass(s) {
  if (!s) return 'bg-riad-100 text-riad-600 border-riad-200'
  if (s === 'libre') return 'bg-green-100 text-green-700 border-green-200'
  if (s.includes('occup')) return 'bg-red-100 text-red-700 border-red-200'
  return 'bg-riad-100 text-riad-600 border-riad-200'
}
</script>

<template>
  <div class="space-y-4 sm:space-y-6">
    <div class="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-between items-stretch sm:items-center">
      <div class="relative flex-1 sm:max-w-md lg:max-w-lg">
        <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-riad-400">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </span>
        <input v-model="searchQuery" type="text" placeholder="Rechercher une chambre..."
          class="pl-10 pr-4 py-2.5 w-full bg-white border border-riad-200 rounded-xl focus:ring-2 focus:ring-gold-500/50 focus:border-gold-500/50 outline-none transition-all text-riad-900 shadow-sm hover:border-riad-300" />
      </div>

      <div class="flex gap-2 overflow-x-auto pb-2 sm:pb-0">
        <button v-for="status in ['all', 'libre', 'occupee']" :key="status"
          @click="filterStatus = status"
          :class="['px-3 sm:px-4 py-2 text-xs sm:text-sm font-bold rounded-lg transition-all duration-300 whitespace-nowrap',
            filterStatus === status ? 'bg-gradient-to-r from-gold-500 to-gold-400 text-riad-950 shadow-lg shadow-gold-500/20' : 'bg-riad-100 text-riad-600 hover:bg-riad-200']">
          {{ status === 'all' ? 'Tous' : status.charAt(0).toUpperCase() + status.slice(1) }}
        </button>
      </div>
    </div>

    <div v-if="riad.loading" class="flex justify-center items-center h-48 sm:h-64">
      <div class="animate-spin rounded-full h-10 w-10 sm:h-12 sm:w-12 border-b-2 border-gold-500"></div>
      <span class="ml-3 text-riad-600 text-sm sm:text-base">Chargement...</span>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-6">
      <div v-for="room in filteredRooms" :key="room.id"
        class="card p-4 sm:p-5 hover:shadow-2xl transition-all duration-300 group hover:-translate-y-1">
        <div class="flex justify-between items-start mb-3 sm:mb-4">
          <div class="flex items-center gap-2 sm:gap-3">
            <div class="w-8 h-8 sm:w-10 sm:h-10 bg-gradient-to-br from-gold-100 to-amber-50 text-gold-700 rounded-lg flex items-center justify-center font-bold shadow-sm">
              {{ room.numero }}
            </div>
            <div>
              <h3 class="text-base sm:text-lg font-bold text-riad-900">Chambre {{ room.numero }}</h3>
              <p class="text-riad-500 text-xs sm:text-sm capitalize">{{ room.type }}</p>
            </div>
          </div>
          <span :class="['badge border', statutClass(room.statut)]">{{ room.statut }}</span>
        </div>

        <div class="space-y-2 mb-3 sm:mb-4">
          <div class="flex justify-between text-xs sm:text-sm">
            <span class="text-riad-500">Prix / nuit</span>
            <span class="font-bold text-gold-600">{{ room.prix }} MAD</span>
          </div>
        </div>

        <p v-if="room.description" class="text-xs sm:text-sm text-riad-400 italic line-clamp-2 border-t border-riad-100 pt-3">
          {{ room.description }}
        </p>
      </div>
    </div>

    <div v-if="filteredRooms.length === 0 && !riad.loading" class="text-center py-8 sm:py-12 card p-6 sm:p-8">
      <span class="text-4xl sm:text-5xl block mb-3">🛏️</span>
      <p class="text-riad-500 text-sm sm:text-base">Aucune chambre trouvée</p>
    </div>
  </div>
</template>
