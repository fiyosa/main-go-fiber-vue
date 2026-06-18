import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const path: readonly RouteRecordRaw[] = [
  { name: 'Home', path: '/', component: () => import('./guest/Home.vue') },
  { name: 'About', path: '/about', component: () => import('./guest/About.vue') },
  { name: 'Logger', path: '/logger', component: () => import('./guest/Logger.vue') },

  { path: '/:pathMatch(.*)*', component: import('./guest/404.vue') },
]

const routes = createRouter({ history: createWebHistory(), routes: path })

export default routes
