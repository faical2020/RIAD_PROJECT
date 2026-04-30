<script setup>
import { ref, onMounted } from 'vue'
import { riadService } from '../services/riadService'

const form = ref({
  userId: '',
  roomId: '',
  start: '',
  end: '',
  amount: 0
})
const message = ref('')
const isSubmitting = ref(false)

const submitReservation = async () => {
  message.value = ''
  isSubmitting.value = true
  try {
    const result = await riadService.createReservation(form.value)
    message.value = result.synced 
      ? '✅ Réservation confirmée et synchronisée !' 
      : '⏳ Enregistrée localement (en attente de synchronisation).'
    if (result.synced) {
      // Reset form on success
      form.value = { userId: '', roomId: '', start: '', end: '', amount: 0 }
    }
  } catch (e) {
    message.value = '❌ Erreur: ' + e
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="p-6 max-w-2xl mx-auto">
    <div class="bg-white rounded-2xl shadow-lg p-8 border border-gray-100">
      <div class="text-center mb-8">
        <h2 class="text-2xl font-bold text-gray-800">Effectuer une Réservation</h2>
        <p class="text-gray-500 text-sm">Remplissez les détails pour bloquer une chambre</p>
      </div>
      
      <form @submit.prevent="submitReservation" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="space-y-1">
            <label class="text-sm font-medium text-gray-700">ID Utilisateur</label>
            <input v-model="form.userId" type="text" required 
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
          </div>
          <div class="space-y-1">
            <label class="text-sm font-medium text-gray-700">ID Chambre</label>
            <input v-model="form.roomId" type="text" required 
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="space-y-1">
            <label class="text-sm font-medium text-gray-700">Date de début</label>
            <input v-model="form.start" type="date" required 
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
          </div>
          <div class="space-y-1">
            <label class="text-sm font-medium text-gray-700">Date de fin</label>
            <input v-model="form.end" type="date" required 
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
          </div>
        </div>

        <div class="space-y-1">
          <label class="text-sm font-medium text-gray-700">Montant Total (DH)</label>
          <input v-model.number="form.amount" type="number" required 
                 class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
        </div>

        <button type="submit" :disabled="isSubmitting"
                class="w-full py-3 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold rounded-lg transition-all shadow-md disabled:bg-gray-400 disabled:cursor-not-allowed">
          <span v-if="!isSubmitting">Confirmer la Réservation</span>
          <span v-else class="flex items-center justify-center gap-2">
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
            Traitement...
          </span>
        </button>
      </form>

      <div v-if="message" 
           class="mt-6 p-4 rounded-xl text-center text-sm font-medium transition-all"
           :class="message.startsWith('❌') ? 'bg-red-50 text-red-600 border border-red-100' : 'bg-emerald-50 text-emerald-600 border border-emerald-100'">
        {{ message }}
      </div>
    </div>
  </div>
</template>
