import { useRiadStore } from '../stores/riad'

export const sseProvider = {
    init() {
        const riadStore = useRiadStore()
        console.log('[Sync] 🌐 Web Mode: Initializing SSE stream...');
        
        const token = localStorage.getItem('token');
        if (!token) {
            console.warn('[Sync] No token found in localStorage, SSE not started.');
            return;
        }

        const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1';
        const eventSource = new EventSource(`${baseUrl}/sync/events?token=${token}`);
        
        eventSource.onopen = () => console.log('[Sync] SSE Connection opened successfully!');
        eventSource.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.type === 'ROOM_UPDATED') riadStore.fetchChambres();
                else if (data.type === 'RESERVATION_UPDATED') riadStore.fetchReservations();
            } catch (e) {
                console.error('[Sync] SSE Parse error:', e);
            }
        };
        eventSource.onerror = (err) => console.error('[Sync] SSE Error:', err);

        return eventSource; // Return for potential cleanup
    },
    destroy(eventSource) {
        if (eventSource) {
            eventSource.close();
            console.log('[Sync] SSE Connection closed.');
        }
    }
}
