import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './pages/routes'
import './main.css'
import secret from './config/secret'

const app = createApp(App)
app.use(createPinia())
app.use(router)

// The agentation-vue dev toolbar boots only in development so its code and
// styles never reach the production bundle. `import.meta.env.DEV` is replaced
// by Vite during transform, which removes the dynamic imports before any chunk
// is generated; the APP_ENV check keeps it off in non-dev modes.
if (import.meta.env.DEV && secret.APP_ENV === 'development') {
  const [{ AgentationVuePlugin }, _styles] = await Promise.all([
    import('agentation-vue'),
    import('agentation-vue/style.css'),
  ])
  app.use(AgentationVuePlugin)
}

app.mount('#app')
