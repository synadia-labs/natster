<template>
  <ul v-for="catalog in catalogs" role="list" class="divide-y divide-gray-100">
    <li v-for="file in catalog.files" :key="file.hash" class="flex justify-between gap-x-6 py-5">
      <div class="flex min-w-0 gap-x-4">
        <img
          class="h-12 w-12 flex-none rounded-full bg-gray-50"
          :src="catalogImage(catalog.name)"
          alt=""
        />
        <div class="min-w-0 flex-auto">
          <p class="text-sm font-semibold leading-6 text-gray-900">
            {{ getFileName(file.path) }}
          </p>
          <p class="mt-1 flex text-xs leading-5 text-gray-500">
            {{ catalog.name }}
          </p>
        </div>
      </div>
      <div class="flex shrink-0 items-center gap-x-6">
        <div class="hidden sm:flex sm:flex-col sm:items-end">
          <p class="text-sm leading-6 text-gray-900">
            {{ file.description }}
          </p>
          <p class="mt-1 flex text-xs leading-5 text-gray-500">
            {{ file.mime_type }} | {{ formatBytes(file.byte_size) }}
          </p>
        </div>
        <Menu as="div" class="relative flex-none">
          <div>
            <MenuButton class="-m-2.5 block p-2.5 text-gray-500 hover:text-gray-900">
              <span class="sr-only">Open options</span>
              <EllipsisVerticalIcon class="h-5 w-5" aria-hidden="true" />
            </MenuButton>
          </div>

          <transition
            enter-active-class="transition duration-100 ease-out"
            enter-from-class="transform scale-95 opacity-0"
            enter-to-class="transform scale-100 opacity-100"
            leave-active-class="transition duration-75 ease-in"
            leave-from-class="transform scale-100 opacity-100"
            leave-to-class="transform scale-95 opacity-0"
          >
            <MenuItems
              class="absolute right-0 mt-2 w-56 origin-top-right divide-y divide-gray-100 rounded-md bg-white shadow-lg ring-1 ring-black/5 focus:outline-none"
            >
              <div class="px-1 py-1">
                <MenuItem v-slot="{ active }">
                  <button
                    :disabled="isDisabled"
                    :class="[
                      active && !isDisabled
                        ? 'bg-violet-500 text-white'
                        : 'text-gray-400 cursor-not-allowed',
                      'group flex w-full items-center rounded-md px-2 py-2 text-sm disabled:bg-blue-100'
                    ]"
                  >
                    <FolderOpenIcon
                      :active="active"
                      class="mr-2 h-5 w-5 text-violet-400"
                      aria-hidden="true"
                    />
                    View
                  </button>
                </MenuItem>
                <MenuItem v-slot="{ active }">
                  <button
                    @click.prevent="uStore.downloadFile(catalog.name, file.hash)"
                    :class="[
                      active ? 'bg-violet-500 text-white' : 'text-gray-900',
                      'group flex w-full items-center rounded-md px-2 py-2 text-sm'
                    ]"
                  >
                    <ArrowDownTrayIcon
                      :active="active"
                      class="mr-2 h-5 w-5 text-violet-400"
                      aria-hidden="true"
                    />
                    Download
                  </button>
                </MenuItem>
              </div>
            </MenuItems>
          </transition>
        </Menu>
      </div>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { jwtDecode } from 'jwt-decode'
import { userStore } from '../stores/user'
import { natsStore } from '../stores/nats'
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { EllipsisVerticalIcon, ArrowDownTrayIcon, FolderOpenIcon } from '@heroicons/vue/20/solid'

const uStore = userStore()
const nStore = natsStore()

const isDisabled = computed(() => true)

const { catalogs } = storeToRefs(uStore)
function catalogImage(name) {
  return 'https://ui-avatars.com/api/?name=+' + name
}
function getFileName(filepath) {
  let sFilePath = filepath.split('/')
  return sFilePath[sFilePath.length - 1]
}
function formatBytes(bytes, decimals = 2) {
  if (!+bytes) return '0 Bytes'

  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
}
</script>
