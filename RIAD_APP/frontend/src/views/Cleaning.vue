<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRiadStore } from '../stores/riad'

const riad = useRiadStore()
const filterCleaning = ref('all')

onMounted(async () => riad.fetchChambres())

async function updateStatus(roomId, newStatus) {
  try {
    await riad.updateCleaningStatus(roomId, newStatus)
  } catch (e) {
    alert('Erreur: ' + e)
  }
}

const filteredRooms = computed(() => {
  if (filterCleaning.value === 'all') return riad.rooms
  return riad.rooms.filter(r => r.cleaning_status === filterCleaning.value)
})

const stats = computed(() => {
  const total = riad.rooms.length
  const propre = riad.rooms.filter(r => r.cleaning_status === 'propre').length
  const sale = riad.rooms.filter(r => r.cleaning_status === 'sale').length
  const enCours = riad.rooms.filter(r => r.cleaning_status === 'en cours').length
  return { total, propre, sale, enCours, pct: total ? Math.round((propre / total) * 100) : 0 }
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-riad-900">Gestion du Ménage</h2>
        <p class="text-riad-400 text-sm">Mise à jour du statut de propreté</p>
      </div>
    </div>

    <!-- Progress -->
    <div class="card p-4 sm:p-5">
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm font-semibold text-riad-700">Propreté globale</span>
        <span class="text-sm font-bold text-riad-900">{{ stats.pct }}%</span>
      </div>
      <div class="h-3 rounded-full bg-riad-100 overflow-hidden flex">
        <div class="bg-green-500 h-full transition-all duration-500" :style="{ width: stats.pct + '%' }"></div>
        <div class="bg-amber-400 h-full transition-all duration-500" :style="{ width: stats.total ? Math.round((stats.enCours / stats.total) * 100) + '%' : '0%' }"></div>
      </div>
      <div class="flex items-center gap-4 mt-3 text-xs text-riad-500">
        <span><span class="inline-block w-2 h-2 rounded-full bg-green-500 mr-1"></span>{{ stats.propre }} propres</span>
        <span><span class="inline-block w-2 h-2 rounded-full bg-amber-400 mr-1"></span>{{ stats.enCours }} en cours</span>
        <span><span class="inline-block w-2 h-2 rounded-full bg-red-500 mr-1"></span>{{ stats.sale }} sales</span>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex gap-2 flex-wrap">
      <button v-for="f in [{ key: 'all', label: 'Toutes' }, { key: 'propre', label: 'Propres' }, { key: 'en cours', label: 'En cours' }, { key: 'sale', label: 'Sales' }]" :key="f.key"
        @click="filterCleaning = f.key"
        :class="['px-4 py-2 rounded-xl text-sm font-semibold transition-all border',
          filterCleaning === f.key
            ? 'bg-riad-900 text-white border-riad-900 shadow-md'
            : 'bg-white text-riad-600 border-riad-200 hover:border-riad-300 hover:text-riad-900']">
        {{ f.label }}
      </button>
    </div>

    <!-- Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <div v-for="room in filteredRooms" :key="room.id"
        class="card p-5 animate-in">
        <div class="flex items-start justify-between mb-4">
          <div>
            <h3 class="font-semibold text-riad-900">Chambre {{ room.numero }}</h3>
            <p class="text-xs text-riad-400 capitalize">{{ room.type }}</p>
          </div>
          <span :class="[
            'px-2.5 py-1 rounded-lg text-xs font-semibold border',
            room.cleaning_status === 'propre' ? 'bg-green-50 text-green-700 border-green-200' :
            room.cleaning_status === 'en cours' ? 'bg-amber-50 text-amber-700 border-amber-200' :
            'bg-red-50 text-red-700 border-red-200'
          ]">
            {{ room.cleaning_status }}
          </span>
        </div>

        <div class="grid grid-cols-3 gap-2">
          <button @click="updateStatus(room.id, 'propre')"
            class="py-2 rounded-lg text-xs font-semibold border border-green-200 text-green-600 hover:bg-green-500 hover:text-white transition-all">
            Propre
          </button>
          <button @click="updateStatus(room.id, 'en cours')"
            class="py-2 rounded-lg text-xs font-semibold border border-amber-200 text-amber-600 hover:bg-amber-500 hover:text-white transition-all">
            En cours
          </button>
          <button @click="updateStatus(room.id, 'sale')"
            class="py-2 rounded-lg text-xs font-semibold border border-red-200 text-red-600 hover:bg-red-500 hover:text-white transition-all">
            Sale
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
