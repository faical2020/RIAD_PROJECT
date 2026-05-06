<template>
  <div class="min-h-screen flex bg-gradient-to-br from-riad-50 via-white to-gold-50/30">
    <!-- Sidebar -->
    <aside :class="['fixed inset-y-0 left-0 z-50 flex flex-col bg-riad-950/95 backdrop-blur-md transition-all duration-300 ease-in-out shadow-2xl',
      sidebarOpen ? 'w-64' : 'w-0 lg:w-20 overflow-hidden']">
      <!-- Logo -->
      <div class="flex items-center gap-3 px-4 py-6 border-b border-riad-800/50 min-w-0">
        <div class="flex-shrink-0 w-10 h-10 rounded-xl bg-gradient-to-br from-gold-500/30 to-gold-600/20 border border-gold-500/30 flex items-center justify-center shadow-lg shadow-gold-500/10">
          <span class="text-2xl">🌙</span>
        </div>
        <div v-show="sidebarOpen" class="overflow-hidden">
          <p class="font-display text-gold-300 text-base font-semibold whitespace-nowrap tracking-wide">Riad Manager</p>
          <p class="text-riad-500 text-xs whitespace-nowrap font-arabic mt-0.5">نظام الإدارة</p>
        </div>
      </div>

      <!-- Nav -->
      <nav class="flex-1 p-3 space-y-1.5 overflow-y-auto scrollbar-thin">
        <RouterLink v-for="item in navItems" :key="item.to" :to="item.to"
          :class="['sidebar-link group relative', { 'active': $route.name === item.name }]">
          <span class="text-xl flex-shrink-0 group-hover:scale-110 transition-transform duration-200">{{ item.icon }}</span>
          <span v-show="sidebarOpen" class="whitespace-nowrap font-medium">{{ item.label }}</span>
          <div v-show="sidebarOpen && $route.name === item.name" class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-6 bg-gold-500 rounded-r-full"></div>
        </RouterLink>
      </nav>

      <!-- User -->
      <div class="p-3 border-t border-riad-800/50 bg-riad-950/50">
        <div class="flex items-center gap-3 px-2 py-3">
          <div class="w-10 h-10 rounded-full bg-gradient-to-br from-gold-500/40 to-gold-600/20 border-2 border-gold-500/50 flex items-center justify-center flex-shrink-0 text-gold-300 font-bold text-base shadow-lg shadow-gold-500/10">
            {{ userInitial }}
          </div>
          <div v-show="sidebarOpen" class="overflow-hidden">
            <p class="text-riad-200 text-sm font-bold whitespace-nowrap truncate">{{ authStore.user?.prenom }} {{ authStore.user?.nom }}</p>
            <p class="text-riad-500 text-xs whitespace-nowrap capitalize">{{ authStore.role }}</p>
          </div>
        </div>
        <button @click="handleLogout"
          class="sidebar-link w-full mt-2 text-terracotta-400 hover:bg-terracotta-600/10 hover:text-terracotta-300">
          <span class="text-xl flex-shrink-0">🚪</span>
          <span v-show="sidebarOpen" class="font-medium">Déconnexion</span>
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <div :class="['flex-1 flex flex-col transition-all duration-300 min-w-0', sidebarOpen ? 'lg:ml-64' : 'lg:ml-20']">
      <!-- Top bar -->
      <header class="glass-header px-4 sm:px-6 py-4 flex items-center gap-4 sticky top-0 z-40">
        <button @click="sidebarOpen = !sidebarOpen" class="p-2 rounded-xl hover:bg-riad-800/50 transition-all duration-300 group">
          <div class="w-5 h-4 flex flex-col justify-between">
            <span class="block h-0.5 bg-riad-300 group-hover:bg-gold-400 transition-colors"></span>
            <span class="block h-0.5 bg-riad-300 group-hover:bg-gold-400 transition-colors"></span>
            <span class="block h-0.5 bg-riad-300 group-hover:bg-gold-400 transition-colors"></span>
          </div>
        </button>
        <div>
          <h1 class="font-display text-riad-100 text-lg sm:text-xl font-semibold">{{ currentTitle }}</h1>
          <p class="text-riad-500 text-xs sm:text-sm">{{ today }}</p>
        </div>
        <div class="ml-auto flex items-center gap-2 sm:gap-3">
          <span class="badge bg-gold-500/20 text-gold-400 capitalize border border-gold-500/30 hidden sm:inline-block">{{ authStore.role }}</span>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 p-4 sm:p-6 lg:p-8 overflow-auto max-w-[1600px] w-full mx-auto">
        <RouterView />
      </main>

      <!-- Footer -->
      <footer class="glass-header border-t border-riad-800/50 py-4 px-6 text-center">
        <p class="text-riad-500 text-xs">&copy; 2026 RIAD_APP. Digital Luxury Experience.</p>
      </footer>
    </div>

    <!-- Mobile overlay -->
    <div v-if="sidebarOpen" @click="sidebarOpen = false"
      class="lg:hidden fixed inset-0 bg-riad-950/80 backdrop-blur-sm z-40 transition-opacity duration-300"></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const sidebarOpen = ref(window.innerWidth >= 1024)

  const allNavItems = [
    { to: '/', name: 'Dashboard', icon: '🏠', label: 'Tableau de bord', roles: null },
    { to: '/calendrier', name: 'Calendar', icon: '📅', label: 'Planning', roles: ['manager', 'receptionniste'] },
    { to: '/cleaning', name: 'Cleaning', icon: '🧹', label: 'Ménage', roles: ['manager', 'receptionniste'] },
    { to: '/chambres', name: 'Chambres', icon: '🛏️', label: 'Chambres', roles: null },
    { to: '/reservations', name: 'Reservations', icon: '📋', label: 'Réservations', roles: ['manager', 'receptionniste'] },
    { to: '/nouvelle-reservation', name: 'NouvelleReservation', icon: '➕', label: 'Réserver', roles: null },
    { to: '/profil', name: 'Profil', icon: '👤', label: 'Mon profil', roles: null },
  ]

const navItems = computed(() =>
  allNavItems.filter(i => !i.roles || i.roles.includes(authStore.role))
)

const currentTitle = computed(() => {
  const item = allNavItems.find(i => i.name === route.name)
  return item?.label || 'Tableau de bord'
})

const userInitial = computed(() =>
  (authStore.user?.prenom?.[0] || authStore.user?.email?.[0] || 'U').toUpperCase()
)

const today = computed(() =>
  new Date().toLocaleDateString('fr-MA', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
)

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

function handleResize() {
  if (window.innerWidth >= 1024) {
    // Keep user preference on desktop
  } else {
    sidebarOpen.value = false
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>
