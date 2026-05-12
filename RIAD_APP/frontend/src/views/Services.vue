<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const API = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1'
const headers = () => ({ Authorization: `Bearer ${localStorage.getItem('token')}` })

const services = ref([])
const showForm = ref(false)
const editForm = ref({ nom: '', categorie: 'restaurant', description: '', prix: 0 })
const isEditing = ref(false)
const selectedId = ref(null)
const filterCat = ref('all')

async function load() {
  try {
    const r = await fetch(`${API}/services`, { headers: headers() })
    if (r.ok) services.value = await r.json()
  } catch {}
}

onMounted(load)

const categories = ['restaurant', 'spa', 'blanchisserie', 'extra']

const catMeta = {
  restaurant: { label: 'Restaurant', icon: '🍽️', color: 'bg-amber-50 text-amber-700 border-amber-200' },
  spa: { label: 'Spa & Bien-être', icon: '💆', color: 'bg-purple-50 text-purple-700 border-purple-200' },
  blanchisserie: { label: 'Blanchisserie', icon: '👕', color: 'bg-blue-50 text-blue-700 border-blue-200' },
  extra: { label: 'Extra', icon: '📦', color: 'bg-gray-50 text-gray-700 border-gray-200' },
}

const filteredServices = computed(() =>
  filterCat.value === 'all' ? services.value : services.value.filter(s => s.categorie === filterCat.value)
)

const grouped = computed(() => {
  const groups = {}
  for (const s of filteredServices.value) {
    const cat = s.categorie || 'extra'
    if (!groups[cat]) groups[cat] = []
    groups[cat].push(s)
  }
  return groups
})

function openAdd(cat) {
  editForm.value = { nom: '', categorie: cat || 'restaurant', description: '', prix: 0 }
  isEditing.value = false; showForm.value = true
}

function openEdit(s) {
  editForm.value = { ...s }
  isEditing.value = true; selectedId.value = s.id; showForm.value = true
}

async function save() {
  const url = isEditing.value
    ? `${API}/services/${selectedId.value}`
    : `${API}/services`
  const method = isEditing.value ? 'PUT' : 'POST'
  try {
    const r = await fetch(url, {
      method, headers: { 'Content-Type': 'application/json', ...headers() },
      body: JSON.stringify(editForm.value)
    })
    if (r.ok) { showForm.value = false; load() }
  } catch {}
}

async function remove(id) {
  if (!confirm('Supprimer ce service ?')) return
  try {
    await fetch(`${API}/services/${id}`, { method: 'DELETE', headers: headers() })
    load()
  } catch {}
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-riad-900">Catalogue des services</h2>
        <p class="text-riad-400 text-sm">{{ services.length }} service(s)</p>
      </div>
      <div class="flex gap-2 w-full sm:w-auto overflow-x-auto">
        <button v-for="c in [{ key: 'all', label: 'Tout' }, ...categories.map(c => ({ key: c, label: catMeta[c]?.icon + ' ' + catMeta[c]?.label }))]"
          :key="c.key" @click="filterCat = c.key"
          :class="['px-3 py-1.5 rounded-lg text-xs font-semibold border whitespace-nowrap transition-colors',
            filterCat === c.key ? 'bg-riad-900 text-white border-riad-900' : 'bg-white text-riad-600 border-riad-200 hover:border-riad-300']">
          {{ c.label }}
        </button>
      </div>
    </div>

    <!-- Form -->
    <div v-if="showForm" class="card p-5 space-y-4 animate-in">
      <h3 class="font-semibold text-riad-900">{{ isEditing ? 'Modifier' : 'Nouveau' }} service</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
        <div>
          <label class="block text-xs font-semibold text-riad-700 mb-1">Nom</label>
          <input v-model="editForm.nom" class="input" placeholder="Thé à la menthe" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-riad-700 mb-1">Catégorie</label>
          <select v-model="editForm.categorie" class="select">
            <option v-for="c in categories" :key="c" :value="c">{{ catMeta[c]?.label || c }}</option>
          </select>
        </div>
        <div>
          <label class="block text-xs font-semibold text-riad-700 mb-1">Prix (MAD)</label>
          <input v-model.number="editForm.prix" type="number" class="input" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-riad-700 mb-1">Description</label>
          <input v-model="editForm.description" class="input" placeholder="Optionnel" />
        </div>
      </div>
      <div class="flex gap-3">
        <button @click="showForm = false" class="btn-ghost btn-sm">Annuler</button>
        <button @click="save" class="btn-primary btn-sm">Enregistrer</button>
      </div>
    </div>

    <!-- Add button -->
    <button v-if="auth.isAdmin" @click="openAdd(filterCat === 'all' ? 'restaurant' : filterCat)"
      class="btn-primary btn-sm">
      <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>
      Ajouter un service
    </button>

    <!-- Grouped by category -->
    <div v-for="(items, cat) in grouped" :key="cat" class="space-y-3">
      <div class="flex items-center gap-2">
        <span class="text-lg">{{ catMeta[cat]?.icon || '📦' }}</span>
        <h3 class="font-semibold text-riad-900">{{ catMeta[cat]?.label || cat }}</h3>
        <span class="text-xs text-riad-400">({{ items.length }})</span>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
        <div v-for="s in items" :key="s.id"
          class="card p-4 hover:-translate-y-0.5 animate-in relative group">
          <div class="flex items-start justify-between mb-2">
            <div>
              <h4 class="font-semibold text-riad-900 text-sm">{{ s.nom }}</h4>
              <p v-if="s.description" class="text-xs text-riad-400 mt-0.5">{{ s.description }}</p>
            </div>
            <span :class="['badge text-[10px]', catMeta[s.categorie]?.color || 'badge-gray']">
              {{ catMeta[s.categorie]?.icon || '' }}
            </span>
          </div>
          <div class="flex items-center justify-between mt-3 pt-3 border-t border-riad-100">
            <span class="text-base font-bold text-gold-600">{{ s.prix }} MAD</span>
            <div v-if="auth.isAdmin" class="flex gap-1 opacity-0 group-hover:opacity-100 transition">
              <button @click="openEdit(s)" class="p-1 rounded hover:bg-riad-100 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-riad-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"/></svg>
              </button>
              <button @click="remove(s.id)" class="p-1 rounded hover:bg-red-50 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-3.5 h-3.5 text-red-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="Object.keys(grouped).length === 0" class="text-center py-12 text-riad-400">
      <p class="text-3xl mb-2">📦</p>
      <p class="text-sm">Aucun service dans cette catégorie</p>
    </div>
  </div>
</template>
