<template>
  <VueSpinnerAudio v-if="isLoading" size="160" class="loading-spinner" />
  <div v-else>
    <VueSpinnerTail v-if="activeDownload" size="160" class="downloading-spinner" />
    <AuthView v-if="isAuthenticated">
      <router-view />
    </AuthView>
    <router-view v-else />
  </div>
</template>

<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue'
import { storeToRefs } from 'pinia'

import AuthView from './views/AuthView.vue'
import { catalogStore } from './stores/catalog'
const { activeDownload } = storeToRefs(catalogStore())

import { VueSpinnerAudio, VueSpinnerTail } from 'vue3-spinners'
const { isAuthenticated, isLoading } = useAuth0()
</script>

<style>
.loading-spinner {
  color: #45c320;
  margin: auto auto;
}

.downloading-spinner {
  position: fixed;
  top: 50%;
  left: 50%;
  margin-top: -50px;
  margin-left: -100px;
  color: #45c320;
  z-index: 50;
}
</style>
