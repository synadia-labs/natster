import { createRouter, createWebHashHistory } from 'vue-router'
import { authGuard, useAuth0 } from '@auth0/auth0-vue'

import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import AuthView from '../views/AuthView.vue'
import Library from '../components/Library.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/login/:code?', name: 'login', component: LoginView },
    { path: '/library', name: 'library', component: Library, beforeEnter: authGuard }
  ]
})

export default router
