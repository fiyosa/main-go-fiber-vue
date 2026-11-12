import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

import GuestLayout from './guest/GuestLayout.vue'
import { authGuard } from './authGuard'

const routes: RouteRecordRaw[] = [
  {
    path: '/admin',
    component: () => import('./admin/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'admin',
        component: () => import('./admin/AdminHome.vue'),
      },
    ],
  },
  {
    path: '/',
    component: GuestLayout,
    children: [
      {
        path: '',
        name: 'home',
        component: () => import('./guest/Home.vue'),
      },
      {
        path: 'tech',
        name: 'tech',
        component: () => import('./guest/Tech.vue'),
      },
      {
        path: 'login',
        name: 'login',
        component: () => import('./guest/Login.vue'),
      },
      ...(import.meta.env.VITE_SHOW_SERVER === 'true'
        ? [
            {
              path: 'server',
              name: 'server',
              component: () => import('./guest/Server.vue'),
            },
          ]
        : []),
      {
        path: ':pathMatch(.*)*',
        name: 'not-found',
        component: () => import('./guest/404.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach(authGuard)

export default router
