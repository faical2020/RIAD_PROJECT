<script setup>
import { computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRiadStore } from '../stores/riad'

const authStore = useAuthStore()
const riad = useRiadStore()

onMounted(async () => {
  await riad.fetchRooms()
  await riad.fetchReservations()
})

const stats = computed(() => {
  const base = [
    { icon: '🛏️', label: 'Chambres total', value: riad.rooms.length, bg: 'from-gold-100 to-amber-50', text: 'text-gold-700', iconBg: 'bg-gold-500/10' },
    { icon: '✅', label: 'Chambres libres', value: riad.availableRooms.length, bg: 'from-green-100 to-emerald-50', text: 'text-green-700', iconBg: 'bg-green-500/10' },
    { icon: '🔴', label: 'Occupées', value: riad.occupiedRooms.length, bg: 'from-red-100 to-rose-50', text: 'text-red-700', iconBg: 'bg-red-500/10' },
  ]
  if (['manager', 'receptionniste'].includes(authStore.role)) {
    base.push({ icon: '📋', label: 'Réservations', value: riad.reservations.length, bg: 'from-blue-100 to-sky-50', text: 'text-blue-700', iconBg: 'bg-blue-500/10' })
  }
  return base
})

function statutChambreClass(s) {
  if (!s) return 'bg-riad-100 text-riad-600'
  if (s === 'libre') return 'bg-green-100 text-green-700 border-green-200'
  if (s.includes('occup')) return 'bg-red-100 text-red-700 border-red-200'
  return 'bg-riad-100 text-riad-600'
}
</script>

<template>
  <div class="space-y-4 sm:space-y-6 lg:space-y-8">
    <!-- Welcome banner -->
    <div class="relative bg-gradient-to-br from-riad-950 via-riad-900 to-riad-950 rounded-2xl sm:rounded-3xl p-6 sm:p-8 overflow-hidden">
      <div class="absolute inset-0 ornament opacity-30"></div>
      <div class="absolute top-0 right-0 w-48 sm:w-64 h-48 sm:h-64 bg-gold-500/10 rounded-full blur-3xl animate-float"></div>
      <div class="absolute bottom-0 left-0 w-32 sm:w-48 h-32 sm:h-48 bg-riad-700/30 rounded-full blur-3xl"></div>
      <div class="relative">
        <p class="text-gold-400 font-arabic text-base sm:text-lg mb-1">مرحباً بكم في رياضنا</p>
        <h2 class="font-display text-white text-xl sm:text-2xl lg:text-3xl mb-2 sm:mb-3">
          Bonjour, {{ authStore.user?.prenom || 'Visiteur' }} 👋
        </h2>
        <p class="text-riad-400 text-sm sm:text-base">Voici un aperçu de l'activité du Riad</p>
      </div>
    </div>

    <!-- Stats cards -->
    <div class="grid grid-cols-2 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 lg:gap-6">
      <div v-for="stat in stats" :key="stat.label" class="group bg-white rounded-xl sm:rounded-2xl p-4 sm:p-5 lg:p-6 border border-riad-100 hover:shadow-2xl hover:shadow-gold-500/5 transition-all duration-300 hover:-translate-y-1">
        <div class="flex items-center gap-2 sm:gap-3 mb-2 sm:mb-3">
          <div :class="['w-8 h-8 sm:w-10 sm:h-10 rounded-lg sm:rounded-xl flex items-center justify-center text-lg sm:text-xl', stat.iconBg]">
            {{ stat.icon }}
          </div>
        </div>
        <p class="font-display text-xl sm:text-2xl lg:text-3xl font-bold text-riad-900">{{ stat.value }}</p>
        <p class="text-riad-500 text-xs sm:text-sm mt-0.5 sm:mt-1">{{ stat.label }}</p>
      </div>
    </div>

    <!-- Two columns -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
      <!-- Recent rooms -->
      <div class="card overflow-hidden">
        <div class="flex items-center justify-between p-4 sm:p-5 border-b border-riad-100">
          <h3 class="font-display text-riad-900 font-semibold text-base sm:text-lg">État des chambres</h3>
          <RouterLink to="/chambres" class="text-gold-600 text-xs sm:text-sm font-bold hover:text-gold-700 transition-all hover:translate-x-1 inline-block">Voir tout →</RouterLink>
        </div>
        <div class="p-4 sm:p-5 space-y-2 sm:space-y-3">
          <div v-if="riad.loading" class="text-center py-6 sm:py-8 text-riad-400">Chargement...</div>
          <div v-else-if="riad.rooms.length === 0" class="text-center py-6 sm:py-8 text-riad-400">
            <span class="text-3xl block mb-2">🛏️</span>
            Aucune chambre enregistrée
          </div>
          <div v-else v-for="c in riad.rooms.slice(0,6)" :key="c.id"
            class="flex items-center justify-between p-2.5 sm:p-3 rounded-lg sm:rounded-xl bg-riad-50 hover:bg-gradient-to-r hover:from-gold-50 hover:to-riad-50 transition-all duration-300 group">
            <div class="flex items-center gap-2 sm:gap-3">
              <span class="w-7 h-7 sm:w-8 sm:h-8 rounded-lg bg-gold-100 flex items-center justify-center text-xs sm:text-sm font-bold text-gold-700 shadow-sm">
                {{ c.numero }}
              </span>
              <div>
                <p class="text-riad-800 text-xs sm:text-sm font-bold capitalize">{{ c.type }}</p>
                <p class="text-riad-400 text-xs">{{ c.prix }} MAD/nuit</p>
              </div>
            </div>
            <span :class="['badge border', statutChambreClass(c.statut)]">{{ c.statut }}</span>
          </div>
        </div>
      </div>

      <!-- Recent reservations -->
      <div class="card overflow-hidden">
        <div class="flex items-center justify-between p-4 sm:p-5 border-b border-riad-100">
          <h3 class="font-display text-riad-900 font-semibold text-base sm:text-lg">Réservations récentes</h3>
          <RouterLink to="/reservations" class="text-gold-600 text-xs sm:text-sm font-bold hover:text-gold-700 transition-all hover:translate-x-1 inline-block">Voir tout →</RouterLink>
        </div>
        <div class="p-4 sm:p-5 space-y-2 sm:space-y-3">
          <div v-if="riad.loading" class="text-center py-6 sm:py-8 text-riad-400">Chargement...</div>
          <div v-else-if="riad.reservations.length === 0" class="text-center py-6 sm:py-8 text-riad-400">
            <span class="text-3xl block mb-2">📋</span>
            Aucune réservation
          </div>
          <div v-else v-for="r in riad.reservations.slice(0,5)" :key="r.id"
            class="flex items-center justify-between p-2.5 sm:p-3 rounded-lg sm:rounded-xl bg-riad-50 hover:bg-gradient-to-r hover:from-blue-50 hover:to-riad-50 transition-all duration-300">
            <div>
              <p class="text-riad-800 text-xs sm:text-sm font-bold">{{ r.date_debut }} → {{ r.date_fin }}</p>
              <p class="text-riad-400 text-xs mt-0.5">{{ r.montant }} MAD</p>
            </div>
            <span class="badge bg-blue-100 text-blue-700 border border-blue-200">{{ r.statut }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
