import { useAuthStore } from '../stores/auth'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
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
        path: '/',
        component: () => import('../views/Dashboard.vue'),
        meta: { requiresAuth: true },
        children: [
            {
                path: '',
                redirect: '/chambres'
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
            }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    const auth = useAuthStore()

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        next('/login')
    } else if (to.meta.guest && auth.isAuthenticated) {
        next('/')
    } else {
        next()
    }
})

export default router
