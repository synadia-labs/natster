<template>
  <Notification />
  <TextFile />
  <TransitionRoot as="template" :show="sidebarOpen">
    <Dialog as="div" class="relative z-50 lg:hidden" @close="sidebarOpen = false">
      <TransitionChild
        as="template"
        enter="transition-opacity ease-linear duration-300"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="transition-opacity ease-linear duration-300"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-gray-900/80" />
      </TransitionChild>

      <div class="fixed inset-0 flex">
        <TransitionChild
          as="template"
          enter="transition ease-in-out duration-300 transform"
          enter-from="-translate-x-full"
          enter-to="translate-x-0"
          leave="transition ease-in-out duration-300 transform"
          leave-from="translate-x-0"
          leave-to="-translate-x-full"
        >
          <DialogPanel class="relative mr-16 flex w-full max-w-xs flex-1">
            <TransitionChild
              as="template"
              enter="ease-in-out duration-300"
              enter-from="opacity-0"
              enter-to="opacity-100"
              leave="ease-in-out duration-300"
              leave-from="opacity-100"
              leave-to="opacity-0"
            >
              <div class="absolute left-full top-0 flex w-16 justify-center pt-5">
                <button type="button" class="-m-2.5 p-2.5" @click="sidebarOpen = false">
                  <span class="sr-only">Close sidebar</span>
                  <XMarkIcon class="h-6 w-6 text-white" aria-hidden="true" />
                </button>
              </div>
            </TransitionChild>
            <!-- Sidebar component, swap this element with another sidebar if you like -->
            <div
              class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6 pb-2 ring-1 ring-white/10"
            >
              <div class="flex h-16 shrink-0 items-center">
                <img class="h-8 w-auto" :src="natsterImg" alt="Natster" />
              </div>
              <nav class="flex flex-1 flex-col">
                <ul role="list" class="flex flex-1 flex-col gap-y-7">
                  <li>
                    <Catalogs />
                  </li>
                </ul>
              </nav>
            </div>
          </DialogPanel>
        </TransitionChild>
      </div>
    </Dialog>
  </TransitionRoot>

  <!-- Static sidebar for desktop -->
  <div class="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-72 lg:flex-col">
    <!-- Sidebar component, swap this element with another sidebar if you like -->
    <div class="flex grow flex-col gap-y-5 overflow-y-auto bg-gray-900 px-6">
      <div class="flex h-20 shrink-0 items-center">
        <img class="h-12 w-auto" :src="natsterImg" alt="Natster" />
      </div>
      <nav class="flex flex-1 flex-col">
        <ul role="list" class="flex flex-1 flex-col gap-y-7">
          <li>
            <Catalogs />
          </li>
          <li class="-mx-6 mt-auto">
            <div
              class="flex flex-row-reverse items-center gap-x-4 px-6 py-3 text-sm inset-y-0 right-0"
            >
              <Avatar
                :user="user.name"
                :natster_account="uStore.getAccount"
                :photo="user.picture"
              />
            </div>
          </li>
        </ul>
      </nav>
    </div>
  </div>

  <div
    class="sticky top-0 z-40 flex items-center gap-x-6 bg-gray-900 px-4 py-4 shadow-sm sm:px-6 lg:hidden"
  >
    <button type="button" class="-m-2.5 p-2.5 text-gray-400 lg:hidden" @click="sidebarOpen = true">
      <span class="sr-only">Open sidebar</span>
      <Bars3Icon class="h-6 w-6" aria-hidden="true" />
    </button>
    <div class="flex-1 text-sm font-semibold leading-6 text-white">Dashboard</div>
    <Avatar :user="user.name" :natster_account="uStore.getAccount" :photo="user.picture" />
  </div>

  <main class="py-10 lg:pl-72">
    <div class="px-4 sm:px-6 lg:px-8">
      <div>
        <TopBar :user="user.name" />
      </div>
      <div>
        <slot></slot>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted, watch } from 'vue'
import { useAuth0 } from '@auth0/auth0-vue'
import { storeToRefs } from 'pinia'
import { Dialog, DialogPanel, TransitionChild, TransitionRoot } from '@headlessui/vue'

import TopBar from '../components/TopBar.vue'
import Avatar from '../components/Avatar.vue'
import Notification from '../components/Notification.vue'
import Catalogs from '../components/Catalogs.vue'
import TextFile from '../components/TextFile.vue'
import { Bars3Icon, XMarkIcon } from '@heroicons/vue/24/outline'

import { userStore } from '../stores/user'
import { natsStore } from '../stores/nats'

const { isLoading, isAuthenticated, user } = useAuth0()
const uStore = userStore()
const nStore = natsStore()
const { jwt, nkey, nats_code } = storeToRefs(uStore)

const natsterImg = new URL('../assets/natster.svg', import.meta.url)

onMounted(() => {
  if (isAuthenticated && user !== null) {
    uStore.setJWT(user.value.natster.jwt)
    uStore.setNkey(user.value.natster.nkey)
    uStore.setNatsCode(user.value.natster.natsCode)
  }
})

watch(jwt, (value) => {
  if (value !== '' && user !== null) {
    nStore.connect()
  }
})

const sidebarOpen = ref(false)
</script>
