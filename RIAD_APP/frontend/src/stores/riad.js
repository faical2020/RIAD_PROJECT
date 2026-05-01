import { defineStore } from 'pinia'
import { riadService } from '../services/riadService'

export const useRiadStore = defineStore('riad', {
    state: () => ({
        rooms: [],
        reservations: [],
        stats: {
            totalRooms: 0,
            occupiedRooms: 0,
            availableRooms: 0,
            pendingSync: 0
        },
        loading: false,
    }),

    getters: {
        availableRooms: (state) => state.rooms.filter(r => r.statut === 'libre'),
        occupiedRooms: (state) => state.rooms.filter(r => r.statut === 'occupee'),
        getRoomById: (state) => (id) => state.rooms.find(r => r.id === id),
    },

    actions: {
        async fetchRooms() {
            this.loading = true
            try {
                this.rooms = await riadService.getRooms()
            } finally {
                this.loading = false
            }
        },

        async fetchStats() {
            this.stats = await riadService.getDashboardStats()

            if (this.stats.totalRooms === 0 && this.rooms.length > 0) {
                this.stats = {
                    totalRooms: this.rooms.length,
                    occupiedRooms: this.rooms.filter(r => r.statut === 'occupee').length,
                    availableRooms: this.rooms.filter(r => r.statut === 'libre').length,
                    pendingSync: this.rooms.filter(r => !r.synced).length
                }
            }
        },

        async createReservation(reservationData) {
            const result = await riadService.createReservation(reservationData)
            await this.fetchRooms()
            await this.fetchStats()
            return result
        },

        async syncAll() {
            await riadService.syncAll()
            await this.fetchRooms()
            await this.fetchStats()
        },

        async initData() {
            await this.fetchRooms()
            await this.fetchStats()
            this.startStatsPolling()
        },

        startStatsPolling() {
            setInterval(async () => {
                await this.fetchRooms()
                await this.fetchStats()
            }, 30000)
        },
    },
})
