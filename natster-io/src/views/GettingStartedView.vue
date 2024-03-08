<template>
  <div class="relative isolate overflow-hidden bg-white px-6 py-24 sm:py-32 lg:overflow-visible lg:px-0">
    <div class="absolute inset-0 -z-10 overflow-hidden">
      <svg
        class="absolute left-[max(50%,25rem)] top-0 h-[64rem] w-[128rem] -translate-x-1/2 stroke-gray-200 [mask-image:radial-gradient(64rem_64rem_at_top,white,transparent)]"
        aria-hidden="true">
        <defs>
          <pattern id="e813992c-7d03-4cc4-a2bd-151760b470a0" width="200" height="200" x="50%" y="-1"
            patternUnits="userSpaceOnUse">
            <path d="M100 200V.5M.5 .5H200" fill="none" />
          </pattern>
        </defs>
        <svg x="50%" y="-1" class="overflow-visible fill-gray-50">
          <path d="M-100.5 0h201v201h-201Z M699.5 0h201v201h-201Z M499.5 400h201v201h-201Z M-300.5 600h201v201h-201Z"
            stroke-width="0" />
        </svg>
        <rect width="100%" height="100%" stroke-width="0" fill="url(#e813992c-7d03-4cc4-a2bd-151760b470a0)" />
      </svg>
    </div>
    <div
      class="mx-auto grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 lg:mx-0 lg:max-w-none lg:grid-cols-2 lg:items-start lg:gap-y-10">
      <div
        class="lg:col-span-2 lg:col-start-1 lg:row-start-1 lg:mx-auto lg:grid lg:w-full lg:max-w-7xl lg:grid-cols-2 lg:gap-x-8 lg:px-8">
        <div class="lg:pr-4">
          <div class="lg:max-w-lg">
            <img :src="natsterImg" class="h-20" />
            <p class="mt-6 text-xl leading-8 text-gray-700">
              Natster is a secure peer-to-multipeer, decentralized media sharing platform.</p>
          </div>
        </div>
      </div>
      <div class="-ml-12 -mt-12 p-12 lg:sticky lg:top-4 lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:overflow-hidden">
        <h2 class="lg:pt-15 pb-2 mt-16 text-2xl font-bold tracking-tight text-gray-900">Getting Started with Natster
        </h2>
        <VCodeBlock :code="code_init" highlightjs lang="bash" theme="tokyo-night-light" copyTab />
        <br />
        <h2 class="pb-2 mt-10 text-2xl font-bold tracking-tight text-gray-900">Create your first Natster share</h2>
        <VCodeBlock :code="code_share" highlightjs lang="bash" theme="tokyo-night-light" copyTab />
      </div>
      <div
        class="lg:col-span-2 lg:col-start-1 lg:row-start-2 lg:mx-auto lg:grid lg:w-full lg:max-w-7xl lg:grid-cols-2 lg:gap-x-8 lg:px-8">
        <div class="lg:pr-4">
          <div class="max-w-xl text-base leading-7 text-gray-700 lg:max-w-lg">
            <p>Natster is an example of the kind of application you can build quickly and easily using nothing but NATS.
            </p>
            <ul role="list" class="mt-8 space-y-8 text-gray-600">
              <li class="flex gap-x-3">
                <CloudArrowUpIcon class="mt-1 h-5 w-5 flex-none text-indigo-600" aria-hidden="true" />
                <span><strong class="font-semibold text-gray-900">Share anywhere.</strong>
                  Privately and securely share your media with anyone, anywhere, without giving up control over your
                  data and its location.
                </span>
              </li>
              <li class="flex gap-x-3">
                <LockClosedIcon class="mt-1 h-5 w-5 flex-none text-indigo-600" aria-hidden="true" />
                <span><strong class="font-semibold text-gray-900">Always encrypted.</strong>
                  Nothing you share can ever be read by anyone but the intended recipient.
                </span>
              </li>
              <li class="flex gap-x-3">
                <CogIcon class="mt-1 h-5 w-5 flex-none text-indigo-600" aria-hidden="true" />
                <span><strong class="font-semibold text-gray-900">Nothing But NATS.</strong>
                  NATS is more than a tool, it's a platform that gives us an easy button for building distributed apps.
                </span>
              </li>
              <li class="flex gap-x-3">
                <ServerIcon class="mt-1 h-5 w-5 flex-none text-indigo-600" aria-hidden="true" />
                <span><strong class="font-semibold text-gray-900">Powered by Nex.</strong>
                  All backend services are managed by the NATS execution engine.
                </span>
              </li>
            </ul>
            <p class="mt-8"></p>
            <h2 class="mt-16 text-2xl font-bold tracking-tight text-gray-900">Need help? Slack us!</h2>
            <p class="mt-6">The maintainers of Natster can be found hanging out in the <a href="https://slack.nats.io"
                class="text-blue-500 underline">NATS.io slack</a>.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CloudArrowUpIcon, LockClosedIcon, CogIcon, ServerIcon } from '@heroicons/vue/20/solid'
import { onMounted } from 'vue'
import { ref } from 'vue'
import Notification from '../components/Notification.vue'
import { useRoute } from 'vue-router'
import { notificationStore } from '../stores/notification'
import VCodeBlock from '@wdns/vue-code-block';
const natsterImg = new URL('@/assets/natster-horizontal.svg', import.meta.url)

const code_init = ref(`# Install the Natster CLI
curl -sSf https://natster.io/install.sh | sh

# Initialize the Natster with your Synadia Cloud Token
natster init --token <SYNADIA CLOUD TOKEN>

# Bind your OAuth ID with your Natster Account
natster login

# Verify your context was successfully bound
natster whoami
`);

const code_share = ref(`# Create a new catalog
natster catalog new

# Serve your catalog
natster catalog serve

# Share with your friends
natster catalog share
`);
const route = useRoute()

onMounted(() => {
  if (route.query.failed) {
    notificationStore().setNotification(
      'Failed to login',
      'You tried to use an unbound oauth_id, please follow these instructions to login',
    )
  }
})

console.log(route)
</script>
