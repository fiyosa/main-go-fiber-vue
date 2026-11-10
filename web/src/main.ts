import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { AgentationVuePlugin } from 'agentation-vue'
import 'agentation-vue/style.css'

import App from './App.vue'
import router from './pages/routes'
import './main.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)

if (import.meta.env.DEV) app.use(AgentationVuePlugin)

app.mount('#app')
