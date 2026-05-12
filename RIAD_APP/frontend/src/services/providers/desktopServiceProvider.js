import { GetLocalRooms, CreateLocalReservation, GetLocalReservations, SetToken, UpdateCleaningStatus } from '../../../bindings/RIAD_APP/riadservice';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1';

export const desktopServiceProvider = {
    async login(credentials) {
        try {
            const response = await fetch(`${API_BASE_URL}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(credentials)
            });
            const data = await response.json();
            if (!response.ok) throw new Error(data.error || 'Erreur lors de la connexion');

            localStorage.setItem('token', data.access_token || data.token);
            localStorage.setItem('refresh_token', data.refresh_token || '');
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
        try {
            const localRooms = await GetLocalRooms();
            return JSON.parse(JSON.stringify(localRooms || []));
        } catch (e) {
            console.warn("Local rooms fetch failed", e);
            return [];
        }
    },
    async getReservations() {
        try {
            return await GetLocalReservations();
        } catch (e) {
            console.warn("Local reservations fetch failed", e);
            return [];
        }
    },
    async getMyReservations() {
        try { return await GetLocalReservations() }
        catch { return [] }
    },
    async createReservation(reservationData) {
        const { client_id, chambre_id, date_debut, date_fin, montant } = reservationData;

        try {
            const localId = await CreateLocalReservation(client_id, chambre_id, date_debut, date_fin, montant);
            return { id: localId, synced: false };
        } catch (e) {
            console.error("Failed to save reservation locally", e);
            throw new Error("Impossible d'enregistrer la réservation");
        }
    },
    async updateReservation(reservationData) {
        const { id, chambre_id, date_debut, date_fin, montant, statut } = reservationData;
        return { id, synced: false };
    },
    async updateCleaningStatus(roomId, status) {
        try {
            return await UpdateCleaningStatus(roomId, status);
        } catch (e) {
            throw new Error('Erreur lors de la mise à jour du ménage: ' + e.message);
        }
    },
    async setToken(token) {
        try {
            await SetToken(token);
        } catch (e) {
            console.error("Failed to set token in Go service", e);
            throw e;
        }
    }
}
