import './scss/main.scss'

import { createApp } from 'vue'
import App from './App.vue'
import i18n from './i18n'
import { createVuestic } from 'vuestic-ui'
import { createGtm } from '@gtm-support/vue-gtm'
import axiosPlugin from './plugins/axios'
import stores from './stores'
import router from './router'
import vuesticGlobalConfig from './services/vuestic-ui/global-config'
import { createPinia } from 'pinia'

const pinia = createPinia()
const app = createApp(App)

app.use(stores)
app.use(router)
app.use(i18n)
app.use(pinia)
app.use(axiosPlugin, {
  baseUrl: '/api',
})
app.use(
  createVuestic({
    config: {
      colors: {
        variables: {
          primary: '#ffffff',
          secondary: '#ffbb22',
          backgroundPrimary: '#000000',
          backgroundSecondary: '#333333',
          textPrimary: '#ffffff',
          textSecondary: '#ffbb22',
          textInverted: '#000000',
          backgroundBorder: '#cccccc',

          // custom colors
          discordBlurple: '#5865F2',
        },
      },
    },
  }),
)

if (import.meta.env.VITE_APP_GTM_ENABLED) {
  app.use(
    createGtm({
      id: import.meta.env.VITE_APP_GTM_KEY,
      debug: false,
      vueRouter: router,
    }),
  )
}

app.mount('#app')
