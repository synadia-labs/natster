import { createRouter, createWebHashHistory } from 'vue-router'
import { authGuard, useAuth0 } from '@auth0/auth0-vue'

import HomeView from '../views/HomeView.vue'
import AuthView from '../views/AuthView.vue'
import Library from '../components/Library.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/:code?', name: 'home', component: HomeView },
    { path: '/library', name: 'library', component: Library, beforeEnter: authGuard }
  ]
})

export default router
