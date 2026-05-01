<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRiadStore } from '../stores/riad'
import { useRoute } from 'vue-router'
import StatsBar from '../components/StatsBar.vue'

const auth = useAuthStore()
const riad = useRiadStore()
const route = useRoute()

onMounted(() => {
    riad.initData()
})
</script>

<template>
    <div class="min-h-screen bg-gray-50 text-gray-900 font-sans">
        <nav class="bg-white border-b border-gray-200 py-3 px-6 flex justify-between items-center sticky top-0 z-20 shadow-sm">
            <div class="flex items-center gap-4">
                <h1 class="text-xl font-bold text-emerald-700 tracking-tight">RIAD <span class="text-gray-400 font-light">Management</span></h1>
                <div class="hidden md:flex items-center gap-2 px-3 py-1 bg-emerald-50 text-emerald-600 rounded-full text-xs font-medium border border-emerald-100">
                    <div class="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
                    Connecté
                </div>
            </div>

            <div class="flex gap-2 items-center">
                <router-link to="/chambres"
                             :class="['px-4 py-2 text-sm font-medium rounded-lg transition-all', route.path === '/chambres' ? 'bg-emerald-600 text-white' : 'text-gray-500 hover:bg-gray-100']">
                    Chambres
                </router-link>
                <router-link to="/reservations"
                             :class="['px-4 py-2 text-sm font-medium rounded-lg transition-all', route.path === '/reservations' ? 'bg-emerald-600 text-white' : 'text-gray-500 hover:bg-gray-100']">
                    Réservations
                </router-link>
                <div class="w-px h-6 bg-gray-200 mx-2"></div>
                <button @click="auth.logout()" class="p-2 text-gray-400 hover:text-red-500 transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                </button>
            </div>
        </nav>

        <main class="container mx-auto py-8 px-4">
            <StatsBar :stats="riad.stats" />
            <router-view />
        </main>

        <footer class="bg-white border-t border-gray-200 py-6 text-center text-gray-400 text-xs">
            &copy; 2026 RIAD_APP Hybrid System. Local-first, Cloud-synced.
        </footer>
    </div>
</template>
