<script setup>
import { ref, onMounted } from 'vue'
import { useRiadStore } from '../stores/riad'
import { useAuthStore } from '../stores/auth'

const riad = useRiadStore()
const auth = useAuthStore()

onMounted(() => {
  riad.fetchReservations()
})
</script>

<template>
  <div class="space-y-4 sm:space-y-6">
    <div class="card overflow-hidden">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between p-4 sm:p-5 border-b border-riad-100 gap-2">
        <h3 class="font-display text-riad-900 font-semibold text-base sm:text-lg">Réservations</h3>
        <RouterLink v-if="auth.isReceptionniste" to="/nouvelle-reservation" 
          class="text-gold-600 text-xs sm:text-sm font-bold hover:text-gold-700 transition-all hover:translate-x-1 inline-block">
          + Nouvelle réservation
        </RouterLink>
      </div>
      <div class="p-4 sm:p-5 space-y-2 sm:space-y-3">
        <div v-if="riad.loading" class="text-center py-6 sm:py-8 text-riad-400">Chargement...</div>
        <div v-else-if="riad.reservations.length === 0" class="text-center py-6 sm:py-8 text-riad-400">
          <span class="text-3xl sm:text-4xl block mb-2">📋</span>
          Aucune réservation
        </div>
        <div v-else v-for="r in riad.reservations" :key="r.id"
          class="flex flex-col sm:flex-row sm:items-center justify-between p-3 sm:p-4 rounded-xl bg-riad-50 hover:bg-gradient-to-r hover:from-blue-50 hover:to-riad-50 transition-all duration-300 gap-2 sm:gap-0 group">
          <div>
            <p class="text-riad-800 text-sm font-bold">{{ r.date_debut }} → {{ r.date_fin }}</p>
            <p class="text-riad-400 text-xs mt-0.5">{{ r.montant }} MAD</p>
          </div>
          <span class="badge bg-blue-100 text-blue-700 border border-blue-200 group-hover:scale-105 transition-transform">{{ r.statut }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
