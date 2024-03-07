<template>
  <VueSpinnerAudio 
    v-if="isLoading"
    size="160"
    class="loading-spinner"
  />
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
import { VueSpinnerAudio } from 'vue3-spinners'

const { isAuthenticated, isLoading } = useAuth0()
</script>

<style>
.loading-spinner {
  color: #45c320;
  margin: 10% auto;
}
</style>
