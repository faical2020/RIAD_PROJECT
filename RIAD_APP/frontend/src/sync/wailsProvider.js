import { Events } from '@wailsio/runtime'
import { useRiadStore } from '../stores/riad'

export const wailsProvider = {
    init() {
        const riadStore = useRiadStore()
        console.log('[Sync] ✅ Desktop Mode: Initializing Wails Events...');

        Events.On('sync:rooms', (data) => {
            console.log('[Sync] Wails Event: sync:rooms', data);
            riadStore.fetchChambres();
        });

        Events.On('sync:reservations', (data) => {
            console.log('[Sync] Wails Event: sync:reservations', data);
            riadStore.fetchReservations();
        });
    },
    destroy() {
        // Wails events are global to the app lifetime usually, 
        // but we could implement Off here if needed.
    }
}
