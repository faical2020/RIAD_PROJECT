import { defineStore } from 'pinia'
import { riadService } from '../services/riadService'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: JSON.parse(localStorage.getItem('user') || 'null'),
        token: localStorage.getItem('token') || null,
        role: localStorage.getItem('role') || null,
        loading: false,
        error: null,
    }),

    getters: {
        isAuthenticated: (state) => !!state.token,
        isAdmin: (state) => state.role === 'manager',
        isStaff: (state) => ['manager', 'receptionniste'].includes(state.role),
        isReceptionniste: (state) => ['manager', 'receptionniste'].includes(state.role),
    },

    actions: {
        async login(email, password) {
            this.loading = true
            this.error = null
            try {
                const result = await riadService.login({ email, password })
                this.token = result.token
                this.role = result.role
                this.user = result.user || { email, role: result.role }
                if (result.user) {
                    this.user = result.user
                }

                localStorage.setItem('token', result.token)
                localStorage.setItem('role', result.role)
                localStorage.setItem('user', JSON.stringify(this.user))
                return true
            } catch (e) {
                this.error = e.message || 'Erreur de connexion'
                return false
            } finally {
                this.loading = false
            }
        },

        async register(userData) {
            this.loading = true
            this.error = null
            try {
                await riadService.register(userData)
                return true
            } catch (e) {
                this.error = e.message || 'Erreur d\'inscription'
                return false
            } finally {
                this.loading = false
            }
        },

        logout() {
            this.token = null
            this.role = null
            this.user = null
            this.error = null

            localStorage.removeItem('token')
            localStorage.removeItem('role')
            localStorage.removeItem('user')
        },
    },
})
