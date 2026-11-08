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
