import { createRouter, createWebHashHistory } from 'vue-router'
import { authGuard, useAuth0 } from '@auth0/auth0-vue'

import HomeView from '../views/HomeView.vue'
import GettingStartedView from '../views/GettingStartedView.vue'
import AuthView from '../views/AuthView.vue'
import Library from '../components/Library.vue'

import { userStore } from '../stores/user'

function isAuthAndHasLocal(to) {
  const { isAuthenticated } = useAuth0()
  const uStore = userStore()

  if (isAuthenticated && uStore.hasJWT && uStore.hasNkey) {
    return { name: 'library' }
  }

  return true
}

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/:code?', name: 'home', component: HomeView, beforeEnter: [isAuthAndHasLocal] },
    { path: '/getting-started', name: 'gettingstarted', component: GettingStartedView, beforeEnter: [isAuthAndHasLocal] },
    { path: '/library', name: 'library', component: Library, beforeEnter: authGuard }
  ]
})

export default router
