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
  <div class="min-h-screen bg-riad-950 ornament flex items-center justify-center p-4 sm:p-6 lg:p-8">
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-24 -right-24 w-96 h-96 rounded-full bg-gold-500/5 blur-3xl animate-float"></div>
    </div>

    <div class="relative w-full max-w-sm sm:max-w-md">
      <div class="text-center mb-6 sm:mb-8">
        <div class="inline-flex items-center justify-center w-14 h-14 sm:w-16 sm:h-16 rounded-2xl bg-gradient-to-br from-gold-500/30 to-gold-600/20 border border-gold-500/30 mb-4 shadow-lg shadow-gold-500/10">
          <span class="text-2xl sm:text-3xl">🌙</span>
        </div>
        <h1 class="font-display text-2xl sm:text-3xl text-gold-300">Riad Manager</h1>
      </div>

      <div class="glass-card border-riad-700/50 p-6 sm:p-8">
        <h2 class="font-display text-xl sm:text-2xl text-white mb-6">Créer un compte</h2>

        <div v-if="authStore.error" class="mb-4 p-3 bg-terracotta-600/20 border border-terracotta-400/40 rounded-xl text-terracotta-400 text-sm backdrop-blur-sm">
          {{ authStore.error }}
        </div>
        <div v-if="success" class="mb-4 p-3 bg-green-900/30 border border-green-500/40 rounded-xl text-green-400 text-sm">
          ✅ Compte créé ! Vous pouvez vous connecter.
        </div>

        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-3">
            <div class="group">
              <label class="block text-riad-300 text-sm font-bold mb-1.5">Prénom</label>
              <input v-model="form.prenom" type="text" placeholder="Mohammed"
                class="input-lux group-hover:border-riad-600/70" />
            </div>
            <div class="group">
              <label class="block text-riad-300 text-sm font-bold mb-1.5">Nom</label>
              <input v-model="form.nom" type="text" placeholder="Alaoui"
                class="input-lux group-hover:border-riad-600/70" />
            </div>
          </div>
          <div class="group">
            <label class="block text-riad-300 text-sm font-bold mb-1.5">Email</label>
            <input v-model="form.email" type="email" placeholder="user@riad.ma"
              class="input-lux group-hover:border-riad-600/70" />
          </div>
          <div class="group">
            <label class="block text-riad-300 text-sm font-bold mb-1.5">Téléphone</label>
            <input v-model="form.telephone" type="tel" placeholder="0612345678"
              class="input-lux group-hover:border-riad-600/70" />
          </div>
          <div class="group">
            <label class="block text-riad-300 text-sm font-bold mb-1.5">Mot de passe</label>
            <input v-model="form.password" type="password" placeholder="••••••••"
              class="input-lux group-hover:border-riad-600/70" />
          </div>
          <div class="group">
            <label class="block text-riad-300 text-sm font-bold mb-1.5">Rôle</label>
            <select v-model="form.role"
              class="input-lux group-hover:border-riad-600/70">
              <option value="client">Client</option>
              <option value="employe">Employé</option>
              <option value="receptionniste">Réceptionniste</option>
              <option value="manager">Manager</option>
            </select>
          </div>
        </div>

        <button @click="handleRegister" :disabled="authStore.loading"
          class="mt-6 w-full bg-gradient-to-r from-gold-500 to-gold-400 hover:from-gold-400 hover:to-gold-300 disabled:opacity-50 text-riad-950 font-body font-bold py-3 sm:py-3.5 rounded-xl transition-all duration-300 shadow-lg shadow-gold-500/20 hover:shadow-xl hover:shadow-gold-500/30 active:scale-[0.98]">
          <span v-if="authStore.loading" class="flex items-center justify-center gap-2">
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-riad-950"></div>
            Création...
          </span>
          <span v-else>Créer le compte</span>
        </button>

        <p class="mt-4 sm:mt-6 text-center text-riad-400 text-sm">
          Déjà un compte ?
          <RouterLink to="/login" class="text-gold-400 hover:text-gold-300 font-bold transition-colors inline-block hover:translate-x-1">Se connecter →</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>
