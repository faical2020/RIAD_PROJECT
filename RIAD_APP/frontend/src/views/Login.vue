<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const form = reactive({ email: '', password: '' })

function fillDemo(email, password) {
  form.email = email
  form.password = password
}

async function handleLogin() {
  const ok = await authStore.login(form.email, form.password)
  if (ok) router.push('/')
}
</script>

<template>
  <div class="min-h-screen flex">
    <!-- Left Panel - Branding -->
    <div class="hidden lg:flex lg:w-1/2 bg-riad-950 relative overflow-hidden items-center justify-center p-12">
      <div class="absolute inset-0">
        <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-gold-500/5 rounded-full blur-3xl"></div>
        <div class="absolute bottom-1/4 right-1/4 w-64 h-64 bg-riad-700/20 rounded-full blur-3xl"></div>
      </div>
      <div class="relative text-center max-w-md">
        <div class="inline-flex items-center justify-center w-20 h-20 rounded-3xl bg-gradient-to-br from-gold-500/20 to-gold-600/10 border border-gold-500/30 mb-6 shadow-lg shadow-gold-500/10">
          <span class="text-4xl">🌙</span>
        </div>
        <h1 class="font-display text-4xl text-gold-300 mb-3">Riad Manager</h1>
        <p class="text-riad-400 text-lg leading-relaxed">Système de gestion hôtelière</p>
        <p class="font-arabic text-riad-600 text-base mt-4">نظام إدارة الرياض الفاخرة</p>
      </div>
    </div>

    <!-- Right Panel - Form -->
    <div class="flex-1 flex items-center justify-center p-4 sm:p-8 lg:p-12 bg-gradient-to-br from-riad-50 to-white">
      <div class="w-full max-w-sm animate-in">
        <div class="text-center lg:hidden mb-8">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-gold-500/30 to-gold-600/20 border border-gold-500/30 mb-4 shadow-lg shadow-gold-500/10">
            <span class="text-3xl">🌙</span>
          </div>
          <h1 class="font-display text-2xl text-riad-900 mb-1">Riad Manager</h1>
          <p class="text-riad-400 text-sm">Système de gestion hôtelière</p>
        </div>

        <h2 class="text-2xl font-bold text-riad-900 mb-1">Bienvenue</h2>
        <p class="text-riad-400 text-sm mb-8">Connectez-vous à votre compte</p>

        <div v-if="authStore.error" class="mb-6 p-3 bg-red-50 border border-red-200 rounded-xl text-red-600 text-sm flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"/></svg>
          {{ authStore.error }}
        </div>

        <form @submit.prevent="handleLogin" class="space-y-5">
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Email</label>
            <input v-model="form.email" type="email" placeholder="vous@exemple.com"
              class="input" autocomplete="email" />
          </div>
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Mot de passe</label>
            <input v-model="form.password" type="password" placeholder="••••••••"
              class="input" autocomplete="current-password" />
          </div>

          <button type="submit" :disabled="authStore.loading"
            class="btn-primary w-full py-3">
            <span v-if="authStore.loading" class="flex items-center gap-2">
              <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
              Connexion...
            </span>
            <span v-else>Se connecter</span>
          </button>
        </form>

        <p class="mt-6 text-center text-sm text-riad-500">
          Pas encore de compte ?
          <RouterLink to="/register" class="text-gold-600 hover:text-gold-700 font-semibold">S'inscrire</RouterLink>
        </p>

        <div class="mt-6 p-4 bg-riad-50 rounded-xl border border-riad-200">
          <p class="font-semibold text-riad-700 text-xs mb-2 uppercase tracking-wider">🔑 Accès rapide</p>
          <div class="space-y-1.5">
            <button @click="fillDemo('manager@riad.ma','123456')" type="button"
              class="block w-full text-left text-xs text-riad-500 hover:text-gold-600 transition-colors py-1 px-2 rounded-lg hover:bg-white">
              manager@riad.ma / 123456
            </button>
            <button @click="fillDemo('client@riad.ma','123456')" type="button"
              class="block w-full text-left text-xs text-riad-500 hover:text-gold-600 transition-colors py-1 px-2 rounded-lg hover:bg-white">
              client@riad.ma / 123456
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
