import './scss/main.scss'

import { createApp } from 'vue'
import App from './App.vue'
import i18n from './i18n'
import { createVuestic } from 'vuestic-ui'
import { createGtm } from '@gtm-support/vue-gtm'
import stores from './stores'
import router from './router'

const app = createApp(App)

app.use(stores)
app.use(i18n)
app.use(router)
app.use(
  createVuestic({
    config: {
      colors: {
        variables: {
          primary: '#ffffff',
          secondary: '#ffbb22',
          backgroundPrimary: '#000000',
          backgroundSecondary: '#333333',
          backgroundElement: '#333333',
          textPrimary: '#ffffff',
          textSecondary: '#ffbb22',
          textInverted: '#000000',
          backgroundBorder: '#cccccc',

          // GSG colours
          gsgYellow: "#ffbb22",
          gsgRed: "#e42222",

          // custom colors
          discordBlurple: '#5865F2',
          mastodonLight: '#6364FF',
          mastodonDark: '#563ACC',
          twitchPurple: '#874af6',
          blueskyBlue: '#3983f7',
          youtubeRed: '#ea333e',
          steamBlue: '#0a183a',
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
