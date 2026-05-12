import { useRiadStore } from '../stores/riad'

function debounce(fn, delay) {
    let timer = null
    return function (...args) {
        clearTimeout(timer)
        timer = setTimeout(() => fn(...args), delay)
    }
}

export const sseProvider = {
    init() {
        const riadStore = useRiadStore()
        console.log('[Sync] 🌐 Web Mode: Initializing SSE stream...');
        
        const token = localStorage.getItem('token');
        if (!token) {
            console.warn('[Sync] No token found in localStorage, SSE not started.');
            return;
        }

        const debouncedRooms = debounce(() => riadStore.fetchChambres(), 500)
        const debouncedReservations = debounce(() => riadStore.fetchReservations(), 500)

        const baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8081/api/v1';
        const eventSource = new EventSource(`${baseUrl}/sync/events?token=${token}`);
        
        eventSource.onopen = () => console.log('[Sync] SSE Connection opened successfully!');
        eventSource.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.type === 'ROOM_UPDATED') debouncedRooms();
                else if (data.type === 'RESERVATION_UPDATED') debouncedReservations();
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
