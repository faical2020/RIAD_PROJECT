<script setup>
import { ref } from 'vue'
import { riadService } from '../services/serviceBridge'

const form = ref({
  email: '',
  password: '',
  nom: '',
  prenom: '',
  telephone: '',
})
const error = ref('')
const success = ref(false)

const handleRegister = async () => {
  error.value = ''
  try {
    await riadService.register(form.value)
    success.value = true
  } catch (e) {
    error.value = e
  }
}
</script>

<template>
  <div class="p-6 max-w-md mx-auto">
    <div class="bg-white rounded-2xl shadow-lg p-8 border border-gray-100">
      <h2 class="text-2xl font-bold mb-6 text-gray-800 text-center">Inscription</h2>
      
      <div v-if="success" class="p-4 bg-green-50 text-green-600 rounded-lg text-center mb-6 font-medium">
        Compte créé avec succès ! Vous pouvez vous connecter.
      </div>

      <form v-else @submit.prevent="handleRegister" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium text-gray-700">Nom</label>
            <input v-model="form.nom" type="text" required 
                   class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium text-gray-700">Prénom</label>
            <input v-model="form.prenom" type="text" required 
                   class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-gray-700">Email</label>
          <input v-model="form.email" type="email" required 
                 class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-gray-700">Mot de passe</label>
          <input v-model="form.password" type="password" required 
                 class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-gray-700">Téléphone</label>
          <input v-model="form.telephone" type="text" 
                 class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition-all" />
        </div>

        <button type="submit" 
                class="w-full py-3 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold rounded-lg transition-colors shadow-sm mt-4">
          S'inscrire
        </button>
      </form>

      <div v-if="error" class="mt-6 p-3 rounded-lg text-center text-sm font-medium bg-red-50 text-red-600">
        {{ error }}
      </div>
    </div>
  </div>
</template>
