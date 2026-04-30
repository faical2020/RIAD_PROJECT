<script setup>
import { computed } from 'vue'

const props = defineProps({
  stats: {
    type: Object,
    required: true
  }
})

const statCards = computed(() => [
  { label: 'Total Chambres', value: props.stats.totalRooms, color: 'bg-blue-50 text-blue-600', icon: '🏨' },
  { label: 'Disponibles', value: props.stats.availableRooms, color: 'bg-green-50 text-green-600', icon: '✅' },
  { label: 'Occupées', value: props.stats.occupiedRooms, color: 'bg-red-50 text-red-600', icon: '🚫' },
  { label: 'Sync Pendante', value: props.stats.pendingSync, color: 'bg-amber-50 text-amber-600', icon: '⏳' },
])
</script>

<template>
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
    <div v-for="card in statCards" :key="card.label" 
         class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm flex items-center gap-4">
      <div :class="['text-2xl p-3 rounded-lg', card.color]">
        {{ card.icon }}
      </div>
      <div>
        <p class="text-xs font-medium text-gray-500 uppercase tracking-wider">{{ card.label }}</p>
        <p class="text-xl font-bold text-gray-900">{{ card.value }}</p>
      </div>
    </div>
  </div>
</template>
