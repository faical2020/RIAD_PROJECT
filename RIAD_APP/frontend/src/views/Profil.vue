<script setup>
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
</script>

<template>
  <div class="card p-6 sm:p-8 max-w-2xl mx-auto">
    <h2 class="font-display text-xl sm:text-2xl text-riad-900 mb-4 sm:mb-6">Mon Profil</h2>

    <div v-if="auth.user" class="space-y-4 sm:space-y-6">
      <div class="flex flex-col sm:flex-row items-center gap-4 sm:gap-6 p-4 sm:p-6 bg-gradient-to-br from-riad-50 to-gold-50 rounded-xl sm:rounded-2xl">
        <div class="w-14 h-14 sm:w-20 sm:h-20 rounded-full bg-gradient-to-br from-gold-500/30 to-gold-600/20 border-2 border-gold-500/50 flex items-center justify-center text-gold-300 font-bold text-2xl sm:text-3xl shadow-lg shadow-gold-500/10">
          {{ (auth.user.prenom?.[0] || auth.user.email?.[0] || 'U').toUpperCase() }}
        </div>
        <div class="text-center sm:text-left">
          <h3 class="font-display text-xl sm:text-2xl text-riad-900">{{ auth.user.prenom }} {{ auth.user.nom }}</h3>
          <p class="text-riad-500 text-sm sm:text-base">{{ auth.user.email }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
        <div class="p-4 sm:p-5 bg-white rounded-xl border border-riad-100 hover:shadow-lg transition-all duration-300">
          <p class="text-riad-500 text-xs font-bold uppercase tracking-wider">Rôle</p>
          <p class="text-riad-900 font-bold capitalize mt-1 sm:mt-2 text-lg">{{ auth.role }}</p>
        </div>
        <div class="p-4 sm:p-5 bg-white rounded-xl border border-riad-100 hover:shadow-lg transition-all duration-300">
          <p class="text-riad-500 text-xs font-bold uppercase tracking-wider">Téléphone</p>
          <p class="text-riad-900 font-bold mt-1 sm:mt-2 text-lg">{{ auth.user.telephone || 'Non renseigné' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
