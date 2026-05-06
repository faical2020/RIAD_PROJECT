const target = import.meta.env.VITE_APP_TARGET;

let provider = null;

const getProvider = async () => {
    if (provider) return provider;
    if (target === 'desktop') {
        const { wailsProvider } = await import('./wailsProvider');
        provider = wailsProvider;
    } else {
        const { sseProvider } = await import('./sseProvider');
        provider = sseProvider;
    }
    return provider;
};

export const syncProvider = {
    async init() { return (await getProvider()).init(); },
    async destroy() { return (await getProvider()).destroy(); },
};

export default syncProvider;
