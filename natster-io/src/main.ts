import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createAuth0 } from '@auth0/auth0-vue'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(
  createAuth0({
    domain: 'login.natster.io',
    clientId: 'veI5fgi7qKMaYc4SRs1CTgpNL2RfgRFK',
    authorizationParams: {
      redirect_uri: window.location.origin + '/#/library'
    }
  })
)

app.mount('#app')
