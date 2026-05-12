import { defineStore } from 'pinia'
import { riadService } from '../services/serviceBridge'

function decodeToken(token) {
    try {
        return JSON.parse(atob(token.split('.')[1]))
    } catch { return null }
}

function validToken(t) {
    return t && t !== 'undefined' && t !== 'null' ? t : null
}

export const useAuthStore = defineStore('auth', {
    state: () => {
        const token = validToken(localStorage.getItem('token'))
        let user = JSON.parse(localStorage.getItem('user') || 'null')
        const role = localStorage.getItem('role') || null

        if (token && (!user || !user.id)) {
            const payload = decodeToken(token)
            if (payload && payload.user_id) {
                if (!user) user = { id: payload.user_id, role: payload.role }
                else user.id = payload.user_id
                localStorage.setItem('user', JSON.stringify(user))
            }
        }

        return { user, token, role, loading: false, error: null }
    },

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
                const token = result.access_token || result.token
                this.token = token
                this.role = result.role
                this.user = result.user || { email, role: result.role }
                if (result.user) {
                    this.user = result.user
                }

                localStorage.setItem('token', token)
                localStorage.setItem('role', result.role)
                localStorage.setItem('user', JSON.stringify(this.user))
                if (result.refresh_token) {
                    localStorage.setItem('refresh_token', result.refresh_token)
                }
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
            localStorage.removeItem('refresh_token')
        },

        async refreshToken() {
            const refreshToken = localStorage.getItem('refresh_token')
            if (!refreshToken) return false
            try {
                const res = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1'}/auth/refresh`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ refresh_token: refreshToken })
                })
                if (!res.ok) throw new Error('refresh failed')
                const data = await res.json()
                this.token = data.access_token
                localStorage.setItem('token', data.access_token)
                localStorage.setItem('refresh_token', data.refresh_token)
                return true
            } catch {
                this.logout()
                return false
            }
        },
    },
})
