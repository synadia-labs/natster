<template>
  <SelectLibrary v-if="getNumCatalogsSelected == 0" />
  <div v-else>
    <div v-for="catalog in getImportedCatalogs" role="list" class="divide-y divide-gray-100">
      <div
        v-for="(files, directory) in getFilesByDirectory(catalog.files)"
        :key="directory"
        :id="'accordion-' + normalizeString(catalog.name + directory)"
        data-accordion="collapse"
      >
        <Directory
          v-if="directory != 'root'"
          :catalog="catalog"
          :directory="directory"
          :files="files"
        />
        <FileComp
          v-else
          v-for="file in files"
          :key="file.hash"
          :catalog="catalog"
          :file="file"
          :image="true"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { initFlowbite } from 'flowbite'
import { storeToRefs } from 'pinia'
import { jwtDecode } from 'jwt-decode'
import { natsStore } from '../stores/nats'
import { catalogStore } from '../stores/catalog'
import type { File } from '../types/types.ts'
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { EllipsisVerticalIcon, ArrowDownTrayIcon, FolderOpenIcon } from '@heroicons/vue/20/solid'
import FileComp from './FileComp.vue'

import SelectLibrary from './SelectLibrary.vue'
import Directory from './Directory.vue'

const nStore = natsStore()
const cStore = catalogStore()

const { getImportedCatalogs, getNumCatalogsSelected } = storeToRefs(cStore)

onMounted(() => {
  initFlowbite()
})

function getFilesByDirectory(files: File[]) {
  type directoryMap = {
    [id: string]: string[]
  }

  let dm: directoryMap = {}

  files.forEach((file) => {
    let tFilePath = file.path.split('/')
    if (tFilePath.length > 1) {
      let directory = tFilePath[0]
      if (directory in dm) {
        dm[directory].push(file)
      } else {
        dm[directory] = [file]
      }
    } else if (tFilePath.length === 1) {
      if ('root' in dm) {
        dm['root'].push(file)
      } else {
        dm['root'] = [file]
      }
    }
  })

  return dm
}

function normalizeString(instring) {
  return instring.replace(/[^a-zA-Z]/g, '')
}
</script>
