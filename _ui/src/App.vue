<template>
  <SideBar v-if="loggedIn">
    <router-view />
  </SideBar>
  <router-view v-else />
</template>

<script setup>
import { watch, onMounted, computed } from "vue";
import { storeToRefs } from "pinia";
import { JSONCodec } from "nats.ws";
import { useRoute } from "vue-router";
import router from "./router";

import SideBar from "./components/SideBar.vue";

import { userStore } from "./stores/user.js";
import { natsStore } from "./stores/nats.js";

const uStore = userStore();
const nStore = natsStore();
const { connection } = storeToRefs(nStore);
const { loggedIn } = storeToRefs(uStore);

const route = useRoute();
const path = computed(() => route.path);

const sleep = (delay) => new Promise((resolve) => setTimeout(resolve, delay));

watch(connection, () => {
  if (nStore.connection !== null) {
    nStore.connection
      .request("natster.global.my.shares", "", { timeout: 5000 })
      .then((m) => {
        uStore.shares = JSONCodec().decode(m.data).data;
      })
      .catch((err) => {
        console.log(`problem with request: ${err.message}`);
      });
  }
});

watch(loggedIn, () => {
  if (loggedIn && path.value == "/login") {
    router.push({ path: "/" });
  }
});

onMounted(() => {
  nStore.connect();
});
</script>
