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
import { watch } from 'vue'
import { useAuth0 } from '@auth0/auth0-vue'
import { storeToRefs } from 'pinia'
import { JSONCodec } from 'nats.ws'
import type { Catalog } from './types/types'

import AuthView from './views/AuthView.vue'
import { userStore } from './stores/user'
import { natsStore } from './stores/nats'

const { isAuthenticated, isLoading } = useAuth0()
const uStore = userStore()
const nStore = natsStore()

const { connection } = storeToRefs(nStore)

watch(connection, () => {
  if (nStore.connection !== null) {
    nStore.getShares(true)
    nStore.ping()
  }
})
</script>
