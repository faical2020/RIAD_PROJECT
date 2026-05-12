<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const reservationId = route.params.reservationId

const API = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1'
const headers = () => ({ 'Content-Type': 'application/json', Authorization: `Bearer ${localStorage.getItem('token')}` })

const facture = ref(null)
const loading = ref(true)
const showPaiement = ref(false)
const paiement = ref({ montant: 0, mode_paiement: 'especes', reference: '' })

async function load() {
  loading.value = true
  try {
    const r = await fetch(`${API}/reservations/${reservationId}/facture`, { headers: headers() })
    if (r.ok) {
      facture.value = await r.json()
      paiement.value.montant = facture.value.restant_du
    }
  } catch {}
  loading.value = false
}

onMounted(load)

async function submitPaiement() {
  if (paiement.value.montant <= 0) return
  try {
    await fetch(`${API}/reservations/${reservationId}/paiement`, {
      method: 'POST', headers: headers(),
      body: JSON.stringify(paiement.value)
    })
    showPaiement.value = false
    load()
  } catch {}
}

const statutBadge = (s) => {
  const map = { 'confirmée': 'badge-blue', 'checkin': 'badge-green', 'checkout': 'badge-gray', 'annulée': 'badge-red' }
  return map[s] || 'badge-gray'
}
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-bold text-riad-900">🧾 Facture</h2>
      <button @click="router.back()" class="btn-ghost btn-sm">Retour</button>
    </div>

    <div v-if="loading" class="text-center py-12 text-riad-400">Chargement...</div>

    <template v-else-if="facture">
      <!-- Header -->
      <div class="card p-5">
        <div class="flex items-center justify-between mb-2">
          <div>
            <h3 class="font-semibold text-riad-900">Riad — Facture</h3>
            <p class="text-xs text-riad-400">Chambre {{ facture.chambre.numero }} — {{ facture.chambre.type }}</p>
          </div>
          <span :class="statutBadge(facture.reservation.statut)">{{ facture.reservation.statut }}</span>
        </div>
        <p class="text-sm text-riad-500">{{ facture.reservation.date_debut }} → {{ facture.reservation.date_fin }}</p>
      </div>

      <!-- Lines -->
      <div class="card overflow-hidden">
        <div class="p-4 sm:p-5 border-b border-riad-100">
          <h3 class="font-semibold">Détail</h3>
        </div>
        <div class="p-4 sm:p-5 space-y-3">
          <!-- Séjour -->
          <div class="flex justify-between items-center py-2">
            <div>
              <p class="font-semibold text-riad-900">Séjour</p>
              <p class="text-xs text-riad-400">{{ facture.chambre.type }} × {{ facture.reservation.date_debut }} → {{ facture.reservation.date_fin }}</p>
            </div>
            <span class="font-bold">{{ facture.reservation.montant }} MAD</span>
          </div>

          <!-- Consommations -->
          <div v-for="c in facture.consommations" :key="c.id" class="flex justify-between items-center py-1.5 text-sm">
            <div>
              <p class="text-riad-800">{{ c.libelle }}</p>
              <p class="text-xs text-riad-400">{{ c.quantite }} × {{ c.prix_unitaire }} MAD</p>
            </div>
            <span class="font-semibold">{{ c.quantite * c.prix_unitaire }} MAD</span>
          </div>

          <!-- Total conso -->
          <div class="flex justify-between pt-3 border-t border-riad-100 text-sm">
            <span class="text-riad-600">Total consommations</span>
            <span class="font-semibold">{{ facture.total_consommations }} MAD</span>
          </div>

          <!-- Total séjour -->
          <div class="flex justify-between text-lg font-bold pt-2">
            <span>Total séjour</span>
            <span>{{ facture.total_sejour }} MAD</span>
          </div>
        </div>
      </div>

      <!-- Paiements -->
      <div class="card overflow-hidden">
        <div class="p-4 sm:p-5 border-b border-riad-100 flex items-center justify-between">
          <h3 class="font-semibold">Paiements</h3>
          <button @click="showPaiement = !showPaiement" class="btn-primary btn-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>
            Encaisser
          </button>
        </div>

        <!-- Paiement form -->
        <div v-if="showPaiement" class="p-4 sm:p-5 border-b border-riad-100 space-y-3 animate-in bg-riad-50">
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="sm:col-span-1">
              <label class="block text-xs font-semibold text-riad-700 mb-1">Montant</label>
              <input v-model.number="paiement.montant" type="number" class="input" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-riad-700 mb-1">Mode</label>
              <select v-model="paiement.mode_paiement" class="select">
                <option value="especes">Espèces</option>
                <option value="carte">Carte bancaire</option>
                <option value="virement">Virement</option>
                <option value="cheque">Chèque</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-riad-700 mb-1">Référence</label>
              <input v-model="paiement.reference" class="input" placeholder="Optionnel" />
            </div>
          </div>
          <div class="flex gap-3">
            <button @click="showPaiement = false" class="btn-ghost btn-sm">Annuler</button>
            <button @click="submitPaiement" class="btn-primary btn-sm">Valider le paiement</button>
          </div>
        </div>

        <div class="p-4 sm:p-5 space-y-2">
          <div v-if="facture.paiements.length === 0" class="text-center py-4 text-riad-400 text-sm">
            Aucun paiement enregistré
          </div>
          <div v-for="p in facture.paiements" :key="p.id"
            class="flex justify-between items-center p-2.5 rounded-lg bg-riad-50">
            <div>
              <p class="text-sm font-semibold text-riad-800">{{ p.mode_paiement }}</p>
              <p v-if="p.reference" class="text-xs text-riad-400">Ref: {{ p.reference }}</p>
            </div>
            <span class="font-bold text-green-600">- {{ p.montant }} MAD</span>
          </div>
        </div>
      </div>

      <!-- Reste dû -->
      <div class="card p-5">
        <div class="flex justify-between items-center">
          <span class="text-lg font-bold">Reste à payer</span>
          <span :class="['text-2xl font-bold', facture.restant_du > 0 ? 'text-red-600' : 'text-green-600']">
            {{ facture.restant_du > 0 ? facture.restant_du + ' MAD' : '✅ Solde' }}
          </span>
        </div>
      </div>
    </template>
  </div>
</template>
