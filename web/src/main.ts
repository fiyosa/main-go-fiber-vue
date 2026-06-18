import './main.css'
import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import routes from './pages/routes.ts'

createApp(App)
  .use(createPinia())
  .use(VueQueryPlugin)
  .use(routes)
  .mount('#app')
