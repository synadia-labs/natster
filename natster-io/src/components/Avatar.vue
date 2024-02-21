<template>
  <Menu as="div" class="relative inline-block text-left">
    <Float placement="top-end">
      <div>
        <MenuButton class="flex items-center no-underline focus:outline-none">
          <span class="relative inline-block">
            <img class="h-12 w-12 rounded-md" :src="userPhotoUrl" alt="" />
            <span
              class="absolute right-0 top-0 block h-4 w-4 -translate-y-1/2 translate-x-1/2 transform rounded-full bg-yellow-300 text-center text-xs align-top text-black ring-1 ring-white"
            >
              2</span
            >
          </span>
        </MenuButton>
      </div>

      <transition
        enter-active-class="transition ease-out duration-100"
        enter-from-class="transform opacity-0 scale-95"
        enter-to-class="transform opacity-100 scale-100"
        leave-active-class="transition ease-in duration-75"
        leave-from-class="transform opacity-100 scale-100"
        leave-to-class="transform opacity-0 scale-95"
      >
        <MenuItems
          class="absolute right-0 z-10 mt-2 w-56 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none bottom-full"
        >
          <div class="py-1">
            <MenuItem v-slot="{ active }">
              <button
                type="submit"
                @click.prevent="signout"
                :class="[
                  active ? 'bg-gray-100 text-gray-900' : 'text-gray-700',
                  'block w-full px-4 py-2 text-left text-sm'
                ]"
              >
                Sign out
              </button>
            </MenuItem>
          </div>
        </MenuItems>
      </transition>
    </Float>
  </Menu>
  <div>
    <p @click="was, getXKeys" class="text-white font-semibold" aria-hidden="true">
      {{ user }}
    </p>
    <p @click="copyAccountIdToClipboard" class="text-gray-500" aria-hidden="true">
      {{ natster_account.substring(0, 8) }}...
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Float } from '@headlessui-float/vue'
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { useAuth0 } from '@auth0/auth0-vue'
import { notificationStore } from '../stores/notification'

const { logout } = useAuth0()
const nStore = notificationStore()

function signout() {
  logout({ logoutParams: { returnTo: window.location.origin } })
}

const userPhotoUrl = computed(() => {
  if (props.photo === undefined || props.photo === '') {
    return 'https://ui-avatars.com/api/?name=' + props.name
  }
  return props.photo
})

const props = defineProps({
  user: String,
  photo: String,
  natster_account: String
})

function copyAccountIdToClipboard() {
  navigator.clipboard.writeText(props.natster_account)
  notificationStore().setNotification('Copied!', 'Account ID copied to clipboard')
}
</script>
