import { useAuthStore } from '../stores/auth'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        name: 'Landing',
        component: () => import('../views/Landing.vue'),
        meta: { guest: true }
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue'),
        meta: { guest: true }
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('../views/Register.vue'),
        meta: { guest: true }
    },
    {
        path: '/app',
        component: () => import('../components/DashboardLayout.vue'),
        meta: { requiresAuth: true },
        children: [
            {
                path: '',
                name: 'Dashboard',
                component: () => import('../views/Dashboard.vue')
            },
            {
                path: 'calendrier',
                name: 'Calendar',
                component: () => import('../views/Calendar.vue')
            },
            {
                path: 'cleaning',
                name: 'Cleaning',
                component: () => import('../views/Cleaning.vue')
            },
            {
                path: 'chambres',
                name: 'Chambres',
                component: () => import('../views/Chambres.vue')
            },
            {
                path: 'reservations',
                name: 'Reservations',
                component: () => import('../views/Reservations.vue')
            },
            {
                path: 'nouvelle-reservation',
                name: 'NouvelleReservation',
                component: () => import('../views/NouvelleReservation.vue')
            },
            {
                path: 'services',
                name: 'Services',
                component: () => import('../views/Services.vue')
            },
            {
                path: 'consommations/:reservationId',
                name: 'Consommations',
                component: () => import('../views/Consommations.vue')
            },
            {
                path: 'facture/:reservationId',
                name: 'Facture',
                component: () => import('../views/Facture.vue')
            },
            {
                path: 'profil',
                name: 'Profil',
                component: () => import('../views/Profil.vue')
            }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from) => {
    const auth = useAuthStore()

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        return '/login'
    }
    if (to.meta.guest && auth.isAuthenticated && to.name !== 'Landing') {
        return '/app'
    }
    if (to.name === 'Landing' && auth.isAuthenticated) {
        return '/app'
    }
    return true
})

export default router
