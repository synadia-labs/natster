<template>
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
                    <ul role="list" class="-mx-2 space-y-1">
                      <li v-for="item in navigation" :key="item.name">
                        <router-link
                          @click.native="toggleSelectedNav()"
                          :to="item.href"
                          :class="[
                            item.current
                              ? 'bg-gray-800 text-white'
                              : 'text-gray-400 hover:text-white hover:bg-gray-800',
                            'group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold'
                          ]"
                        >
                          <component :is="item.icon" class="h-6 w-6 shrink-0" aria-hidden="true" />
                          {{ item.name }}
                        </router-link>
                      </li>
                    </ul>
                  </li>
                  <li>
                    <div class="text-xs font-semibold leading-6 text-gray-400">Your Libraries</div>
                    <ul role="list" class="-mx-2 mt-2 space-y-1">
                      <li v-for="team in teams" :key="team.name">
                        <div
                          class="group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold text-gray-400"
                        >
                          <span
                            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-lg border border-gray-700 bg-gray-800 text-[0.625rem] font-medium text-gray-400"
                            >{{ team.initial }}</span
                          >
                          <span class="truncate">{{ team.name }}</span>
                        </div>
                      </li>
                    </ul>
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
            <ul role="list" class="-mx-2 space-y-1">
              <li v-for="item in navigation" :key="item.name">
                <router-link
                  @click="toggleSelectedNav()"
                  :to="item.href"
                  :class="[
                    item.current
                      ? 'bg-gray-800 text-white'
                      : 'text-gray-400 hover:text-white hover:bg-gray-800',
                    'group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold'
                  ]"
                >
                  <component :is="item.icon" class="h-6 w-6 shrink-0" aria-hidden="true" />
                  {{ item.name }}
                </router-link>
              </li>
            </ul>
          </li>
          <li>
            <div class="text-xs font-semibold leading-6 text-gray-400">Natster Libraries</div>
            <ul role="list" class="-mx-2 mt-2 space-y-1">
              <li v-for="team in teams" :key="team.name">
                <div
                  class="group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold text-gray-400"
                >
                  <span
                    class="relative inline-block flex h-6 w-6 shrink-0 items-center justify-center rounded-lg border border-gray-700 bg-gray-800 text-[0.625rem] font-medium text-gray-400"
                  >
                    <span
                      class="absolute right-0 top-0 block h-1.5 w-1.5 -translate-y-1/2 translate-x-1/2 transform rounded-full ring-1 ring-white"
                      :class="[team.online ? 'bg-green-500' : 'bg-gray-500']"
                    />
                    <span class="">{{ team.initial }} </span>
                  </span>
                  <span class="truncate">{{ team.name }} </span>
                </div>
              </li>
            </ul>
          </li>
          <li class="-mx-6 mt-auto">
            <div
              class="flex flex-row-reverse items-center gap-x-4 px-6 py-3 text-sm inset-y-0 right-0"
            >
              <Avatar>
                <span class="relative inline-block">
                  <img class="h-12 w-12 rounded-md" :src="uStore.getUserPhotoUrl" alt="" />
                  <span
                    class="absolute right-0 top-0 block h-4 w-4 -translate-y-1/2 translate-x-1/2 transform rounded-full bg-yellow-300 text-center text-xs align-top text-black ring-1 ring-white"
                  >
                    2</span
                  >
                </span>
              </Avatar>
              <div>
                <p class="text-white font-semibold" aria-hidden="true">
                  {{ uStore.getUserName }}
                </p>
                <p class="text-gray-500" aria-hidden="true">
                  {{ uStore.getAccount.substring(0, 8) }}...
                </p>
              </div>
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
    <a href="#">
      <span class="sr-only">Your profile</span>
      <img
        class="h-8 w-8 rounded-full bg-gray-800"
        src="https://avatars.githubusercontent.com/u/15827604?v=4"
        alt=""
      />
    </a>
  </div>

  <main class="py-10 lg:pl-72">
    <div class="px-4 sm:px-6 lg:px-8">
      <div>
        <TopBar :user="libraryName" />
        {{ uStore.getUser }}
        {{ uStore.getNatsCode }}
      </div>
      <div>
        <slot></slot>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import TopBar from '../components/TopBar.vue'
import Avatar from '../components/Avatar.vue'
import { userStore } from '../stores/user.js'
const uStore = userStore()

import { useAuth0 } from '@auth0/auth0-vue'
const { isLoading } = 'useAuth0'

import { Dialog, DialogPanel, TransitionChild, TransitionRoot } from '@headlessui/vue'
import {
  Bars3Icon,
  CalendarIcon,
  ChartPieIcon,
  DocumentDuplicateIcon,
  FolderIcon,
  HomeIcon,
  UsersIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

const natsterImg = new URL('../assets/natster.svg', import.meta.url)
const libraryName = computed(() => uStore.getUserName + "'s Library")
const navigation = reactive([
  { name: 'My Files', href: '/library', icon: HomeIcon, current: true },
  { name: 'My Shares', href: '/shares', icon: FolderIcon, current: false }
])
const teams = [
  {
    id: 1,
    name: 'Synadia Global Share',
    href: '#',
    initial: 'S',
    online: true,
    current: false
  },
  {
    id: 2,
    name: 'KevBuzz',
    href: '#',
    initial: 'KB',
    online: false,
    current: false
  },
  {
    id: 3,
    name: 'KylesSlowJams',
    href: '#',
    initial: 'KSJ',
    online: true,
    current: false
  },
  {
    id: 4,
    name: 'Dope Filez',
    href: '#',
    initial: 'DF',
    online: true,
    current: false
  },
  {
    id: 5,
    name: 'MySecrets',
    href: '#',
    initial: 'MS',
    online: true,
    current: false
  }
]

function toggleSelectedNav() {
  var arrayLength = navigation.length
  for (var i = 0; i < arrayLength; i++) {
    navigation[i].current = false
  }
}

const sidebarOpen = ref(false)
</script>
