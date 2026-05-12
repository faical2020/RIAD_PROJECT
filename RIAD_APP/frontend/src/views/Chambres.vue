<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRiadStore } from '../stores/riad'

const riad = useRiadStore()
const searchQuery = ref('')
const filterStatus = ref('all')

onMounted(() => riad.fetchRooms())

const filteredRooms = computed(() =>
  riad.rooms.filter(room => {
    const q = searchQuery.value.toLowerCase()
    const matchesSearch = !q || room.numero.toString().includes(q) ||
      (room.type && room.type.toLowerCase().includes(q))
    const matchesStatus = filterStatus.value === 'all' || room.statut === filterStatus.value
    return matchesSearch && matchesStatus
  })
)

const totalRooms = computed(() => riad.rooms.length)
const libreCount = computed(() => riad.rooms.filter(r => r.statut === 'libre').length)

function statutClass(s) {
  if (!s) return 'badge-gray'
  if (s === 'libre') return 'badge-green'
  if (s.includes('occup')) return 'badge-red'
  return 'badge-gray'
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-riad-900">Chambres</h2>
        <p class="text-riad-400 text-sm">{{ libreCount }} disponibles sur {{ totalRooms }}</p>
      </div>
      <div class="flex items-center gap-2 w-full sm:w-auto">
        <div class="relative flex-1 sm:flex-none sm:w-64">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-riad-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/></svg>
          <input v-model="searchQuery" type="text" placeholder="Rechercher..."
            class="input pl-9 py-2 text-sm" />
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex gap-2 flex-wrap">
      <button v-for="s in ['all', 'libre', 'occupee', 'occupe']" :key="s"
        @click="filterStatus = s"
        :class="['px-4 py-2 rounded-xl text-sm font-semibold transition-all border',
          filterStatus === s
            ? 'bg-riad-900 text-white border-riad-900 shadow-md'
            : 'bg-white text-riad-600 border-riad-200 hover:border-riad-300 hover:text-riad-900']">
        {{ s === 'all' ? 'Toutes' : s.charAt(0).toUpperCase() + s.slice(1) }}
      </button>
    </div>

    <!-- Loading -->
    <div v-if="riad.loading" class="flex items-center justify-center py-16">
      <svg class="animate-spin h-8 w-8 text-gold-500" viewBox="0 0 24 24" fill="none"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
      <span class="ml-3 text-riad-500 text-sm">Chargement...</span>
    </div>

    <!-- Grid -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <div v-for="room in filteredRooms" :key="room.id"
        class="card p-5 hover:-translate-y-0.5 animate-in">
        <div class="flex items-start justify-between mb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-gold-50 to-amber-50 border border-gold-200 flex items-center justify-center text-sm font-bold text-gold-700">
              {{ room.numero }}
            </div>
            <div>
              <h3 class="font-semibold text-riad-900">Chambre {{ room.numero }}</h3>
              <p class="text-xs text-riad-500 capitalize">{{ room.type || 'Standard' }}</p>
            </div>
          </div>
          <span :class="statutClass(room.statut)">{{ room.statut }}</span>
        </div>

        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-riad-400">Prix</span>
            <span class="font-semibold text-riad-900">{{ room.prix }} MAD</span>
          </div>
        </div>

        <p v-if="room.description" class="text-xs text-riad-400 mt-3 pt-3 border-t border-riad-100 line-clamp-2">
          {{ room.description }}
        </p>
      </div>
    </div>

    <!-- Empty -->
    <div v-if="!riad.loading && filteredRooms.length === 0"
      class="card p-8 text-center">
      <p class="text-4xl mb-3">🛏️</p>
      <p class="text-riad-500 text-sm">Aucune chambre trouvée</p>
    </div>
  </div>
</template>
