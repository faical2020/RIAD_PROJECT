const target = import.meta.env.VITE_APP_TARGET;

let service = null;

const getService = async () => {
    if (service) return service;
    if (target === 'desktop') {
        const { desktopServiceProvider } = await import('./providers/desktopServiceProvider');
        service = desktopServiceProvider;
    } else {
        const { webServiceProvider } = await import('./providers/webServiceProvider');
        service = webServiceProvider;
    }
    return service;
};

export const riadService = {
    async login(credentials) { return (await getService()).login(credentials); },
    async register(userData) { return (await getService()).register(userData); },
    async getRooms() { return (await getService()).getRooms(); },
    async getReservations() { return (await getService()).getReservations(); },
    async createReservation(reservationData) { return (await getService()).createReservation(reservationData); },
    async setToken(token) { return (await getService()).setToken(token); },
};

export default riadService;
