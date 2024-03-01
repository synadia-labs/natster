<template>
  <li :key="file.hash" class="flex items-center justify-between gap-x-6 py-5 pl-16">
    <div class="min-w-0">
      <div class="flex items-start gap-x-3">
        <p class="text-sm font-semibold leading-6 text-gray-900">{{ getFileName(file.path) }}</p>
      </div>
      <div class="mt-1 flex items-center gap-x-2 text-xs leading-5 text-gray-500">
        <p class="whitespace-nowrap">{{ file.hash }}</p>
      </div>
    </div>
    <div class="flex flex-none items-center gap-x-4">
    <div class="min-w-0">
      <div class="flex justify-end gap-x-3">
        <p class="text-sm leading-6 text-gray-900">{{ file.description }}</p>
      </div>
      <div class="mt-1 flex items-center gap-x-2 text-xs leading-5 text-gray-500">
        <p class="whitespace-nowrap">{{ file.mime_type }} | {{ formatBytes(file.byte_size)}}</p>
      </div>
    </div>

      <Menu as="div" class="relative inline-block flex-none z-10">
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
                  <ArrowDownTrayIcon
                    :active="active"
                    class="mr-2 h-5 w-5 text-violet-400"
                    aria-hidden="true"
                  />
                  Download
                </button>
              </MenuItem>
              <MenuItem v-slot="{ active }">
                <button
                  @click.prevent="cStore.viewFile(getFileName(file.path), catalog.name, file.hash)"
                  :class="[
                    active ? 'bg-violet-500 text-white' : 'text-gray-900',
                    'group flex w-full items-center rounded-md px-2 py-2 text-sm'
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
            </div>
          </MenuItems>
        </transition>
      </Menu>
    </div>
  </li>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { initFlowbite } from 'flowbite'
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { EllipsisVerticalIcon, ArrowDownTrayIcon, FolderOpenIcon } from '@heroicons/vue/20/solid'
import type File from '../types/types'
import { catalogStore } from '../stores/catalog'

const cStore = catalogStore()
const isDisabled = ref(true)

onMounted(() => {
  initFlowbite()
})

const props = defineProps<{
  file: File
  catalog: String
}>()

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
