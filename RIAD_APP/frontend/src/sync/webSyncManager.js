import { localCacheService } from '../services/localCacheService'

export const webSyncManager = {
    async syncDrafts() {
        const drafts = await localCacheService.getDrafts()
        if (drafts.length === 0) return

        console.log(`[WebSync] Found ${drafts.length} drafts to synchronize...`)
        const token = localStorage.getItem('token')
        if (!token) return

        const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1'
        
        for (const draft of drafts) {
            try {
                const response = await fetch(`${API_BASE_URL}/reservations`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    },
                    body: JSON.stringify({
                        user_id: draft.user_id,
                        chambre_id: draft.chambre_id,
                        date_debut: draft.date_debut,
                        date_fin: draft.date_fin,
                        montant: draft.montant
                    })
                })

                if (response.ok) {
                    console.log(`[WebSync] Draft ${draft.id} synced successfully`)
                }
            } catch (e) {
                console.error(`[WebSync] Failed to sync draft ${draft.id}:`, e)
                break // Stop if server is still down
            }
        }
        
        // Clear drafts after attempt (or implement a per-item delete)
        await localCacheService.clearDrafts()
    },

    startAutoSync() {
        setInterval(() => this.syncDrafts(), 30000) // Check every 30 seconds
    }
}
