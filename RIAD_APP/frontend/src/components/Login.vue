<script setup>
import { ref } from 'vue'
import { riadService } from '../services/serviceBridge'

const form = ref({
  email: '',
  password: '',
})
const error = ref('')
const emit = defineEmits(['login-success'])

const handleLogin = async () => {
  error.value = ''
  try {
    const result = await riadService.login(form.value)
    emit('login-success', result)
  } catch (e) {
    error.value = e
  }
}
</script>

<template>
  <div class="p-6 max-w-md mx-auto">
    <div class="bg-white rounded-2xl shadow-lg p-8 border border-gray-100">
      <h2 class="text-2xl font-bold mb-6 text-gray-800 text-center">Connexion</h2>
      
      <form @submit.prevent="handleLogin" class="space-y-4">
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
        <button type="submit" 
                class="w-full py-3 bg-emerald-600 hover:bg-emerald-700 text-white font-semibold rounded-lg transition-colors shadow-sm mt-4">
          Se connecter
        </button>
      </form>

      <div v-if="error" class="mt-6 p-3 rounded-lg text-center text-sm font-medium bg-red-50 text-red-600">
        {{ error }}
      </div>
    </div>
  </div>
</template>
