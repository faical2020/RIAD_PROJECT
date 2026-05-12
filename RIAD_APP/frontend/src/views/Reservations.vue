<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRiadStore } from '../stores/riad'
import { useAuthStore } from '../stores/auth'

const riad = useRiadStore()
const auth = useAuthStore()
const searchQuery = ref('')

onMounted(() => riad.fetchReservations())

const groupedReservations = computed(() => {
  const groups = { upcoming: [], active: [], past: [] }
  const today = new Date().toISOString().split('T')[0]
  for (const r of riad.reservations) {
    if (r.date_debut > today) groups.upcoming.push(r)
    else if (r.date_fin < today) groups.past.push(r)
    else groups.active.push(r)
  }
  return groups
})

function statutBadge(s) {
  const map = { 'confirmée': 'badge-blue', 'en attente': 'badge-amber', 'checkin': 'badge-green', 'checkout': 'badge-gray', 'annulée': 'badge-red' }
  return map[s] || 'badge-gray'
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-riad-900">Réservations</h2>
        <p class="text-riad-400 text-sm">{{ riad.reservations.length }} réservation(s)</p>
      </div>
      <RouterLink v-if="auth.isReceptionniste" to="/nouvelle-reservation" class="btn-primary btn-sm">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>
        Nouvelle
      </RouterLink>
    </div>

    <!-- Loading -->
    <div v-if="riad.loading" class="flex items-center justify-center py-16">
      <svg class="animate-spin h-8 w-8 text-gold-500" viewBox="0 0 24 24" fill="none"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
    </div>

    <!-- Content -->
    <template v-else>
      <!-- Active -->
      <div v-if="groupedReservations.active.length" class="card overflow-hidden">
        <div class="p-4 sm:p-5 border-b border-riad-100">
          <h3 class="font-semibold text-riad-900 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-green-500"></span>
            En cours ({{ groupedReservations.active.length }})
          </h3>
        </div>
        <div class="p-4 sm:p-5 space-y-2">
          <div v-for="r in groupedReservations.active" :key="r.id"
            class="flex flex-col sm:flex-row sm:items-center justify-between p-3 rounded-xl bg-green-50/50 gap-2">
            <div>
              <p class="text-sm font-semibold text-riad-800">{{ r.date_debut }} → {{ r.date_fin }}</p>
              <p class="text-xs text-riad-400 mt-0.5">{{ r.montant }} MAD</p>
            </div>
            <span :class="statutBadge(r.statut)">{{ r.statut }}</span>
          </div>
        </div>
      </div>

      <!-- Upcoming -->
      <div v-if="groupedReservations.upcoming.length" class="card overflow-hidden">
        <div class="p-4 sm:p-5 border-b border-riad-100">
          <h3 class="font-semibold text-riad-900">À venir ({{ groupedReservations.upcoming.length }})</h3>
        </div>
        <div class="p-4 sm:p-5 space-y-2">
          <div v-for="r in groupedReservations.upcoming" :key="r.id"
            class="flex flex-col sm:flex-row sm:items-center justify-between p-3 rounded-xl bg-riad-50 hover:bg-blue-50/50 transition-colors gap-2">
            <div>
              <p class="text-sm font-semibold text-riad-800">{{ r.date_debut }} → {{ r.date_fin }}</p>
              <p class="text-xs text-riad-400 mt-0.5">{{ r.montant }} MAD</p>
            </div>
            <span :class="statutBadge(r.statut)">{{ r.statut }}</span>
          </div>
        </div>
      </div>

      <!-- Past -->
      <div v-if="groupedReservations.past.length" class="card overflow-hidden">
        <div class="p-4 sm:p-5 border-b border-riad-100">
          <h3 class="font-semibold text-riad-900">Passées ({{ groupedReservations.past.length }})</h3>
        </div>
        <div class="p-4 sm:p-5 space-y-2">
          <div v-for="r in groupedReservations.past" :key="r.id"
            class="flex flex-col sm:flex-row sm:items-center justify-between p-3 rounded-xl bg-riad-50 gap-2 opacity-60">
            <div>
              <p class="text-sm font-semibold text-riad-800">{{ r.date_debut }} → {{ r.date_fin }}</p>
              <p class="text-xs text-riad-400 mt-0.5">{{ r.montant }} MAD</p>
            </div>
            <span :class="statutBadge(r.statut)">{{ r.statut }}</span>
          </div>
        </div>
      </div>

      <div v-if="riad.reservations.length === 0" class="card p-8 text-center">
        <p class="text-4xl mb-3">📋</p>
        <p class="text-riad-500 text-sm">Aucune réservation</p>
      </div>
    </template>
  </div>
</template>
