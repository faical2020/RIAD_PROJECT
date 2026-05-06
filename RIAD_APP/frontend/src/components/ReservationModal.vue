<script setup>
import { ref, reactive } from 'vue'
import { useRiadStore } from '../stores/riad'

const riadStore = useRiadStore()

const props = defineProps({
  reservation: {
    type: Object,
    required: false,
    default: null
  },
  isOpen: {
    type: Boolean,
    required: true
  }
})

const emit = defineEmits(['close', 'update'])

const form = reactive({ ...(props.reservation || {}) })

async function handleSubmit() {
  emit('update', { ...form })
}
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-riad-950/60 backdrop-blur-sm transition-opacity duration-300">
    <div class="bg-white w-full max-w-lg rounded-3xl shadow-2xl overflow-hidden animate-in fade-in zoom-in duration-200">
      <div class="bg-gradient-to-r from-gold-600 to-gold-500 p-6 text-white flex justify-between items-center">
        <div>
          <h3 class="font-display text-xl font-bold">Modifier la Réservation</h3>
          <p class="text-gold-100 text-xs opacity-80">Mise à jour des détails du séjour</p>
        </div>
        <button @click="emit('close')" class="p-2 rounded-full hover:bg-white/20 transition-colors">
          <span class="text-2xl">&times;</span>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-6 space-y-5">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="group">
            <label class="block text-riad-700 text-xs font-bold mb-1.5 uppercase tracking-wider">Chambre</label>
            <select v-model="form.chambre_id" required
              class="w-full bg-riad-50 border border-riad-200 rounded-xl px-4 py-2 text-riad-900 focus:ring-2 focus:ring-gold-500/50 outline-none transition-all">
              <option v-for="c in riadStore.rooms" :key="c.id" :value="c.id">
                Ch. {{ c.numero }} - {{ c.type }}
              </option>
            </select>
          </div>
          <div class="group">
            <label class="block text-riad-700 text-xs font-bold mb-1.5 uppercase tracking-wider">Montant (MAD)</label>
            <input v-model.number="form.montant" type="number" required
              class="w-full bg-riad-50 border border-riad-200 rounded-xl px-4 py-2 text-riad-900 focus:ring-2 focus:ring-gold-500/50 outline-none transition-all" />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="group">
            <label class="block text-riad-700 text-xs font-bold mb-1.5 uppercase tracking-wider">Date début</label>
            <input v-model="form.date_debut" type="date" required
              class="w-full bg-riad-50 border border-riad-200 rounded-xl px-4 py-2 text-riad-900 focus:ring-2 focus:ring-gold-500/50 outline-none transition-all" />
          </div>
          <div class="group">
            <label class="block text-riad-700 text-xs font-bold mb-1.5 uppercase tracking-wider">Date fin</label>
            <input v-model="form.date_fin" type="date" required
              class="w-full bg-riad-50 border border-riad-200 rounded-xl px-4 py-2 text-riad-900 focus:ring-2 focus:ring-gold-500/50 outline-none transition-all" />
          </div>
        </div>

        <div class="group">
          <label class="block text-riad-700 text-xs font-bold mb-1.5 uppercase tracking-wider">Statut</label>
          <select v-model="form.statut" required
            class="w-full bg-riad-50 border border-riad-200 rounded-xl px-4 py-2 text-riad-900 focus:ring-2 focus:ring-gold-500/50 outline-none transition-all">
            <option value="confirmée">Confirmée</option>
            <option value="en attente">En attente</option>
            <option value="annulée">Annulée</option>
          </select>
        </div>

        <div class="flex items-center gap-3 pt-4">
          <button type="button" @click="emit('close')"
            class="flex-1 py-3 px-4 rounded-xl text-riad-500 font-bold hover:bg-riad-100 transition-all duration-200">
            Annuler
          </button>
          <button type="submit"
            class="flex-1 py-3 px-4 rounded-xl bg-gradient-to-r from-gold-600 to-gold-500 text-white font-bold shadow-lg shadow-gold-500/30 hover:scale-[1.02] transition-all duration-200">
            Enregistrer
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
