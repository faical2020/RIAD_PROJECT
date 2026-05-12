import { Events } from '@wailsio/runtime'
import { useRiadStore } from '../stores/riad'

function debounce(fn, delay) {
    let timer = null
    return function (...args) {
        clearTimeout(timer)
        timer = setTimeout(() => fn(...args), delay)
    }
}

export const wailsProvider = {
    init() {
        const riadStore = useRiadStore()
        console.log('[Sync] ✅ Desktop Mode: Initializing Wails Events...');

        const debouncedRooms = debounce(() => riadStore.fetchChambres(), 500)
        const debouncedReservations = debounce(() => riadStore.fetchReservations(), 500)

        Events.On('sync:rooms', (data) => {
            console.log('[Sync] Wails Event: sync:rooms', data);
            debouncedRooms();
        });

        Events.On('sync:reservations', (data) => {
            console.log('[Sync] Wails Event: sync:reservations', data);
            debouncedReservations();
        });
    },
    destroy() {
        // Wails events are global to the app lifetime usually, 
        // but we could implement Off here if needed.
    }
}
