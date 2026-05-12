<script setup>
import { reactive } from 'vue'
import { useRiadStore } from '../stores/riad'

const riadStore = useRiadStore()

const props = defineProps({
  reservation: { type: Object, required: false, default: null },
  isOpen: { type: Boolean, required: true }
})

const emit = defineEmits(['close', 'update'])

const form = reactive({ ...(props.reservation || {}) })

function handleSubmit() {
  emit('update', { ...form })
}
</script>

<template>
  <Teleport to="body">
    <div v-if="isOpen" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-riad-950/50 backdrop-blur-sm" @click.self="emit('close')">
      <div class="bg-white w-full max-w-lg rounded-2xl shadow-2xl animate-in overflow-hidden" @click.stop>
        <!-- Header -->
        <div class="flex items-center justify-between p-5 border-b border-riad-100">
          <div>
            <h3 class="text-lg font-bold text-riad-900">Modifier la réservation</h3>
            <p class="text-xs text-riad-400 mt-0.5">Mise à jour des détails du séjour</p>
          </div>
          <button @click="emit('close')" class="p-1.5 rounded-lg hover:bg-riad-100 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-riad-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="p-5 space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-semibold text-riad-700 mb-1.5">Chambre</label>
              <select v-model="form.chambre_id" required class="select">
                <option v-for="c in riadStore.rooms" :key="c.id" :value="c.id">Ch. {{ c.numero }} - {{ c.type }}</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-semibold text-riad-700 mb-1.5">Montant (MAD)</label>
              <input v-model.number="form.montant" type="number" required class="input" />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-semibold text-riad-700 mb-1.5">Date début</label>
              <input v-model="form.date_debut" type="date" required class="input" />
            </div>
            <div>
              <label class="block text-sm font-semibold text-riad-700 mb-1.5">Date fin</label>
              <input v-model="form.date_fin" type="date" required class="input" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Statut</label>
            <select v-model="form.statut" required class="select">
              <option value="confirmée">Confirmée</option>
              <option value="en attente">En attente</option>
              <option value="annulée">Annulée</option>
            </select>
          </div>

          <div class="flex items-center gap-3 pt-2">
            <button type="button" @click="emit('close')" class="btn-ghost flex-1 py-2.5">Annuler</button>
            <button type="submit" class="btn-primary flex-1 py-2.5">Enregistrer</button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>
