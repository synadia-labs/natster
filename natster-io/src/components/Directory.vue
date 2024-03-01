<template>
  <div class="relative flex justify-between gap-x-6 py-5">
    <div
      class="flex min-w-0 gap-x-4"
      :id="'accordion-heading-' + normalizeString(catalog.name + directory)"
      data-accordion="collapse"
      data-active-classes="dark:bg-gray-900 text-gray-900 dark:text-white"
      data-inactive-classes="text-gray-500 dark:text-gray-400"
    >
      <img
        class="h-12 w-12 flex-none rounded-full bg-gray-50"
        :src="catalogImage(catalog.name)"
        alt=""
      />
      <div class="min-w-0 flex-auto">
        <p class="text-sm font-semibold leading-6 text-gray-900">
          {{ directory == 'root' ? '/' : directory }}
        </p>
        <p class="mt-1 flex text-xs leading-5 text-gray-500">{{ files.length != 1 ? files.length + ' files' : files.length + ' file'}}</p>
      </div>
    </div>
    <div class="flex shrink-0 items-center gap-x-4">
      <div class="hidden sm:flex sm:flex-col sm:items-end">
        <p class="text-sm leading-6 text-gray-900">{{ catalog.description }}</p>
        <p class="mt-1 text-xs leading-5 text-gray-500">
          {{ catalog.name }}
        </p>
      </div>
      <button
        type="button"
        aria-expanded="false"
        :data-accordion-target="'#accordion-flush-' + normalizeString(catalog.name + directory)"
        :aria-controls="'accordion-flush-' + normalizeString(catalog.name + directory)"
      >
        <svg
          data-accordion-icon
          class="h-5 w-5 rotate-180 flex-none text-gray-400"
          viewBox="0 0 10 6"
          fill="none"
          aria-hidden="true"
        >
          <path fill-rule="evenodd" stroke="currentColor" d="M9 5 5 1 1 5" clip-rule="evenodd" />
        </svg>
      </button>
    </div>
  </div>
  <div
    :id="'accordion-flush-' + normalizeString(catalog.name + directory)"
    class="hidden"
    :aria-labelledby="'accordion-heading-' + normalizeString(catalog.name + directory)"
  >
    <ul role="list" class="divide-y divide-gray-100">
      <FileComp v-for="file in files" :key="file.hash" :catalog="catalog" :file="file" />
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { initFlowbite } from 'flowbite'
import type { Catalog, File } from '../types/types.ts'
import FileComp from './FileComp.vue'

onMounted(() => {
  initFlowbite()
})

const props = defineProps<{
  catalog: Catalog
  directory: String
  files: File[]
}>()

function catalogImage(name) {
  return 'https://ui-avatars.com/api/?name=+' + name
}

function normalizeString(instring) {
  return instring.replace(/[^a-zA-Z]/g, '')
}
</script>
