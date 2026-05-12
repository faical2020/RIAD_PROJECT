<script setup>
import { computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRiadStore } from '../stores/riad'

const authStore = useAuthStore()
const riad = useRiadStore()

onMounted(async () => {
  await riad.fetchRooms()
  if (authStore.isStaff) {
    await riad.fetchReservations()
  } else {
    await riad.fetchMyReservations()
  }
})

const stats = computed(() => [
  { icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"/></svg>', label: 'Chambres', value: riad.rooms.length, color: 'from-gold-500 to-amber-500', bg: 'bg-gold-50 text-gold-600' },
  { icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>', label: 'Libres', value: riad.availableRooms.length, color: 'from-green-500 to-emerald-500', bg: 'bg-green-50 text-green-600' },
  { icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/></svg>', label: 'Occupées', value: riad.occupiedRooms.length, color: 'from-red-500 to-rose-500', bg: 'bg-red-50 text-red-600' },
  ...(authStore.isStaff
    ? [{ icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/></svg>', label: 'Réservations', value: riad.reservations.length, color: 'from-blue-500 to-indigo-500', bg: 'bg-blue-50 text-blue-600' }]
    : [{ icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/></svg>', label: 'Mes réservations', value: riad.reservations.length, color: 'from-blue-500 to-indigo-500', bg: 'bg-blue-50 text-blue-600' }]
  )
])

const cleaningStats = computed(() => {
  const total = riad.rooms.length
  if (!total) return { propre: 0, sale: 0, enCours: 0, proprePct: 0 }
  const propre = riad.rooms.filter(r => r.cleaning_status === 'propre').length
  const sale = riad.rooms.filter(r => r.cleaning_status === 'sale').length
  const enCours = riad.rooms.filter(r => r.cleaning_status === 'en cours').length
  return { propre, sale, enCours, proprePct: Math.round((propre / total) * 100) }
})

function statutClass(s) {
  if (!s) return 'badge-gray'
  if (s === 'libre') return 'badge-green'
  if (s.includes('occup')) return 'badge-red'
  return 'badge-gray'
}
</script>

<template>
  <div class="space-y-6 lg:space-y-8">
    <!-- Welcome -->
    <div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-riad-950 via-riad-900 to-riad-950 p-6 sm:p-8">
      <div class="absolute top-0 right-0 w-64 h-64 bg-gold-500/10 rounded-full blur-3xl"></div>
      <div class="relative">
        <p class="text-gold-400 font-arabic text-sm sm:text-base mb-1">مرحباً بكم في رياضنا</p>
        <h2 class="text-white text-xl sm:text-2xl lg:text-3xl font-bold mb-1">
          Bonjour, {{ authStore.user?.prenom || 'Visiteur' }}
        </h2>
        <p class="text-riad-400 text-sm">Voici un aperçu de l'activité du Riad</p>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
      <div v-for="stat in stats" :key="stat.label"
        class="card p-4 sm:p-5 hover:-translate-y-0.5">
        <div class="flex items-center gap-3 mb-3">
          <div :class="['w-10 h-10 rounded-xl flex items-center justify-center', stat.bg]" v-html="stat.icon"></div>
        </div>
        <p class="text-2xl sm:text-3xl font-bold text-riad-900">{{ stat.value }}</p>
        <p class="text-riad-500 text-sm mt-0.5">{{ stat.label }}</p>
        <div class="mt-3 h-1.5 rounded-full bg-riad-100 overflow-hidden">
          <div :class="['h-full rounded-full bg-gradient-to-r transition-all duration-500', stat.color]" :style="{ width: riad.rooms.length ? Math.round((stat.value / riad.rooms.length) * 100) + '%' : '0%' }"></div>
        </div>
      </div>
    </div>

    <!-- Cleaning quick overview -->
    <div class="card p-4 sm:p-5">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-semibold text-riad-900">État du ménage</h3>
        <span class="text-xs text-riad-400">{{ cleaningStats.proprePct }}% propre</span>
      </div>
      <div class="h-2.5 rounded-full bg-riad-100 overflow-hidden flex">
        <div class="bg-green-500 h-full transition-all" :style="{ width: cleaningStats.proprePct + '%' }"></div>
        <div class="bg-amber-400 h-full transition-all" :style="{ width: riad.rooms.length ? Math.round((cleaningStats.enCours / riad.rooms.length) * 100) + '%' : '0%' }"></div>
      </div>
      <div class="flex items-center gap-4 mt-3 text-xs text-riad-500">
        <span class="flex items-center gap-1.5"><span class="w-2 h-2 rounded-full bg-green-500"></span> {{ cleaningStats.propre }} propres</span>
        <span class="flex items-center gap-1.5"><span class="w-2 h-2 rounded-full bg-amber-400"></span> {{ cleaningStats.enCours }} en cours</span>
        <span class="flex items-center gap-1.5"><span class="w-2 h-2 rounded-full bg-red-500"></span> {{ cleaningStats.sale }} sales</span>
      </div>
    </div>

    <!-- Two columns -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
      <!-- Rooms -->
      <div class="card overflow-hidden">
        <div class="flex items-center justify-between p-4 sm:p-5 border-b border-riad-100">
          <h3 class="font-semibold text-riad-900">État des chambres</h3>
          <RouterLink to="/chambres" class="text-sm text-gold-600 hover:text-gold-700 font-semibold">Voir tout →</RouterLink>
        </div>
        <div class="p-4 sm:p-5 space-y-2">
          <div v-if="riad.loading" class="text-center py-8 text-riad-400 text-sm">Chargement...</div>
          <div v-else-if="riad.rooms.length === 0" class="text-center py-8 text-riad-400">
            <p class="text-4xl mb-2">🛏️</p>
            Aucune chambre
          </div>
          <div v-else v-for="c in riad.rooms.slice(0,6)" :key="c.id"
            class="flex items-center justify-between p-3 rounded-xl bg-riad-50 hover:bg-gold-50/50 transition-colors">
            <div class="flex items-center gap-3">
              <span class="w-8 h-8 rounded-lg bg-white text-gold-700 flex items-center justify-center text-sm font-bold shadow-sm border border-riad-200">
                {{ c.numero }}
              </span>
              <div>
                <p class="text-sm font-semibold text-riad-800 capitalize">{{ c.type }}</p>
                <p class="text-xs text-riad-400">{{ c.prix }} MAD/nuit</p>
              </div>
            </div>
            <span :class="statutClass(c.statut)">{{ c.statut }}</span>
          </div>
        </div>
      </div>

      <!-- Reservations -->
      <div class="card overflow-hidden">
        <div class="flex items-center justify-between p-4 sm:p-5 border-b border-riad-100">
          <h3 class="font-semibold text-riad-900">{{ authStore.isStaff ? 'Réservations récentes' : 'Mes réservations' }}</h3>
          <RouterLink v-if="authStore.isStaff" to="/reservations" class="text-sm text-gold-600 hover:text-gold-700 font-semibold">Voir tout →</RouterLink>
        </div>
        <div class="p-4 sm:p-5 space-y-2">
          <div v-if="riad.loading" class="text-center py-8 text-riad-400 text-sm">Chargement...</div>
          <div v-else-if="riad.reservations.length === 0" class="text-center py-8 text-riad-400">
            <p class="text-4xl mb-2">📋</p>
            Aucune réservation
          </div>
          <div v-else v-for="r in riad.reservations.slice(0,5)" :key="r.id"
            class="flex items-center justify-between p-3 rounded-xl bg-riad-50 hover:bg-blue-50/50 transition-colors">
            <div>
              <p class="text-sm font-semibold text-riad-800">{{ r.date_debut }} → {{ r.date_fin }}</p>
              <p class="text-xs text-riad-400 mt-0.5">{{ r.montant }} MAD</p>
            </div>
            <span class="badge-blue">{{ r.statut }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
