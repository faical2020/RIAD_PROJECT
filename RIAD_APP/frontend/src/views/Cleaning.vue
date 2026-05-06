<script setup>
import { ref, onMounted } from 'vue'
import { useRiadStore } from '../stores/riad'

const riad = useRiadStore()

onMounted(async () => {
    await riad.fetchChambres()
})

async function updateStatus(roomId, newStatus) {
    try {
        await riad.updateCleaningStatus(roomId, newStatus)
    } catch (e) {
        alert('Erreur lors de la mise à jour du statut: ' + e)
    }
}
</script>

<template>
    <div class="space-y-6">
        <div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white p-4 rounded-2xl shadow-sm border border-riad-100">
            <div>
                <h2 class="font-display text-xl font-bold text-riad-900">Gestion du Ménage</h2>
                <p class="text-riad-500 text-sm">Mise à jour du statut de propreté des chambres</p>
            </div>
            <div class="flex items-center gap-4 text-xs font-medium">
                <div class="flex items-center gap-1.5">
                    <span class="w-3 h-3 rounded-full bg-green-500"></span>
                    <span class="text-riad-500">Propre</span>
                </div>
                <div class="flex items-center gap-1.5">
                    <span class="w-3 h-3 rounded-full bg-amber-400"></span>
                    <span class="text-riad-500">En cours</span>
                </div>
                <div class="flex items-center gap-1.5">
                    <span class="w-3 h-3 rounded-full bg-red-500"></span>
                    <span class="text-riad-500">Sale</span>
                </div>
            </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <div v-for="room in riad.rooms" :key="room.id" 
                class="bg-white p-5 rounded-2xl border border-riad-100 shadow-sm hover:shadow-md transition-all duration-300 group">
                <div class="flex justify-between items-start mb-4">
                    <div>
                        <h3 class="font-bold text-riad-900">Chambre {{ room.numero }}</h3>
                        <p class="text-riad-400 text-xs">{{ room.type }}</p>
                    </div>
                    <span :class="[
                        'px-2 py-1 rounded-lg text-[10px] font-bold uppercase tracking-wider border transition-colors',
                        room.cleaning_status === 'propre' ? 'bg-green-50 text-green-600 border-green-200' : 
                        room.cleaning_status === 'en cours' ? 'bg-amber-50 text-amber-600 border-amber-200' : 
                        'bg-red-50 text-red-600 border-red-200'
                    ]">
                        {{ room.cleaning_status }}
                    </span>
                </div>

                <div class="grid grid-cols-3 gap-2">
                    <button @click="updateStatus(room.id, 'propre')" 
                        class="py-2 px-1 rounded-lg text-[10px] font-bold transition-all duration-200 border border-green-200 text-green-600 hover:bg-green-500 hover:text-white">
                        Propre
                    </button>
                    <button @click="updateStatus(room.id, 'en cours')" 
                        class="py-2 px-1 rounded-lg text-[10px] font-bold transition-all duration-200 border border-amber-200 text-amber-600 hover:bg-amber-500 hover:text-white">
                        En cours
                    </button>
                    <button @click="updateStatus(room.id, 'sale')" 
                        class="py-2 px-1 rounded-lg text-[10px] font-bold transition-all duration-200 border border-red-200 text-red-600 hover:bg-red-500 hover:text-white">
                        Sale
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
