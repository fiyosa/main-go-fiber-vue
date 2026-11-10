import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

import GuestLayout from './guest/GuestLayout.vue'

const routes: RouteRecordRaw[] = [
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
        path: 'about',
        name: 'about',
        component: () => import('./guest/About.vue'),
      },
      {
        path: 'tech',
        name: 'tech',
        component: () => import('./guest/Tech.vue'),
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

export default router
