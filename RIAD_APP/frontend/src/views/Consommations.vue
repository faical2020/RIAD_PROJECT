<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const reservationId = route.params.reservationId

const API = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1'
const headers = () => ({ 'Content-Type': 'application/json', Authorization: `Bearer ${localStorage.getItem('token')}` })

const reservation = ref(null)
const consommations = ref([])
const services = ref([])
const showForm = ref(false)
const form = ref({ service_id: '', libelle: '', quantite: 1, prix_unitaire: 0 })

async function load() {
  try {
    const r = await fetch(`${API}/reservations/${reservationId}/consommations`, { headers: headers() })
    if (r.ok) consommations.value = await r.json()
  } catch {}
  try {
    const s = await fetch(`${API}/services`, { headers: headers() })
    if (s.ok) services.value = await s.json()
  } catch {}
}

async function loadReservation() {
  try {
    const r = await fetch(`${API}/reservations/mine`, { headers: headers() })
    if (r.ok) {
      const all = await r.json()
      reservation.value = all.find(x => x.id === reservationId)
    }
  } catch {}
}

onMounted(() => { load(); loadReservation() })

const total = computed(() =>
  consommations.value.reduce((acc, c) => acc + c.quantite * c.prix_unitaire, 0)
)

function selectService(s) {
  form.value.service_id = s.id
  form.value.libelle = s.nom
  form.value.prix_unitaire = s.prix
}

async function addConso() {
  if (!form.value.libelle) return
  try {
    await fetch(`${API}/reservations/${reservationId}/consommations`, {
      method: 'POST', headers: headers(),
      body: JSON.stringify(form.value)
    })
    form.value = { service_id: '', libelle: '', quantite: 1, prix_unitaire: 0 }
    showForm.value = false
    load()
  } catch {}
}

async function removeConso(id) {
  if (!confirm('Supprimer cette consommation ?')) return
  try {
    await fetch(`${API}/consommations/${id}`, { method: 'DELETE', headers: headers() })
    load()
  } catch {}
}

function goToFacture() {
  router.push(`/app/facture/${reservationId}`)
}
</script>

<template>
  <div class="space-y-6 max-w-3xl">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-riad-900">Consommations</h2>
        <p v-if="reservation" class="text-riad-400 text-sm">
          Réservation du {{ reservation.date_debut }} → {{ reservation.date_fin }}
        </p>
      </div>
      <div class="flex gap-2">
        <button @click="showForm = !showForm" class="btn-primary btn-sm">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>
          Ajouter
        </button>
        <button @click="goToFacture" class="btn-ghost btn-sm">🧾 Facture</button>
      </div>
    </div>

    <!-- Add form -->
    <div v-if="showForm" class="card p-5 space-y-4 animate-in">
      <h3 class="font-semibold">Nouvelle consommation</h3>

      <!-- Quick select from catalog -->
      <div v-if="services.length" class="flex gap-2 flex-wrap">
        <button v-for="s in services" :key="s.id" @click="selectService(s)"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold border transition-colors"
          :class="form.service_id === s.id ? 'bg-riad-900 text-white border-riad-900' : 'bg-white text-riad-600 border-riad-200 hover:border-riad-300'">
          {{ s.nom }} ({{ s.prix }} MAD)
        </button>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div class="sm:col-span-2">
          <label class="block text-xs font-semibold text-riad-700 mb-1">Libellé</label>
          <input v-model="form.libelle" class="input" placeholder="Thé à la menthe" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-riad-700 mb-1">Quantité</label>
          <input v-model.number="form.quantite" type="number" min="1" class="input" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-riad-700 mb-1">Prix unitaire</label>
          <input v-model.number="form.prix_unitaire" type="number" class="input" />
        </div>
      </div>
      <div class="flex gap-3">
        <button @click="showForm = false" class="btn-ghost">Annuler</button>
        <button @click="addConso" class="btn-primary">Ajouter</button>
      </div>
    </div>

    <!-- List -->
    <div class="card overflow-hidden">
      <div class="p-4 sm:p-5 border-b border-riad-100 flex items-center justify-between">
        <h3 class="font-semibold">{{ consommations.length }} consommation(s)</h3>
        <span class="text-sm font-bold text-riad-900">{{ total }} MAD</span>
      </div>
      <div class="p-4 sm:p-5 space-y-2">
        <div v-if="consommations.length === 0" class="text-center py-8 text-riad-400 text-sm">
          Aucune consommation enregistrée
        </div>
        <div v-else v-for="c in consommations" :key="c.id"
          class="flex items-center justify-between p-3 rounded-xl bg-riad-50 hover:bg-riad-100/50 transition-colors">
          <div class="flex items-center gap-3">
            <div>
              <p class="text-sm font-semibold text-riad-800">{{ c.libelle }}</p>
              <p class="text-xs text-riad-400">{{ c.quantite }} × {{ c.prix_unitaire }} MAD</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-sm font-bold text-riad-900">{{ c.quantite * c.prix_unitaire }} MAD</span>
            <button @click="removeConso(c.id)" class="p-1 rounded-lg hover:bg-red-50 transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-red-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
          </div>
        </div>
        <div v-if="consommations.length" class="flex justify-between items-center pt-3 border-t border-riad-200">
          <span class="font-semibold">Total</span>
          <span class="font-bold text-lg text-riad-900">{{ total }} MAD</span>
        </div>
      </div>
    </div>
  </div>
</template>
