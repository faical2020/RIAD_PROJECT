<template>
    <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import { Events } from '@wailsio/runtime'
import { useRiadStore } from './stores/riad'

const riadStore = useRiadStore()

const setupSync = () => {
    // Robust detection of Wails environment
    const isWails = !!(window.go || window.Wails || window._wails);

    if (isWails) {
        console.log('[Sync] ✅ Desktop Mode detected. Using Wails Events...');
        
        Events.On('sync:rooms', (data) => {
            console.log('[Sync] Wails Event: sync:rooms', data);
            riadStore.fetchChambres();
        });

        Events.On('sync:reservations', (data) => {
            console.log('[Sync] Wails Event: sync:reservations', data);
            riadStore.fetchReservations();
        });
    } else {
        console.log('[Sync] 🌐 Web Mode detected. Using SSE stream...');
        const token = localStorage.getItem('token');
        if (token) {
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
        }
    }
}

onMounted(() => {
    console.log('[Sync] Mounting App.vue...');
    
    // Give Wails a moment to inject window.go if it's not already there
    if (!(window.go || window.Wails || window._wails)) {
        console.log('[Sync] Wails not detected immediately, retrying in 100ms...');
        setTimeout(setupSync, 100);
    } else {
        setupSync();
    }
})
</script>


<style>
body {
    margin: 0;
    padding: 0;
    font-family: system-ui, -apple-system, sans-serif;
}
</style>
