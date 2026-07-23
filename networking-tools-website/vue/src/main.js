import { createApp } from 'vue'
import { createVuestic } from 'vuestic-ui'
import 'vuestic-ui/css'
import './assets/theme.css'
import App from './App.vue'
import router from './router'

createApp(App)
  .use(router)
  .use(createVuestic({
    config: {
      colors: {
        variables: {
          primary: '#49d7c4',
          background: '#0a1420',
          backgroundElement: '#101f30',
          backgroundSecondary: '#16283c',
          textPrimary: '#e4f7fa',
          textInverted: '#06201c',
        }
      }
    }
  }))
  .mount('#app')