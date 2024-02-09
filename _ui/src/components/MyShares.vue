<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="mt-8 flow-root">
      <div class="-mx-4 -my-2 sm:-mx-6 lg:-mx-8">
        <div class="inline-block min-w-full py-2 align-middle">
          <table class="min-w-full border-separate border-spacing-0">
            <thead>
              <tr>
                <th
                  scope="col"
                  class="sticky top-0 z-10 border-b border-gray-300 bg-white bg-opacity-75 py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter sm:pl-6 lg:pl-8"
                >
                  Catalog
                </th>
                <th
                  scope="col"
                  class="sticky top-0 z-10 hidden border-b border-gray-300 bg-white bg-opacity-75 px-3 py-3.5 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter sm:table-cell"
                >
                  From Account
                </th>
                <th
                  scope="col"
                  class="sticky top-0 z-10 hidden border-b border-gray-300 bg-white bg-opacity-75 px-3 py-3.5 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter lg:table-cell"
                >
                  To Account
                </th>
                <th
                  scope="col"
                  class="sticky top-0 z-10 hidden border-b border-gray-300 bg-white bg-opacity-75 px-3 py-3.5 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter lg:table-cell"
                >
                  Shared
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(catalog, i) in uStore.myShares" :key="i">
                <td
                  :class="[
                    i !== uStore.myShares.length - 1
                      ? 'border-b border-gray-200'
                      : '',
                    'whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6 lg:pl-8',
                  ]"
                >
                  {{ catalog.catalog }}
                </td>
                <td
                  :class="[
                    i !== uStore.myShares.length - 1
                      ? 'border-b border-gray-200'
                      : '',
                    'whitespace-nowrap hidden px-3 py-4 text-sm text-gray-500 sm:table-cell',
                  ]"
                >
                  {{ catalog.from_account.substring(0, 8) }}...
                </td>
                <td
                  :class="[
                    i !== uStore.myShares.length - 1
                      ? 'border-b border-gray-200'
                      : '',
                    'whitespace-nowrap hidden px-3 py-4 text-sm text-gray-500 lg:table-cell',
                  ]"
                >
                  {{ catalog.to_account.substring(0, 8) }}...
                </td>
                <td
                  :class="[
                    i !== uStore.myShares.length - 1
                      ? 'border-b border-gray-200'
                      : '',
                    'whitespace-nowrap hidden px-3 py-4 text-sm text-gray-500 lg:table-cell',
                  ]"
                >
                  <span
                    :class="[
                      catalog.to_account
                        ? 'ring-green-150 bg-green-100'
                        : 'ring-red-150 bg-red-100',
                    ]"
                    class="inline-flex items-center rounded-md px-2 py-1 text-xs font-medium text-green-700 ring-1 ring-inset"
                    >{{ catalog.to_account == "" ? "No" : "Yes" }}</span
                  >
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { userStore } from "../stores/user.js";

const uStore = userStore();
</script>
