<script setup>
import { reactive, ref } from 'vue'
import { useRiadStore } from '../stores/riad'
import { useAuthStore } from '../stores/auth'

const riad = useRiadStore()
const auth = useAuthStore()
const message = ref('')
const isSubmitting = ref(false)

const form = reactive({
  client_id: auth.user?.id || '',
  chambre_id: '',
  date_debut: '',
  date_fin: '',
  montant: 0
})

async function submitReservation() {
  message.value = ''
  isSubmitting.value = true
  try {
    const result = await riad.createReservation(form)
    message.value = result.synced
      ? 'Réservation confirmée et synchronisée !'
      : 'Enregistrée localement (en attente de synchronisation).'
    Object.assign(form, { client_id: auth.user?.id || '', chambre_id: '', date_debut: '', date_fin: '', montant: 0 })
  } catch (e) {
    message.value = 'Erreur: ' + e
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <!-- Header -->
    <div>
      <h2 class="text-xl font-bold text-riad-900">Nouvelle réservation</h2>
      <p class="text-riad-400 text-sm">Remplissez les détails pour réserver une chambre</p>
    </div>

    <div class="card p-6 sm:p-8">
      <form @submit.prevent="submitReservation" class="space-y-5">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Chambre</label>
            <select v-model="form.chambre_id" required class="select">
              <option value="">Sélectionner...</option>
              <option v-for="c in riad.libres" :key="c.id" :value="c.id">
                Chambre {{ c.numero }} - {{ c.type }} ({{ c.prix }} MAD)
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Montant (MAD)</label>
            <input v-model.number="form.montant" type="number" required class="input" />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Date de début</label>
            <input v-model="form.date_debut" type="date" required class="input" />
          </div>
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Date de fin</label>
            <input v-model="form.date_fin" type="date" required class="input" />
          </div>
        </div>

        <button type="submit" :disabled="isSubmitting"
          class="btn-primary w-full py-3">
          <span v-if="isSubmitting" class="flex items-center gap-2">
            <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
            Traitement...
          </span>
          <span v-else>Confirmer la réservation</span>
        </button>
      </form>

      <div v-if="message"
        class="mt-4 p-3 rounded-xl text-sm font-medium"
        :class="message.startsWith('Erreur') ? 'bg-red-50 text-red-700 border border-red-200' : 'bg-green-50 text-green-700 border border-green-200'">
        {{ message }}
      </div>
    </div>
  </div>
</template>
