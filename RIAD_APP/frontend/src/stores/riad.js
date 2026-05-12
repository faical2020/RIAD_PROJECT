import { defineStore } from 'pinia'
import { riadService } from '../services/serviceBridge'

let _lastRoomsFetch = 0
let _lastReservationsFetch = 0

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
        
        // Helper for Calendar
        reservationsForRoomAndDate: (state) => (roomId, dateStr) => {
            return state.reservations.find(res => 
                res.chambre_id === roomId && 
                dateStr >= res.date_debut && 
                dateStr <= res.date_fin
            );
        },
    },

    actions: {
        async fetchChambres() {
            const now = Date.now()
            if (now - _lastRoomsFetch < 1000) return
            _lastRoomsFetch = now
            console.log('[Store] fetchChambres called');
            this.loading = true
            try {
                const rooms = await riadService.getRooms()
                console.log(`[Store] fetchChambres received ${rooms.length} rooms`, rooms);
                this.chambres = rooms
            } catch (e) {
                console.error('[Store] fetchChambres error:', e);
            } finally {
                this.loading = false
            }
        },

        async fetchRooms() {
            return this.fetchChambres()
        },

        async fetchReservations() {
            const now = Date.now()
            if (now - _lastReservationsFetch < 1000) return
            _lastReservationsFetch = now
            try {
                this.reservations = await riadService.getReservations()
            } catch (e) {
                console.error('Failed to fetch reservations', e)
            }
        },

        async fetchMyReservations() {
            try {
                this.reservations = await riadService.getMyReservations()
            } catch (e) {
                console.error('Failed to fetch my reservations', e)
            }
        },

        async createReservation(reservationData) {
            try {
                const result = await riadService.createReservation(reservationData)
                await this.fetchReservations()
                return result
            } catch (e) {
                console.error('[Store] createReservation error:', e);
                throw e;
            }
        },

        async updateReservation(reservationData) {
            try {
                const result = await riadService.updateReservation(reservationData)
                await this.fetchReservations()
                return result
            } catch (e) {
                console.error('[Store] updateReservation error:', e);
                throw e;
            }
        },

        async updateCleaningStatus(roomId, status) {
            try {
                await riadService.updateCleaningStatus(roomId, status)
                await this.fetchChambres()
            } catch (e) {
                console.error('[Store] updateCleaningStatus error:', e);
                throw e;
            }
        },
    },
})
