const API_BASE_URL = 'http://localhost:8081/api/v1';

const hasWailsBindings = () => typeof window !== 'undefined' && window.go && window.go.main && window.go.main.RiadService;

export const riadService = {
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
        let rooms = [];

        if (hasWailsBindings()) {
            try {
                rooms = await window.go.main.RiadService.GetLocalRooms();
            } catch (e) {
                console.warn("Local rooms fetch failed", e);
            }
        }

        if (navigator.onLine) {
            const token = localStorage.getItem('token');
            try {
                const response = await fetch(`${API_BASE_URL}/chambres`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                if (response.ok) {
                    const serverRooms = await response.json();

                    if (hasWailsBindings()) {
                        for (const room of serverRooms) {
                            await window.go.main.RiadService.UpdateLocalRoom(
                                room.id, room.numero, room.type, room.prix, room.description, room.equipements, room.statut
                            );
                        }
                    }

                    rooms = serverRooms;
                }
            } catch (e) {
                console.warn("Cloud sync failed, using local data", e);
            }
        }

        return rooms;
    },

    async getReservations() {
        if (navigator.onLine) {
            const token = localStorage.getItem('token');
            try {
                const response = await fetch(`${API_BASE_URL}/reservations`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                if (response.ok) {
                    return await response.json();
                }
            } catch (e) {
                console.warn("Failed to fetch reservations from server", e);
            }
        }
        return [];
    },

    async createReservation(reservationData) {
        const { client_id, chambre_id, date_debut, date_fin, montant } = reservationData;

        if (navigator.onLine) {
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
            } catch (e) {
                console.warn("Failed to create reservation", e);
            }
        }

        return { synced: false };
    },
};
