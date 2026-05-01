import { defineStore } from 'pinia'
import { riadService } from '../services/riadService'

export const useRiadStore = defineStore('riad', {
    state: () => ({
        chambres: [],
        reservations: [],
        loading: false,
    }),

    getters: {
        rooms: (state) => state.chambres,
        libres: (state) => state.chambres.filter(r => r.statut === 'libre'),
        occupees: (state) => state.chambres.filter(r => r.statut === 'occupee' || r.statut === 'occupe'),
        availableRooms: (state) => state.chambres.filter(r => r.statut === 'libre'),
        occupiedRooms: (state) => state.chambres.filter(r => r.statut === 'occupee' || r.statut === 'occupe'),
        getRoomById: (state) => (id) => state.chambres.find(r => r.id === id),
    },

    actions: {
        async fetchChambres() {
            this.loading = true
            try {
                this.chambres = await riadService.getRooms()
            } finally {
                this.loading = false
            }
        },

        async fetchRooms() {
            return this.fetchChambres()
        },

        async fetchReservations() {
            try {
                this.reservations = await riadService.getReservations()
            } catch (e) {
                console.error('Failed to fetch reservations', e)
            }
        },

        async fetchStats() {
            // Stats sont calculées via les getters
        },
    },
})
