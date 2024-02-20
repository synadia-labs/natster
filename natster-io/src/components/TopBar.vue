<template>
  <div class="md:flex md:items-center md:justify-between pb-3">
    <div class="min-w-0 flex-1">
      <h2 class="text-2xl font-bold leading-7 text-black sm:truncate sm:text-3xl sm:tracking-tight">
        {{ libraryName }}
      </h2>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const libraryName = computed(() => props.user + "'s Library")

const props = defineProps({
  user: String
})
</script>
