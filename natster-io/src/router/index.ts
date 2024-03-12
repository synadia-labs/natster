import { createRouter, createWebHashHistory } from 'vue-router'
import { authGuard, useAuth0 } from '@auth0/auth0-vue'

import HomeView from '../views/HomeView.vue'
import GettingStartedView from '../views/GettingStartedView.vue'
import AuthView from '../views/AuthView.vue'
import Library from '../components/Library.vue'

import { userStore } from '../stores/user'

function isNotLoggedIn(to, from, next) {
  const uStore = userStore()

  if (!uStore.hasJWT && !uStore.hasNkey) {
    next()
    return
  }

  next('library')
}

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/:code?', name: 'home', component: HomeView, beforeEnter: isNotLoggedIn },
    { path: '/getting-started', name: 'gettingstarted', component: GettingStartedView, beforeEnter: isNotLoggedIn },
    { path: '/library', name: 'library', component: Library, beforeEnter: authGuard }
  ]
})

export default router
