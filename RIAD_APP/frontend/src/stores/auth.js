import { defineStore } from 'pinia'
import { riadService } from '../services/riadService'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: JSON.parse(localStorage.getItem('user') || 'null'),
        token: localStorage.getItem('token') || null,
        role: localStorage.getItem('role') || null,
    }),

    getters: {
        isAuthenticated: (state) => !!state.token,
        isAdmin: (state) => state.role === 'manager',
        isStaff: (state) => ['manager', 'receptionniste'].includes(state.role),
    },

    actions: {
        async login(email, password) {
            const result = await riadService.login({ email, password })
            this.token = result.token
            this.role = result.role
            this.user = { email, role: result.role }

            localStorage.setItem('token', result.token)
            localStorage.setItem('role', result.role)
            localStorage.setItem('user', JSON.stringify(this.user))
        },

        async register(userData) {
            await riadService.register(userData)
        },

        logout() {
            this.token = null
            this.role = null
            this.user = null

            localStorage.removeItem('token')
            localStorage.removeItem('role')
            localStorage.removeItem('user')
        },
    },
})
