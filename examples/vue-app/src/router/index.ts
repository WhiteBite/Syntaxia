import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView
        },
        {
            path: '/users',
            name: 'users',
            component: () => import('../views/UsersView.vue')
        },
        {
            path: '/users/:id',
            name: 'user-detail',
            component: () => import('../views/UserDetailView.vue')
        },
        {
            path: '/about',
            name: 'about',
            component: () => import('../views/AboutView.vue')
        }
    ]
})

export default router
