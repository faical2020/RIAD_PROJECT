<script setup>
import { reactive, ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const success = ref(false)

const form = reactive({
  prenom: '', nom: '', email: '', password: '', telephone: '', role: 'client'
})

async function handleRegister() {
  success.value = false
  const ok = await authStore.register(form)
  if (ok) {
    success.value = true
    Object.assign(form, { prenom: '', nom: '', email: '', password: '', telephone: '', role: 'client' })
  }
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
        <p class="text-riad-400 text-lg leading-relaxed">Créez votre compte et rejoignez notre établissement</p>
        <p class="font-arabic text-riad-600 text-base mt-4">مرحباً بكم في فريقنا</p>
      </div>
    </div>

    <!-- Right Panel -->
    <div class="flex-1 flex items-center justify-center p-4 sm:p-8 lg:p-12 bg-gradient-to-br from-riad-50 to-white">
      <div class="w-full max-w-sm animate-in">
        <div class="text-center lg:hidden mb-8">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-gold-500/30 to-gold-600/20 border border-gold-500/30 mb-4 shadow-lg shadow-gold-500/10">
            <span class="text-3xl">🌙</span>
          </div>
          <h1 class="font-display text-2xl text-riad-900 mb-1">Riad Manager</h1>
        </div>

        <h2 class="text-2xl font-bold text-riad-900 mb-1">Créer un compte</h2>
        <p class="text-riad-400 text-sm mb-8">Inscrivez-vous pour commencer</p>

        <div v-if="authStore.error" class="mb-6 p-3 bg-red-50 border border-red-200 rounded-xl text-red-600 text-sm flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"/></svg>
          {{ authStore.error }}
        </div>
        <div v-if="success" class="mb-6 p-3 bg-green-50 border border-green-200 rounded-xl text-green-700 text-sm flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 flex-shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          Compte créé ! Vous pouvez vous connecter.
        </div>

        <form @submit.prevent="handleRegister" class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-semibold text-riad-700 mb-1.5">Prénom</label>
              <input v-model="form.prenom" type="text" placeholder="Mohammed" class="input" />
            </div>
            <div>
              <label class="block text-sm font-semibold text-riad-700 mb-1.5">Nom</label>
              <input v-model="form.nom" type="text" placeholder="Alaoui" class="input" />
            </div>
          </div>
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Email</label>
            <input v-model="form.email" type="email" placeholder="vous@exemple.com" class="input" />
          </div>
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Téléphone</label>
            <input v-model="form.telephone" type="tel" placeholder="0612345678" class="input" />
          </div>
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Mot de passe</label>
            <input v-model="form.password" type="password" placeholder="••••••••" class="input" />
          </div>
          <div>
            <label class="block text-sm font-semibold text-riad-700 mb-1.5">Rôle</label>
            <select v-model="form.role" class="select">
              <option value="client">Client</option>
              <option value="employe">Employé</option>
              <option value="receptionniste">Réceptionniste</option>
              <option value="manager">Manager</option>
            </select>
          </div>

          <button type="submit" :disabled="authStore.loading"
            class="btn-primary w-full py-3 mt-2">
            <span v-if="authStore.loading" class="flex items-center gap-2">
              <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/></svg>
              Création...
            </span>
            <span v-else>Créer le compte</span>
          </button>
        </form>

        <p class="mt-6 text-center text-sm text-riad-500">
          Déjà un compte ?
          <RouterLink to="/login" class="text-gold-600 hover:text-gold-700 font-semibold">Se connecter</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>
