<template>
  <div class="text-xs font-semibold leading-6 text-gray-400">Natster Shares</div>
  <ul role="list" class="-mx-2 mt-2 space-y-1">
    <li v-for="(catalog, i) in getImportedCatalogs" :key="i">
      <div
        @click="uStore.setCatalogSelected(catalog)"
        :class="[
          catalog.selected
            ? 'bg-gray-800 text-white'
            : 'text-gray-400 hover:text-white hover:bg-gray-800',
          'group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold'
        ]"
      >
        <span
          class="relative inline-block flex h-6 w-6 shrink-0 items-center justify-center rounded-lg border border-gray-700 bg-gray-800 text-[0.625rem] font-medium text-gray-400"
        >
          <span
            class="absolute right-0 top-0 block h-1.5 w-1.5 -translate-y-1/2 translate-x-1/2 transform rounded-full ring-1 ring-white"
            :class="[catalog.online ? 'bg-green-500' : 'bg-gray-500']"
          />
          <span class=""> {{ catalog.name.substring(0, 1).toUpperCase() }} </span>
        </span>
        <span class="truncate">{{ catalog.name }} </span>
      </div>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'

import { userStore } from '../stores/user'

const { getImportedCatalogs } = storeToRefs(userStore())
</script>
