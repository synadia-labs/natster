import { createRouter, createWebHistory } from 'vue-router'
import { authGuard, useAuth0 } from '@auth0/auth0-vue'

import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import AuthView from '../views/AuthView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/login/:code?', name: 'login', component: LoginView },
    { path: '/library', name: 'library', component: AuthView, beforeEnter: authGuard },
    { path: '/shares', name: 'shares', component: AuthView, beforeEnter: authGuard },
    {
      path: '/callback',
      redirect: (to) => {
        console.log(to.path, to.query.code, to.query.state)
        return { path: '/library', query: { code: to.query.code, state: to.query.state } }
      }
    }
  ]
})

export default router
