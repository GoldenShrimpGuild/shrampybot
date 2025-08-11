import { createApp } from 'vue'
import App from './App.vue'
import i18n from './i18n'
import { createVuesticEssential } from 'vuestic-ui'
import stores from './stores'
import router from './router'

// CSS stuff
import 'vuestic-ui/styles/essential.css'
import 'vuestic-ui/styles/typography.css'
import './style/icon-fonts/vuestic-icons/vuestic-icons.css'
import './style/main.css'
import './style/shrampybot.css'

const app = createApp(App)

app.use(stores)
app.use(i18n)
app.use(router)
app.use(
  createVuesticEssential({
    components: {},
    plugins: {},
    config: {
      components: {
        VaButton: {
          size: "small",
          round: true,
          color: "#ffbb22",
          style: "font-weight: bold;"
        },
        VaSwitch: {
          size: "small",
          color: "#ffbb22",
          offColor: "#715411",
          style: "line-height: 12px;"
        }
      },
      breakpoint: {
        enabled: false,
        bodyClass: true,
        thresholds: {
          xs: 0,
          sm: 0,
          md: 0,
          lg: 0,
          xl: 0,
        },
      },
      colors: {
        variables: {
          primary: '#ffffff',
          secondary: '#ffbb22',
          backgroundPrimary: '#000000',
          backgroundSecondary: '#333333',
          backgroundElement: '#694e00ff',
          textPrimary: '#ffffff',
          textSecondary: '#ffbb22',
          textInverted: '#000000',
          backgroundBorder: '#cccccc',

          // GSG colours
          gsgYellow: '#ffbb22',
          gsgRed: '#e42222',

          // custom colors
          gsgDarkYellow: '#715411',
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

app.mount('#app')
