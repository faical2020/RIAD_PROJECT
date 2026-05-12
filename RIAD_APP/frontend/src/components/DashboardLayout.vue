<template>
  <div class="min-h-screen flex bg-riad-50">
    <!-- Desktop Sidebar -->
    <aside :class="[
      'fixed inset-y-0 left-0 z-50 flex flex-col bg-riad-950 transition-all duration-300 ease-out hidden lg:flex',
      sidebarOpen ? 'w-64' : 'w-20'
    ]">
      <!-- Logo -->
      <div class="flex items-center h-16 px-4 border-b border-riad-800/50" :class="sidebarOpen ? 'gap-3' : 'justify-center'">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-gold-500/30 to-gold-600/20 border border-gold-500/30 flex items-center justify-center flex-shrink-0 shadow-lg shadow-gold-500/10">
          <span class="text-lg">🌙</span>
        </div>
        <div v-show="sidebarOpen" class="overflow-hidden transition-all duration-300">
          <p class="font-display text-gold-300 text-sm font-semibold whitespace-nowrap">Riad Manager</p>
          <p class="text-riad-500 text-[10px] whitespace-nowrap font-arabic">نظام الإدارة</p>
        </div>
      </div>

      <!-- Nav -->
      <nav class="flex-1 p-3 space-y-1 overflow-y-auto">
        <RouterLink v-for="item in navItems" :key="item.to" :to="item.to"
          :class="['sidebar-link', { 'active': $route.name === item.name }]"
          :title="!sidebarOpen ? item.label : ''">
          <span class="text-xl flex-shrink-0" v-html="item.icon"></span>
          <span v-show="sidebarOpen" class="whitespace-nowrap text-sm">{{ item.label }}</span>
        </RouterLink>
      </nav>

      <!-- User -->
      <div class="p-3 border-t border-riad-800/50">
        <div :class="['flex items-center', sidebarOpen ? 'gap-3 px-2' : 'justify-center']">
          <div class="w-9 h-9 rounded-full bg-gradient-to-br from-gold-500/40 to-gold-600/20 border-2 border-gold-500/50 flex items-center justify-center flex-shrink-0 text-gold-300 font-bold text-sm shadow-lg shadow-gold-500/10">
            {{ userInitial }}
          </div>
          <div v-show="sidebarOpen" class="overflow-hidden min-w-0">
            <p class="text-riad-200 text-sm font-semibold truncate">{{ authStore.user?.prenom }} {{ authStore.user?.nom }}</p>
            <p class="text-riad-500 text-xs capitalize">{{ authStore.role }}</p>
          </div>
        </div>
        <button @click="handleLogout"
          class="sidebar-link w-full mt-2 text-riad-500 hover:text-red-400 hover:bg-red-500/10" :class="{ 'justify-center': !sidebarOpen }">
          <span class="text-xl flex-shrink-0">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/></svg>
          </span>
          <span v-show="sidebarOpen" class="text-sm">Déconnexion</span>
        </button>
      </div>
    </aside>

    <!-- Mobile Bottom Nav -->
    <nav class="fixed bottom-0 inset-x-0 z-50 lg:hidden bg-white/90 backdrop-blur-xl border-t border-riad-200 safe-area-bottom">
      <div class="flex items-center justify-around px-2 py-1">
        <RouterLink v-for="item in mobileNav" :key="item.to" :to="item.to"
          :class="['nav-link', { 'active': $route.name === item.name }]">
          <span class="nav-icon" v-html="item.icon"></span>
          <span>{{ item.label }}</span>
        </RouterLink>
      </div>
    </nav>

    <!-- Main area -->
    <div :class="['flex-1 flex flex-col min-w-0 transition-all duration-300', sidebarOpen ? 'lg:ml-64' : 'lg:ml-20']">
      <!-- Top bar -->
      <header class="sticky top-0 z-40 bg-white/80 backdrop-blur-xl border-b border-riad-200">
        <div class="flex items-center gap-3 px-4 sm:px-6 h-16">
          <button @click="sidebarOpen = !sidebarOpen" class="p-2 rounded-xl hover:bg-riad-100 transition-colors hidden lg:flex">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-riad-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"/></svg>
          </button>
          <div class="min-w-0">
            <h1 class="font-semibold text-riad-900 text-lg sm:text-xl truncate">{{ currentTitle }}</h1>
            <p class="text-riad-400 text-xs sm:text-sm">{{ today }}</p>
          </div>
          <div class="ml-auto flex items-center gap-2">
            <span class="hidden sm:inline-flex badge bg-gold-50 text-gold-700 border border-gold-200 capitalize">{{ authStore.role }}</span>
            <div class="flex lg:hidden w-8 h-8 rounded-full bg-gradient-to-br from-gold-500/40 to-gold-600/20 border-2 border-gold-500/50 flex items-center justify-center text-gold-600 font-bold text-xs">
              {{ userInitial }}
            </div>
          </div>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 p-4 sm:p-6 lg:p-8 overflow-auto max-w-7xl w-full mx-auto pb-24 lg:pb-8">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const sidebarOpen = ref(true)

const allNavItems = [
  { to: '/app', name: 'Dashboard', icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"/></svg>', label: 'Accueil', roles: null },
  { to: '/app/calendrier', name: 'Calendar', icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5"/></svg>', label: 'Planning', roles: ['manager', 'receptionniste'] },
  { to: '/app/cleaning', name: 'Cleaning', icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 00-5.78 1.128 2.25 2.25 0 01-2.4 2.245 4.5 4.5 0 008.4-2.245c0-.399-.078-.78-.22-1.128zm0 0a15.998 15.998 0 003.388-1.62m-5.043-.025a15.994 15.994 0 011.622-3.395m3.42 3.42a15.995 15.995 0 004.764-4.648l3.876-5.814a1.151 1.151 0 00-1.597-1.597L14.146 6.32a15.996 15.996 0 00-4.649 4.763m3.42 3.42a6.776 6.776 0 00-3.42-3.42"/></svg>', label: 'Ménage', roles: ['manager', 'receptionniste'] },
  { to: '/app/chambres', name: 'Chambres', icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"/></svg>', label: 'Chambres', roles: null },
  { to: '/app/reservations', name: 'Reservations', icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/></svg>', label: 'Réservations', roles: ['manager', 'receptionniste'] },
  { to: '/app/nouvelle-reservation', name: 'NouvelleReservation', icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>', label: 'Réserver', roles: null },
  { to: '/app/services', name: 'Services', icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9.53 16.122a3 3 0 00-5.78 1.128 2.25 2.25 0 01-2.4 2.245 4.5 4.5 0 008.4-2.245c0-.399-.078-.78-.22-1.128zm0 0a15.998 15.998 0 003.388-1.62m-5.043-.025a15.994 15.994 0 011.622-3.395m3.42 3.42a15.995 15.995 0 004.764-4.648l3.876-5.814a1.151 1.151 0 00-1.597-1.597L14.146 6.32a15.996 15.996 0 00-4.649 4.763m3.42 3.42a6.776 6.776 0 00-3.42-3.42"/></svg>', label: 'Services', roles: ['manager', 'receptionniste'] },
  { to: '/app/profil', name: 'Profil', icon: '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z"/></svg>', label: 'Profil', roles: null },
]

const navItems = computed(() =>
  allNavItems.filter(i => !i.roles || i.roles.includes(authStore.role))
)

const mobileNav = computed(() =>
  navItems.value.filter(i => i.name !== 'Profil').slice(0, 5)
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
  router.push('/')
}

function handleResize() {
  if (window.innerWidth < 1024) sidebarOpen.value = false
  else sidebarOpen.value = true
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>
