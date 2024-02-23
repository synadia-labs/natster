<template>
  <div v-if="isLoading">Loading...</div>
  <div v-else>
    <AuthView v-if="isAuthenticated">
      <router-view />
    </AuthView>
    <router-view v-else />
  </div>
</template>

<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue'
import AuthView from './views/AuthView.vue'

const { isAuthenticated, isLoading } = useAuth0()
</script>
