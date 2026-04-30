const API_BASE_URL = 'http://localhost:8081/api/v1';

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
            throw e.message || e;
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
            throw e.message || e;
        }
    },

    async getRooms() {
        let rooms = [];
        try {
            rooms = await window.go.main.RiadService.GetLocalRooms();
        } catch (e) {
            console.warn("Local rooms fetch failed", e);
        }

        if (navigator.onLine) {
            try {
                const response = await fetch(`${API_BASE_URL}/chambres`, {
                    headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
                });
                if (response.ok) {
                    const serverRooms = await response.json();
                    for (const room of serverRooms) {
                        await window.go.main.RiadService.UpdateLocalRoom(
                            room.id, room.numero, room.type, room.prix, room.description, room.equipements, room.statut
                        );
                    }
                    rooms = serverRooms;
                }
            } catch (e) {
                console.warn("Cloud sync failed, using local data", e);
            }
        }
        return rooms;
    },

    async createReservation(reservationData) {
        const { userId, roomId, start, end, amount } = reservationData;
        
        const localId = await window.go.main.RiadService.CreateLocalReservation(userId, roomId, start, end, amount);

        if (navigator.onLine) {
            try {
                const response = await fetch(`${API_BASE_URL}/reservations`, {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token')}` 
                    },
                    body: JSON.stringify({
                        user_id: userId,
                        chambre_id: roomId,
                        date_debut: start,
                        date_fin: end,
                        montant: amount
                    })
                });

                if (response.ok) {
                    const serverRes = await response.json();
                    await window.go.main.RiadService.MarkAsSynced('reservations', localId);
                    return { id: serverRes.id, synced: true };
                }
            } catch (e) {
                console.warn("Offline mode: reservation pending sync", e);
            }
        }

        return { id: localId, synced: false };
    },

    async syncAll() {
        if (!navigator.onLine) return;

        const tables = ['reservations', 'rooms'];
        for (const table of tables) {
            const unsynced = await window.go.main.RiadService.GetUnsynced(table);
            for (const item of unsynced) {
                try {
                    if (table === 'reservations') {
                        await fetch(`${API_BASE_URL}/reservations`, {
                            method: 'POST',
                            headers: { 
                                'Content-Type': 'application/json',
                                'Authorization': `Bearer ${localStorage.getItem('token')}` 
                            },
                            body: JSON.stringify({
                                user_id: item.user_id,
                                chambre_id: item.chambre_id,
                                date_debut: item.date_debut,
                                date_fin: item.date_fin,
                                montant: item.montant
                            })
                        });
                    }
                    await window.go.main.RiadService.MarkAsSynced(table, item.id);
                } catch (e) {
                    console.error(`Failed to sync ${item.id}`, e);
                }
            }
        }
    },

    async getDashboardStats() {
        try {
            const rooms = await window.go.main.RiadService.GetLocalRooms();
            const unsyncedRes = await window.go.main.RiadService.GetUnsynced('reservations');
            const unsyncedRooms = await window.go.main.RiadService.GetUnsynced('rooms');

            const totalRooms = rooms.length;
            const occupiedRooms = rooms.filter(r => r.statut === 'occupee').length;

            return {
                totalRooms,
                occupiedRooms,
                availableRooms: totalRooms - occupiedRooms,
                pendingSync: unsyncedRes.length + unsyncedRooms.length
            };
        } catch (e) {
            console.error("Stats fetch failed", e);
            return { totalRooms: 0, occupiedRooms: 0, availableRooms: 0, pendingSync: 0 };
        }
    }
};

