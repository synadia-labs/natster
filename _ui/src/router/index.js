import { createRouter, createWebHistory } from "vue-router";

import MyLibrary from "../components/MyLibrary.vue";
import MyShares from "../components/MyShares.vue";
import Login from "../components/Login.vue";
import Home from "../components/Home.vue";

import { userStore } from "../stores/user";

const routes = [
  { path: "/", name: "Home", component: Home },
  { path: "/login", name: "Login", component: Login },
  { path: "/library", name: "Library", component: MyLibrary },
  { path: "/shares", name: "Shares", component: MyShares },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, from) => {
  const uStore = userStore();
  if (!uStore.loggedIn && to.name === "Home") {
    return;
  }
  if (!uStore.loggedIn && to.name !== "Login") {
    return { name: "Login" };
  }
  if (uStore.loggedIn && to.name == "Login") {
    return { name: "Library" };
  }
});

export default router;
