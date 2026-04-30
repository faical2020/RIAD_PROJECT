<script setup>
import { ref, onMounted } from 'vue'
import RoomList from './components/RoomList.vue'
import ReservationForm from './components/ReservationForm.vue'
import StatsBar from './components/StatsBar.vue'
import Login from './components/Login.vue'
import Register from './components/Register.vue'
import { riadService } from './services/riadService'

const isAuthenticated = ref(!!localStorage.getItem('token'))
const authMode = ref('login') 
const stats = ref({ totalRooms: 0, occupiedRooms: 0, availableRooms: 0, pendingSync: 0 })
const activeTab = ref('rooms') // 'rooms' or 'reservations'

const handleLoginSuccess = (user) => {
  isAuthenticated.value = true
  updateStats()
}

const toggleAuthMode = () => {
  authMode.value = authMode.value === 'login' ? 'register' : 'login'
}

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('role')
  isAuthenticated.value = false
}

const updateStats = async () => {
  stats.value = await riadService.getDashboardStats()
}

onMounted(() => {
  if (isAuthenticated.value) {
    updateStats()
    // Update stats every 30 seconds
    setInterval(updateStats, 30000)
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 font-sans">
    <!-- Navbar -->
    <nav v-if="isAuthenticated" class="bg-white border-b border-gray-200 py-3 px-6 flex justify-between items-center sticky top-0 z-20 shadow-sm">
      <div class="flex items-center gap-4">
        <h1 class="text-xl font-bold text-emerald-700 tracking-tight">RIAD <span class="text-gray-400 font-light">Management</span></h1>
        <div class="hidden md:flex items-center gap-2 px-3 py-1 bg-emerald-50 text-emerald-600 rounded-full text-xs font-medium border border-emerald-100">
          <div class="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
          Connecté
        </div>
      </div>
      
      <div class="flex gap-2 items-center">
        <button @click="activeTab = 'rooms'" 
                :class="['px-4 py-2 text-sm font-medium rounded-lg transition-all', activeTab === 'rooms' ? 'bg-emerald-600 text-white' : 'text-gray-500 hover:bg-gray-100']">
          Chambres
        </button>
        <button @click="activeTab = 'reservations'" 
                :class="['px-4 py-2 text-sm font-medium rounded-lg transition-all', activeTab === 'reservations' ? 'bg-emerald-600 text-white' : 'text-gray-500 hover:bg-gray-100']">
          Réservations
        </button>
        <div class="w-px h-6 bg-gray-200 mx-2"></div>
        <button @click="logout" class="p-2 text-gray-400 hover:text-red-500 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </button>
      </div>
    </nav>

    <main class="container mx-auto py-8 px-4">
      <!-- Auth Screens -->
      <div v-if="!isAuthenticated" class="max-w-md mx-auto py-12">
        <div class="text-center mb-8">
          <div class="w-16 h-16 bg-emerald-600 text-white rounded-2xl flex items-center justify-center text-3xl font-bold mx-auto mb-4 shadow-lg">R</div>
          <h1 class="text-3xl font-bold text-gray-800 mb-2">Bienvenue au Riad</h1>
          <p class="text-gray-500">Gestion hybride Local-First & Cloud</p>
        </div>

        <Login v-if="authMode === 'login'" @login-success="handleLoginSuccess" />
        <Register v-else />

        <div class="text-center mt-6">
          <button @click="toggleAuthMode" class="text-sm text-emerald-600 font-medium hover:text-emerald-700 transition-colors">
            {{ authMode === 'login' ? "Pas de compte ? Inscrivez-vous" : "Déjà un compte ? Connectez-vous" }}
          </button>
        </div>
      </div>

      <!-- Dashboard Content -->
      <div v-else class="space-y-8">
        <!-- Stats Header -->
        <StatsBar :stats="stats" />

        <!-- Tab Content -->
        <div v-if="activeTab === 'rooms'" class="animate-in fade-in slide-in-from-bottom-4 duration-500">
          <RoomList />
        </div>
        <div v-else class="animate-in fade-in slide-in-from-bottom-4 duration-500">
          <ReservationForm />
        </div>
      </div>
    </main>

    <footer v-if="isAuthenticated" class="bg-white border-t border-gray-200 py-6 text-center text-gray-400 text-xs">
      &copy; 2026 RIAD_APP Hybrid System. Local-first, Cloud-synced.
    </footer>
  </div>
</template>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>
