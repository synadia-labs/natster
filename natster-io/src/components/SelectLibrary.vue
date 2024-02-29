<template>
  <div class="relative pt-10">
    <div class="relative h-80 overflow-hidden md:absolute md:left-0 md:h-full md:w-1/3 lg:w-1/2">
      <img class="h-full w-full" :src="natsterImg" alt="" />
    </div>
    <div class="relative mx-auto max-w-7xl py-24 sm:py-32 lg:px-8 lg:py-40">
      <div class="pl-6 pr-6 md:ml-auto md:w-2/3 md:pl-16 lg:w-1/2 lg:pl-24 lg:pr-0 xl:pl-32">
        <h2 class="text-base font-semibold leading-7 text-indigo-400">Thanks for using Natster!</h2>
        <p class="mt-2 text-3xl font-bold tracking-tight sm:text-4xl">
          Choose from one of your shared catalogs
        </p>
        <p class="mt-6 text-base leading-7">
          Any catalog that has been shared and you have imported will show up in the left sidebar.
          You can select a catalog to view the files in that catalog.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const natsterImg = new URL('@/assets/natster.svg', import.meta.url)
</script>
