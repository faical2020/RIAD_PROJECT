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
                return await response.json();
            }
            throw new Error('Failed to fetch rooms');
        } catch (e) {
            console.error("Web: Cloud fetch failed", e);
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
                return await response.json();
            } else if (response.status === 403) {
                console.warn("Server: Access forbidden to reservations list.");
                return [];
            }
            throw new Error('Failed to fetch reservations');
        } catch (e) {
            console.error("Web: Cloud fetch failed", e);
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
            throw new Error("Aucune connexion serveur et stockage local indisponible (Mode Web)");
        }
    },
    async setToken(token) {
        // No-op in web mode: token is sent via Authorization headers in each request
        return Promise.resolve();
    }
}
