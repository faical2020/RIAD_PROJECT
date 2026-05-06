<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRiadStore } from '../stores/riad'
import { useRouter } from 'vue-router'
import ReservationModal from '../components/ReservationModal.vue'

const riad = useRiadStore()
const router = useRouter()

const currentDate = ref(new Date())
const selectedRes = ref(null)
const isModalOpen = ref(false)

const viewedMonth = computed(() => currentDate.value.getMonth())
const viewedYear = computed(() => currentDate.value.getFullYear())

const daysInMonth = computed(() => {
    const year = viewedYear.value
    const month = viewedMonth.value
    const date = new Date(year, month, 1)
    const days = []
    while (date.getMonth() === month) {
        days.push(new Date(date))
        date.setDate(date.getDate() + 1)
    }
    return days
})

const monthName = computed(() => {
    return currentDate.value.toLocaleDateString('fr-MA', { month: 'long', year: 'numeric' })
})

function changeMonth(offset) {
    const newDate = new Date(currentDate.value)
    newDate.setMonth(newDate.getMonth() + offset)
    currentDate.value = newDate
}

function goToToday() {
    currentDate.value = new Date()
}

function formatDate(date) {
    return date.toISOString().split('T')[0]
}

async function handleReservationClick(res) {
    if (!res) return
    selectedRes.value = { ...res }
    isModalOpen.value = true
}

async function handleUpdateReservation(updatedRes) {
    try {
        await riad.updateReservation(updatedRes)
        isModalOpen.value = false
    } catch (e) {
        alert('Erreur lors de la mise à jour: ' + e)
    }
}

onMounted(async () => {
    if (riad.rooms.length === 0) {
        await riad.fetchChambres()
    }
    if (riad.reservations.length === 0) {
        await riad.fetchReservations()
    }
})
</script>

<template>
    <div class="space-y-6">
        <!-- Toolbar -->
        <div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white p-4 rounded-2xl shadow-sm border border-riad-100">
            <div class="flex items-center gap-4">
                <h2 class="font-display text-xl font-bold text-riad-900">{{ monthName }}</h2>
                <div class="flex items-center gap-2">
                    <button @click="changeMonth(-1)" class="p-2 rounded-lg hover:bg-riad-100 text-riad-600 transition-colors">
                        <span class="text-lg">◀</span>
                    </button>
                    <button @click="goToToday" class="px-3 py-1 text-xs font-bold uppercase tracking-wider bg-gold-100 text-gold-700 rounded-lg hover:bg-gold-200 transition-colors">
                        Aujourd'hui
                    </button>
                    <button @click="changeMonth(1)" class="p-2 rounded-lg hover:bg-riad-100 text-riad-600 transition-colors">
                        <span class="text-lg">▶</span>
                    </button>
                </div>
            </div>
            <div class="flex items-center gap-3 text-xs font-medium">
                <div class="flex items-center gap-1.5">
                    <span class="w-3 h-3 rounded-full bg-blue-500"></span>
                    <span class="text-riad-500">Confirmée</span>
                </div>
                <div class="flex items-center gap-1.5">
                    <span class="w-3 h-3 rounded-full bg-amber-400"></span>
                    <span class="text-riad-500">En attente</span>
                </div>
            </div>
        </div>

        <!-- Calendar Grid -->
        <div class="relative overflow-x-auto shadow-xl rounded-2xl border border-riad-100 bg-white">
            <table class="w-full border-collapse table-fixed min-w-[800px]">
                <thead>
                    <tr class="bg-riad-950 text-gold-300">
                        <th class="sticky left-0 z-20 bg-riad-950 p-3 text-left text-xs font-bold uppercase tracking-wider border-r border-riad-800 w-32">
                            Chambres
                        </th>
                        <th v-for="day in daysInMonth" :key="day.getDate()" 
                            class="p-2 text-center text-[10px] font-bold uppercase tracking-tighter border-r border-riad-800 w-10"
                            :class="day.getDate() === new Date().getDate() && day.getMonth() === new Date().getMonth() ? 'bg-gold-600 text-white' : ''">
                            {{ day.getDate() }}<br/>
                            <span class="opacity-60 font-normal lowercase">{{ day.toLocaleDateString('fr-MA', { weekday: 'short' }).slice(0,3) }}</span>
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="room in riad.rooms" :key="room.id" class="border-b border-riad-100 group hover:bg-riad-50/50 transition-colors">
                        <td class="sticky left-0 z-10 bg-white group-hover:bg-riad-50 p-3 text-sm font-bold text-riad-800 border-r border-riad-100 whitespace-nowrap">
                            Ch. {{ room.numero }} <span class="text-[10px] font-normal text-riad-400 ml-1">{{ room.type }}</span>
                        </td>
                        <td v-for="day in daysInMonth" :key="day.getDate()" 
                            class="p-0 border-r border-riad-100 relative h-12 group/cell">
                            
                            <div v-if="riad.reservationsForRoomAndDate(room.id, formatDate(day))"
                                @click="handleReservationClick(riad.reservationsForRoomAndDate(room.id, formatDate(day)))"
                                class="absolute inset-y-1 left-0 right-0 rounded-md cursor-pointer transition-all duration-200 hover:brightness-110 z-10 flex items-center justify-center overflow-hidden px-1 text-[9px] font-bold text-white truncate shadow-sm"
                                :class="riad.reservationsForRoomAndDate(room.id, formatDate(day)).statut === 'confirmée' ? 'bg-blue-500' : 'bg-amber-400'">
                                {{ riad.reservationsForRoomAndDate(room.id, formatDate(day)).id.slice(0,4) }}
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <ReservationModal 
            :isOpen="isModalOpen" 
            :reservation="selectedRes" 
            @close="isModalOpen = false"
            @update="handleUpdateReservation"
        />
    </div>
</template>

<style scoped>
.scrollbar-thin::-webkit-scrollbar {
    height: 6px;
    width: 6px;
}
.scrollbar-thin::-webkit-scrollbar-track {
    background: transparent;
}
.scrollbar-thin::-webkit-scrollbar-thumb {
    background: #d1d5db;
    border-radius: 10px;
}
</style>
