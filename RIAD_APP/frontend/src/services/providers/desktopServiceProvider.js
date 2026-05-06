import { GetLocalRooms, UpdateLocalRoom, CreateLocalReservation, GetLocalReservations, SetToken } from '../../../bindings/RIAD_APP/riadservice';

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

        try {
            const localRooms = await GetLocalRooms();
            rooms = JSON.parse(JSON.stringify(localRooms || []));
        } catch (e) {
            console.warn("Local rooms fetch failed", e);
        }

        if (navigator.onLine) {
            const token = localStorage.getItem('token');
            try {
                const response = await fetch(`${API_BASE_URL}/chambres`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                if (response.ok) {
                    const serverRooms = await response.json();

                    for (const room of serverRooms) {
                        try {
                            await UpdateLocalRoom(
                                room.id, room.numero, room.type, room.prix, room.description, room.equipements, room.statut
                            );
                        } catch (e) {
                            console.warn(`Failed to update local room ${room.id}`, e);
                        }
                    }
                    rooms = JSON.parse(JSON.stringify(serverRooms));
                }
            } catch (e) {
                console.warn("Cloud sync failed, using local data", e);
            }
        }
        return rooms;
    },
    async getReservations() {
        let reservations = [];

        try {
            reservations = await GetLocalReservations();
        } catch (e) {
            console.warn("Local reservations fetch failed", e);
        }

        if (navigator.onLine) {
            const token = localStorage.getItem('token');
            try {
                const response = await fetch(`${API_BASE_URL}/reservations`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                if (response.ok) {
                    const serverRes = await response.json();
                    const allRes = [...reservations, ...serverRes];
                    const uniqueRes = Array.from(new Map(allRes.map(item => [item.id, item])).values());
                    return uniqueRes;
                } else if (response.status === 403) {
                    console.warn("Server: Access forbidden to reservations list. Showing local data only.");
                }
            } catch (e) {
                console.warn("Failed to fetch reservations from server", e);
            }
        }
        return reservations;
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
                console.warn("Server unavailable, saving locally...", e);
            }
        }

        try {
            const localId = await CreateLocalReservation(client_id, chambre_id, date_debut, date_fin, montant);
            return { id: localId, synced: false };
        } catch (e) {
            console.error("Failed to save reservation locally", e);
            throw new Error("Impossible d'enregistrer la réservation");
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
