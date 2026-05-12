<script setup>
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <h2 class="text-xl font-bold text-riad-900">Mon Profil</h2>

    <div v-if="auth.user" class="space-y-4">
      <!-- Avatar card -->
      <div class="card p-6 sm:p-8 flex flex-col sm:flex-row items-center gap-6">
        <div class="w-20 h-20 rounded-full bg-gradient-to-br from-gold-500/30 to-gold-600/20 border-2 border-gold-500/50 flex items-center justify-center text-gold-600 font-bold text-3xl shadow-lg shadow-gold-500/10">
          {{ (auth.user.prenom?.[0] || auth.user.email?.[0] || 'U').toUpperCase() }}
        </div>
        <div class="text-center sm:text-left">
          <h3 class="text-xl font-bold text-riad-900">{{ auth.user.prenom }} {{ auth.user.nom }}</h3>
          <p class="text-riad-500 text-sm">{{ auth.user.email }}</p>
          <span class="mt-1 inline-flex badge bg-gold-50 text-gold-700 border border-gold-200 capitalize text-xs">{{ auth.role }}</span>
        </div>
      </div>

      <!-- Info cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="card p-5">
          <p class="text-xs font-semibold text-riad-400 uppercase tracking-wider mb-1">Rôle</p>
          <p class="text-lg font-bold text-riad-900 capitalize">{{ auth.role }}</p>
        </div>
        <div class="card p-5">
          <p class="text-xs font-semibold text-riad-400 uppercase tracking-wider mb-1">Téléphone</p>
          <p class="text-lg font-bold text-riad-900">{{ auth.user.telephone || 'Non renseigné' }}</p>
        </div>
        <div class="card p-5">
          <p class="text-xs font-semibold text-riad-400 uppercase tracking-wider mb-1">Email</p>
          <p class="text-lg font-bold text-riad-900 break-all">{{ auth.user.email }}</p>
        </div>
        <div class="card p-5">
          <p class="text-xs font-semibold text-riad-400 uppercase tracking-wider mb-1">Nom complet</p>
          <p class="text-lg font-bold text-riad-900">{{ auth.user.prenom }} {{ auth.user.nom }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
