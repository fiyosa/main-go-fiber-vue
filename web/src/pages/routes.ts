import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

import GuestLayout from './guest/GuestLayout.vue'
import AdminLayout from './admin/AdminLayout.vue'
import AdminNotFound from './admin/AdminNotFound.vue'
import { authGuard } from './authGuard'

const routes: RouteRecordRaw[] = [
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'admin',
        component: () => import('./admin/AdminHome.vue'),
      },
      {
        path: 'auth/register',
        name: 'admin-register',
        component: () => import('./admin/AuthRegister.vue'),
      },
      {
        path: 'auth/change-password',
        name: 'admin-change-password',
        component: () => import('./admin/AuthChangePassword.vue'),
      },
      {
        path: 'work/cron',
        name: 'admin-cron',
        component: () => import('./admin/CronManager.vue'),
      },
      {
        path: 'work/cron/create',
        name: 'admin-cron-create',
        component: () => import('./admin/CronForm.vue'),
      },
      {
        path: 'work/cron/:id/edit',
        name: 'admin-cron-edit',
        component: () => import('./admin/CronForm.vue'),
      },
    ],
  },
  {
    // Top-level admin catch-all: unknown /admin/* paths render a standalone
    // 404 page without the AdminLayout header and sidebar.
    path: '/admin/:pathMatch(.*)*',
    name: 'admin-not-found',
    component: AdminNotFound,
    meta: { requiresAuth: true },
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
