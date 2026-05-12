<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRiadStore } from '../stores/riad'
import ReservationModal from '../components/ReservationModal.vue'

const riad = useRiadStore()
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

const monthName = computed(() =>
  currentDate.value.toLocaleDateString('fr-MA', { month: 'long', year: 'numeric' })
)

function changeMonth(offset) {
  const d = new Date(currentDate.value)
  d.setMonth(d.getMonth() + offset)
  currentDate.value = d
}

function goToToday() { currentDate.value = new Date() }

function formatDate(date) { return date.toISOString().split('T')[0] }

function handleReservationClick(res) {
  if (!res) return
  selectedRes.value = { ...res }
  isModalOpen.value = true
}

async function handleUpdateReservation(updatedRes) {
  try {
    await riad.updateReservation(updatedRes)
    isModalOpen.value = false
  } catch (e) {
    alert('Erreur: ' + e)
  }
}

function resForCell(roomId, dateStr) {
  return riad.reservationsForRoomAndDate(roomId, dateStr)
}

function resBadgeClass(statut) {
  return statut === 'confirmée' ? 'bg-blue-500' : 'bg-amber-400'
}

onMounted(async () => {
  if (!riad.rooms.length) await riad.fetchChambres()
  if (!riad.reservations.length) await riad.fetchReservations()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Toolbar -->
    <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
      <div class="flex items-center gap-4">
        <h2 class="text-xl font-bold text-riad-900">{{ monthName }}</h2>
        <div class="flex items-center gap-1">
          <button @click="changeMonth(-1)" class="p-1.5 rounded-lg hover:bg-riad-100 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-riad-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5"/></svg>
          </button>
          <button @click="goToToday" class="px-3 py-1.5 text-xs font-semibold bg-gold-50 text-gold-700 rounded-lg hover:bg-gold-100 transition-colors border border-gold-200">
            Aujourd'hui
          </button>
          <button @click="changeMonth(1)" class="p-1.5 rounded-lg hover:bg-riad-100 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-riad-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5"/></svg>
          </button>
        </div>
      </div>
      <div class="flex items-center gap-3 text-xs text-riad-500">
        <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-sm bg-blue-500"></span> Confirmée</span>
        <span class="flex items-center gap-1.5"><span class="w-2.5 h-2.5 rounded-sm bg-amber-400"></span> En attente</span>
      </div>
    </div>

    <!-- Calendar -->
    <div class="overflow-x-auto rounded-xl border border-riad-200 bg-white shadow-sm">
      <table class="w-full border-collapse table-fixed min-w-[700px]">
        <thead>
          <tr>
            <th class="sticky left-0 z-10 bg-riad-950 p-3 text-left text-xs font-semibold text-gold-300 uppercase tracking-wider border-r border-riad-800 w-28">Chambres</th>
            <th v-for="day in daysInMonth" :key="day.getDate()"
              class="p-2 text-center text-[10px] font-semibold text-riad-500 border-r border-riad-100 w-10"
              :class="day.toDateString() === new Date().toDateString() ? 'bg-gold-500 text-white!' : 'bg-riad-50'">
              {{ day.getDate() }}<br><span class="opacity-60 font-normal lowercase">{{ day.toLocaleDateString('fr-MA', { weekday: 'short' }).slice(0,3) }}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="room in riad.rooms" :key="room.id" class="border-b border-riad-100 hover:bg-riad-50/50 transition-colors">
            <td class="sticky left-0 z-10 bg-white hover:bg-riad-50 p-3 text-sm font-semibold text-riad-800 border-r border-riad-100">
              Ch. {{ room.numero }}
              <span class="text-[10px] font-normal text-riad-400 ml-1">{{ room.type }}</span>
            </td>
            <td v-for="day in daysInMonth" :key="day.getDate()"
              class="p-0 border-r border-riad-100 relative h-12">
              <div v-if="resForCell(room.id, formatDate(day))"
                @click="handleReservationClick(resForCell(room.id, formatDate(day)))"
                :class="['absolute inset-y-1 left-0.5 right-0.5 rounded cursor-pointer transition-all hover:brightness-110 z-10 flex items-center justify-center text-[8px] font-bold text-white truncate shadow-sm', resBadgeClass(resForCell(room.id, formatDate(day)).statut)]">
                {{ resForCell(room.id, formatDate(day)).id.slice(0,4) }}
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <ReservationModal :isOpen="isModalOpen" :reservation="selectedRes" @close="isModalOpen = false" @update="handleUpdateReservation" />
  </div>
</template>
