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
import { Catalog } from './types/types'

import AuthView from './views/AuthView.vue'
import { userStore } from './stores/user'
import { natsStore } from './stores/nats'

const { isAuthenticated, isLoading } = useAuth0()
const uStore = userStore()
const nStore = natsStore()

const { connection } = storeToRefs(nStore)

watch(connection, () => {
  if (nStore.connection !== null) {
    nStore.connection
      .request('natster.global.my.shares', '', { timeout: 5000 })
      .then((m) => {
        const catalogs = JSONCodec().decode(m.data).data
        catalogs.forEach((c, i) => {
          if (c.to_account === uStore.getAccount) {
            const catalog: Catalog = {
              selected: false,
              to: c.to_account,
              from: c.from_account,
              name: c.catalog,
              online: c.catalog_online,
              files: []
            }
            uStore.catalogs.push(catalog)

            if (i == 0) {
              uStore.setCatalogSelected(catalog)
            }
          }
        })
        nStore.ping()
      })
      .catch((err) => {
        console.log(`problem with request: ${err.message}`)
      })
  }
})
</script>
