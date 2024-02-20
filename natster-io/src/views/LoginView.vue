<template>
  <div class="bg-gray-900">
    <div class="relative isolate overflow-hidden">
      <img
        src="https://images.unsplash.com/photo-1521737604893-d14cc237f11d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2830&q=80&blend=111827&sat=-100&exp=15&blend-mode=multiply"
        alt=""
        class="absolute inset-0 -z-10 h-full w-full object-cover"
      />
      <div
        class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
        aria-hidden="true"
      >
        <div
          class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-20 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"
          style="
            clip-path: polygon(
              74.1% 44.1%,
              100% 61.6%,
              97.5% 26.9%,
              85.5% 0.1%,
              80.7% 2%,
              72.5% 32.5%,
              60.2% 62.4%,
              52.4% 68.1%,
              47.5% 58.3%,
              45.2% 34.5%,
              27.5% 76.7%,
              0.1% 64.9%,
              17.9% 100%,
              27.6% 76.8%,
              76.1% 97.7%,
              74.1% 44.1%
            );
          "
        />
      </div>
      <div class="container relative mx-auto mx-center py-32 sm:py-48 lg:py-56">
        <div class="mx-full w-3/6 items-center content-center">
          <h1 class="text-4xl font-bold tracking-tight text-white sm:text-6xl">
            Login to start using Natster
          </h1>
          You will need to initiate the login workflow by running
          <pre>
            <code class="text-green-500 bg-gray-800">
curl -lSs https://natster.com/install.sh | sh
natster init --token YOUR_SYNADIA_CLOUD_TOKEN
natster web-login</code>
            </pre>
          from the command line. This will provide you with a link/code you can enter below.
        </div>
        <div v-if="!codeProvided" class="mt-10 2xl:h-48 xl:h-20 w-auto gap-x-6">
          <div class="flex flex-col-6 lg:h-full gap-5">
            <div
              class="flex-1 rounded-md shadow-sm ring-1 ring-inset ring-gray-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-indigo-600 sm:max-w-md bg-gray-50 bg-opacity-50"
            >
              <input
                type="text"
                name="code1"
                id="code1"
                class="block flex-1 border-0 bg-transparent py-1.5 pl-1 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6 text-center text-xl"
                pattern="[a-zA-Z0-9]{1}$"
                required
              />
            </div>
          </div>
          <div class="flex justify-end pt-10">
            <button
              type="button"
              class="rounded-md bg-indigo-500 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
            >
              Login
            </button>
          </div>
        </div>
      </div>
    </div>
    <div
      class="absolute inset-x-0 top-[calc(100%-13rem)] -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[calc(100%-30rem)]"
      aria-hidden="true"
    >
      <div
        class="relative left-[calc(50%+3rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-20 sm:left-[calc(50%+36rem)] sm:w-[72.1875rem]"
        style="
          clip-path: polygon(
            74.1% 44.1%,
            100% 61.6%,
            97.5% 26.9%,
            85.5% 0.1%,
            80.7% 2%,
            72.5% 32.5%,
            60.2% 62.4%,
            52.4% 68.1%,
            47.5% 58.3%,
            45.2% 34.5%,
            27.5% 76.7%,
            0.1% 64.9%,
            17.9% 100%,
            27.6% 76.8%,
            76.1% 97.7%,
            74.1% 44.1%
          );
        "
      />
    </div>
  </div>
  <div></div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Dialog, DialogPanel } from '@headlessui/vue'
import { Bars3Icon, XMarkIcon } from '@heroicons/vue/24/outline'
import { userStore } from '../stores/user.js'
import { useAuth0 } from '@auth0/auth0-vue'

const codeProvided = computed(() => {
  const route = useRoute()
  if (route.params.code === undefined || route.params.code === '') {
    return false
  }

  const { loginWithRedirect } = useAuth0()
  loginWithRedirect({
    appState: {
      target: '/library',
      nats_code: route.params.code
    },
    authorizationParams: {
      nats_code: route.params.code
    }
  })

  return true
})
</script>
