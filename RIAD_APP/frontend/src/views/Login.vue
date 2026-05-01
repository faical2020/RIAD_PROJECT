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
  <div class="min-h-screen bg-riad-950 ornament flex items-center justify-center p-4 sm:p-6 lg:p-8">
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-24 -right-24 w-96 h-96 rounded-full bg-gold-500/5 blur-3xl animate-float"></div>
      <div class="absolute -bottom-24 -left-24 w-96 h-96 rounded-full bg-riad-700/20 blur-3xl animate-float" style="animation-delay: -3s"></div>
    </div>

    <div class="relative w-full max-w-sm sm:max-w-md lg:max-w-lg">
      <div class="text-center mb-6 sm:mb-8">
        <div class="inline-flex items-center justify-center w-14 h-14 sm:w-16 sm:h-16 rounded-2xl bg-gradient-to-br from-gold-500/30 to-gold-600/20 border border-gold-500/30 mb-4 shadow-lg shadow-gold-500/10">
          <span class="text-2xl sm:text-3xl">🌙</span>
        </div>
        <h1 class="font-display text-2xl sm:text-3xl text-gold-300 mb-1">Riad Manager</h1>
        <p class="font-arabic text-riad-400 text-sm">نظام إدارة الرياض</p>
      </div>

      <div class="glass-card border-riad-700/50 p-6 sm:p-8">
        <h2 class="font-display text-xl sm:text-2xl text-white mb-6">Connexion</h2>

        <div v-if="authStore.error" class="mb-4 p-3 bg-terracotta-600/20 border border-terracotta-400/40 rounded-xl text-terracotta-400 text-sm backdrop-blur-sm">
          {{ authStore.error }}
        </div>

        <div class="space-y-4 sm:space-y-5">
          <div class="group">
            <label class="block text-riad-300 text-sm font-bold mb-1.5">Email</label>
            <input v-model="form.email" type="email" placeholder="user@riad.ma"
              class="input-lux group-hover:border-riad-600/70" />
          </div>
          <div class="group">
            <label class="block text-riad-300 text-sm font-bold mb-1.5">Mot de passe</label>
            <input v-model="form.password" type="password" placeholder="••••••••"
              class="input-lux group-hover:border-riad-600/70" />
          </div>
        </div>

        <button @click="handleLogin" :disabled="authStore.loading"
          class="mt-6 w-full bg-gradient-to-r from-gold-500 to-gold-400 hover:from-gold-400 hover:to-gold-300 disabled:opacity-50 text-riad-950 font-body font-bold py-3 sm:py-3.5 rounded-xl transition-all duration-300 flex items-center justify-center gap-2 shadow-lg shadow-gold-500/20 hover:shadow-xl hover:shadow-gold-500/30 active:scale-[0.98]">
          <span v-if="authStore.loading" class="flex items-center gap-2">
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-riad-950"></div>
            Connexion...
          </span>
          <span v-else>Se connecter</span>
        </button>

        <p class="mt-4 sm:mt-6 text-center text-riad-400 text-sm">
          Pas encore de compte ?
          <RouterLink to="/register" class="text-gold-400 hover:text-gold-300 font-bold transition-colors inline-block hover:translate-x-1">S'inscrire →</RouterLink>
        </p>
      </div>

      <div class="mt-4 p-4 bg-riad-900/40 backdrop-blur-sm border border-riad-700/30 rounded-xl text-xs text-riad-400">
        <p class="font-bold text-riad-300 mb-2">🔑 Comptes démo</p>
        <div class="space-y-1">
          <button @click="fillDemo('manager@riad.ma','123456')" class="block w-full text-left hover:text-gold-400 transition-all duration-200 hover:translate-x-1">
            Manager: manager@riad.ma / 123456
          </button>
          <button @click="fillDemo('client@riad.ma','123456')" class="block w-full text-left hover:text-gold-400 transition-all duration-200 hover:translate-x-1">
            Client: client@riad.ma / 123456
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
