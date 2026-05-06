import { localCacheService } from '../localCacheService';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1';

export const webServiceProvider = {
    async login(credentials) {
        try {
            const response = await fetch(`${API_BASE_URL}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(credentials)
            });
            const data = await response.json();
            if (!response.ok) throw new Error(data.error || 'Erreur lors de la connexion');

            localStorage.setItem('token', data.token);
            localStorage.setItem('role', data.role);
            return data;
        } catch (e) {
            throw new Error(e.message || 'Erreur de connexion');
        }
    },
    async register(userData) {
        try {
            const response = await fetch(`${API_BASE_URL}/auth/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(userData)
            });
            const data = await response.json();
            if (!response.ok) throw new Error(data.error || 'Erreur lors de l\'inscription');
            return data;
        } catch (e) {
            throw new Error(e.message || 'Erreur d\'inscription');
        }
    },
    async getRooms() {
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_BASE_URL}/chambres`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                await localCacheService.setAll('rooms', data);
                return data;
            }
            throw new Error('Failed to fetch rooms');
        } catch (e) {
            console.warn("Web: Server unavailable, loading rooms from cache...", e);
            const cached = await localCacheService.getAll('rooms');
            if (cached && cached.length > 0) return cached;
            throw e;
        }
    },
    async getReservations() {
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_BASE_URL}/reservations`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                await localCacheService.setAll('reservations', data);
                return data;
            }
            throw new Error('Failed to fetch reservations');
        } catch (e) {
            console.warn("Web: Server unavailable, loading reservations from cache...", e);
            const cached = await localCacheService.getAll('reservations');
            if (cached && cached.length > 0) return cached;
            throw e;
        }
    },
    async createReservation(reservationData) {
        const { client_id, chambre_id, date_debut, date_fin, montant } = reservationData;
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_BASE_URL}/reservations`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    user_id: client_id,
                    chambre_id: chambre_id,
                    date_debut: date_debut,
                    date_fin: date_fin,
                    montant: montant
                })
            });

            if (response.ok) {
                const serverRes = await response.json();
                return { id: serverRes.id, synced: true };
            }
            throw new Error('Server rejected reservation');
        } catch (e) {
            console.warn("Web: Server unavailable, saving as draft in IndexedDB...", e);
            const draftId = await localCacheService.saveDraft({
                user_id: client_id,
                chambre_id: chambre_id,
                date_debut: date_debut,
                date_fin: date_fin,
                montant: montant,
                status: 'draft'
            });
            return { id: draftId, synced: false };
        }
    },
    async updateReservation(reservationData) {
        const { id, client_id, chambre_id, date_debut, date_fin, montant, statut } = reservationData;
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_BASE_URL}/reservations/${id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    user_id: client_id,
                    chambre_id: chambre_id,
                    date_debut: date_debut,
                    date_fin: date_fin,
                    montant: montant,
                    statut: statut
                })
            });

            if (response.ok) {
                const updatedRes = await response.json();
                await localCacheService.setAll('reservations', await localCacheService.getAll('reservations')); // Simplified cache update
                return { id: updatedRes.id, synced: true };
            }
            throw new Error('Server rejected update');
        } catch (e) {
            console.warn("Web: Server unavailable, updating cache only...", e);
            // Update local cache as draft/unsynced
            const resList = await localCacheService.getAll('reservations');
            const idx = resList.findIndex(r => r.id === id);
            if (idx !== -1) {
                resList[idx] = { ...resList[idx], ...reservationData, synced: false };
                await localCacheService.setAll('reservations', resList);
            }
            return { id, synced: false };
        }
    },
    async updateCleaningStatus(roomId, status) {
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_BASE_URL}/chambres/${roomId}/cleaning`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ cleaning_status: status })
            });

            if (response.ok) {
                const data = await response.json();
                // Refresh cache
                const rooms = await this.getRooms();
                await localCacheService.setAll('rooms', rooms);
                return data;
            }
            throw new Error('Server rejected cleaning status update');
        } catch (e) {
            console.warn("Web: Server unavailable, updating cache only...", e);
            const rooms = await localCacheService.getAll('rooms');
            const idx = rooms.findIndex(r => r.id === roomId);
            if (idx !== -1) {
                rooms[idx].cleaning_status = status;
                await localCacheService.setAll('rooms', rooms);
            }
            return { id: roomId, cleaning_status: status };
        }
    },
    async setToken(token) {
        return Promise.resolve();
    }
}
