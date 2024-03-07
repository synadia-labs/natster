<template>
  <TransitionRoot as="template" :show="show">
    <Dialog as="div" class="relative z-10" @close="close()">
      <TransitionChild
        as="template"
        enter="ease-out duration-300"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="ease-in duration-200"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
      </TransitionChild>

      <div class="fixed inset-0 z-10 w-screen overflow-y-auto">
        <div
          class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0"
        >
          <TransitionChild
            as="template"
            enter="ease-out duration-300"
            enter-from="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            enter-to="opacity-100 translate-y-0 sm:scale-100"
            leave="ease-in duration-200"
            leave-from="opacity-100 translate-y-0 sm:scale-100"
            leave-to="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          >
            <DialogPanel
              class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6"
            >
              <div class="absolute right-0 top-0 hidden pr-4 pt-4 sm:block">
                <button
                  type="button"
                  class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                  @click="close()"
                >
                  <span class="sr-only">Close</span>
                  <XMarkIcon class="h-6 w-6" aria-hidden="true" />
                </button>
              </div>
              <div class="sm:flex sm:items-start">
                <div class="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
                  <DialogTitle as="h3" class="text-base font-semibold leading-6 text-gray-900">
                    {{ catalog.name }} | {{ title }}</DialogTitle
                  >
                  <div class="mt-2">
                    <p v-if="!!body" class="text-sm text-gray-500">
                      {{ body }}
                    </p>
                  </div>
                </div>
              </div>
              <div class="relative">
                <VueSpinnerAudio 
                  v-if="loading"
                  size="80"
                  class="loading-spinner"
                />

                <video
                  v-if="!!mediaUrl && mimeType.toLowerCase() == 'video/mp4'"
                  v-show="!loading"
                  id="video"
                  :type="mimeType"
                  :src="mediaUrl"
                  width="640"
                  height="360"
                  autoplay
                  controls
                ></video>
                
                <AudioPlayer
                  v-if="!!mediaUrl && mimeType.toLowerCase() == 'audio/mpeg'"
                  v-show="!loading"
                  @loadedmetadata="playAudio(e)"
                  :option="
                    getAudioOptions(
                      mediaUrl,
                      description == '' ? title : description,
                      catalog.image
                    )
                  "
                  :title="title"
                />
              </div>
              <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
                <button
                  type="button"
                  class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto"
                  @click="close()"
                >
                  Close
                </button>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { storeToRefs } from 'pinia'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'
import { VueSpinnerAudio } from 'vue3-spinners'

import AudioPlayer from 'vue3-audio-player'
import 'vue3-audio-player/dist/style.css'

import { fileStore } from '../stores/file'
const fStore = fileStore()
const { body, title, show, loading, mimeType, mediaUrl, catalog, description } = storeToRefs(fStore)

function close() {
  console.log('closing file view')
  fStore.show = false
  fStore.reset()
}

watch(mimeType, (newVal, oldVal) => {
  console.log(`mime type changed... ${newVal}`)
  if (!!newVal && newVal.toLowerCase().indexOf('video/') === 0) {
    console.log('video incoming')
  }
})

function getAudioOptions(inSrc, inTitle, inCover) {
  return {
    src: inSrc,
    title: inTitle,
    coverImage: inCover
  }
}

function playAudio() {
  document.querySelector('audio').play()
}

watch(mediaUrl, (newVal, oldVal) => {
  if (!!newVal) {
    setTimeout(() => {
      try {
        document.querySelector('audio').addEventListener('play', (event) => {
          if (fStore.loading) {
            fStore.loading = false
          }
        })
      } catch (e) {}

      try {
        document.querySelector('video').addEventListener('play', (event) => {
          if (fStore.loading) {
            fStore.loading = false
          }
        })
      } catch (e) {}
    }, 50)
  }
})
</script>

<style>
.audio__player-play {
  display: block;
  margin-left: auto;
  margin-right: auto;
  width: 55%;
}

.loading-spinner {
  color: #45c320;
  margin: 50px auto;
}
</style>
