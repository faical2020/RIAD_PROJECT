import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './style.css'
import { SetToken } from '../bindings/RIAD_APP/riadservice'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)

// Sync token with Go service on startup
const syncTokenWithGo = async () => {
    const token = localStorage.getItem('token');
    if (token) {
        try {
            await SetToken(token);
            console.log(`[Sync] Token synced successfully with Go service`);
        } catch (e) {
            console.error(`[Sync] Failed to sync token:`, e);
        }
    }
};

syncTokenWithGo();

app.use(router)
app.mount('#app')
