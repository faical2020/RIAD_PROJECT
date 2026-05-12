import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './style.css'
import { riadService } from './services/serviceBridge'
import { webSyncManager } from './sync/webSyncManager'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)

const syncTokenWithGo = async () => {
    const raw = localStorage.getItem('token');
    const token = raw && raw !== 'undefined' && raw !== 'null' ? raw : null;
    if (token) {
        try {
            await riadService.setToken(token);
            console.log(`[Sync] Token synced successfully with service`);
        } catch (e) {
            console.error(`[Sync] Failed to sync token:`, e);
        }
    }
};

syncTokenWithGo();

// Start Web-specific background sync for drafts
const target = import.meta.env.VITE_APP_TARGET;
if (target === 'web') {
    webSyncManager.startAutoSync();
}

app.use(router)
app.mount('#app')
