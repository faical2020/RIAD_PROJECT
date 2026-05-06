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
          ? '✅ Réservation confirmée et synchronisée !'
          : '⏳ Enregistrée localement (en attente de synchronisation).'
        
        // Reset form regardless of sync status
        Object.assign(form, { client_id: auth.user?.id || '', chambre_id: '', date_debut: '', date_fin: '', montant: 0 })
    } catch (e) {
    message.value = '❌ Erreur: ' + e
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="card p-6 sm:p-8 max-w-2xl mx-auto">
    <div class="text-center mb-6 sm:mb-8">
      <h2 class="font-display text-xl sm:text-2xl text-riad-900 mb-2">Nouvelle Réservation</h2>
      <p class="text-riad-500 text-sm">Remplissez les détails pour réserver une chambre</p>
    </div>

    <form @submit.prevent="submitReservation" class="space-y-4 sm:space-y-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
        <div class="group">
          <label class="block text-riad-700 text-sm font-bold mb-1.5">Chambre</label>
          <select v-model="form.chambre_id" required
            class="w-full bg-white border border-riad-200 rounded-xl px-4 py-2.5 text-riad-900 focus:outline-none focus:ring-2 focus:ring-gold-500/50 transition-all duration-300 hover:border-riad-300">
            <option value="">Sélectionner...</option>
            <option v-for="c in riad.libres" :key="c.id" :value="c.id">
              Chambre {{ c.numero }} - {{ c.type }} ({{ c.prix }} MAD)
            </option>
          </select>
        </div>
        <div class="group">
          <label class="block text-riad-700 text-sm font-bold mb-1.5">Montant (MAD)</label>
          <input v-model.number="form.montant" type="number" required
            class="w-full bg-white border border-riad-200 rounded-xl px-4 py-2.5 text-riad-900 focus:outline-none focus:ring-2 focus:ring-gold-500/50 transition-all duration-300 hover:border-riad-300" />
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
        <div class="group">
          <label class="block text-riad-700 text-sm font-bold mb-1.5">Date de début</label>
          <input v-model="form.date_debut" type="date" required
            class="w-full bg-white border border-riad-200 rounded-xl px-4 py-2.5 text-riad-900 focus:outline-none focus:ring-2 focus:ring-gold-500/50 transition-all duration-300 hover:border-riad-300" />
        </div>
        <div class="group">
          <label class="block text-riad-700 text-sm font-bold mb-1.5">Date de fin</label>
          <input v-model="form.date_fin" type="date" required
            class="w-full bg-white border border-riad-200 rounded-xl px-4 py-2.5 text-riad-900 focus:outline-none focus:ring-2 focus:ring-gold-500/50 transition-all duration-300 hover:border-riad-300" />
        </div>
      </div>

      <button type="submit" :disabled="isSubmitting"
        class="w-full bg-gradient-to-r from-gold-500 to-gold-400 hover:from-gold-400 hover:to-gold-300 disabled:opacity-50 text-riad-950 font-body font-bold py-3 sm:py-3.5 rounded-xl transition-all duration-300 shadow-lg shadow-gold-500/20 hover:shadow-xl hover:shadow-gold-500/30 active:scale-[0.98]">
        <span v-if="!isSubmitting">Confirmer la Réservation</span>
        <span v-else class="flex items-center justify-center gap-2">
          <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-riad-950"></div>
          Traitement...
        </span>
      </button>
    </form>

    <div v-if="message"
      class="mt-4 sm:mt-6 p-3 sm:p-4 rounded-xl text-center text-sm font-medium transition-all duration-300"
      :class="message.startsWith('❌') ? 'bg-terracotta-600/20 text-terracotta-400 border border-terracotta-400/40' : 'bg-green-900/30 text-green-400 border border-green-500/40'">
      {{ message }}
    </div>
  </div>
</template>
