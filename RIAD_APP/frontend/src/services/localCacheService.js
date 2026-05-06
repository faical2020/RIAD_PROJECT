const DB_NAME = 'riad_web_cache';
const DB_VERSION = 1;
const STORES = {
    rooms: 'rooms',
    reservations: 'reservations',
    drafts: 'drafts'
};

export const localCacheService = {
    async openDB() {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open(DB_NAME, DB_VERSION);

            request.onupgradeneeded = (event) => {
                const db = event.target.result;
                if (!db.objectStoreNames.contains(STORES.rooms)) {
                    db.createObjectStore(STORES.rooms, { keyPath: 'id' });
                }
                if (!db.objectStoreNames.contains(STORES.reservations)) {
                    db.createObjectStore(STORES.reservations, { keyPath: 'id' });
                }
                if (!db.objectStoreNames.contains(STORES.drafts)) {
                    db.createObjectStore(STORES.drafts, { keyPath: 'id', autoIncrement: true });
                }
            };

            request.onsuccess = () => resolve(request.result);
            request.onerror = () => reject(request.error);
        });
    },

    async setAll(storeName, data) {
        const db = await this.openDB();
        const tx = db.transaction(storeName, 'readwrite');
        const store = tx.objectStore(storeName);
        
        // Clear old data
        await store.clear();
        
        // Add new data
        data.forEach(item => store.put(item));
        return new Promise((resolve) => {
            tx.oncomplete = () => resolve();
        });
    },

    async getAll(storeName) {
        const db = await this.openDB();
        return new Promise((resolve, reject) => {
            const tx = db.transaction(storeName, 'readonly');
            const store = tx.objectStore(storeName);
            const request = store.getAll();
            request.onsuccess = () => resolve(request.result);
            request.onerror = () => reject(request.error);
        });
    },

    async saveDraft(data) {
        const db = await this.openDB();
        const tx = db.transaction(STORES.drafts, 'readwrite');
        const store = tx.objectStore(STORES.drafts);
        return new Promise((resolve) => {
            const request = store.add({ ...data, created_at: Date.now() });
            request.onsuccess = () => resolve(request.result);
        });
    },

    async getDrafts() {
        const db = await this.openDB();
        return new Promise((resolve, reject) => {
            const tx = db.transaction(STORES.drafts, 'readonly');
            const store = tx.objectStore(STORES.drafts);
            const request = store.getAll();
            request.onsuccess = () => resolve(request.result);
            request.onerror = () => reject(request.error);
        });
    },

    async clearDrafts() {
        const db = await this.openDB();
        const tx = db.transaction(STORES.drafts, 'readwrite');
        const store = tx.objectStore(STORES.drafts);
        store.clear();
    }
};
