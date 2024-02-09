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
                  File
                </th>
                <th
                  scope="col"
                  class="sticky top-0 z-10 hidden border-b border-gray-300 bg-white bg-opacity-75 px-3 py-3.5 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter sm:table-cell"
                >
                  Catalog
                </th>
                <th
                  scope="col"
                  class="sticky top-0 z-10 hidden border-b border-gray-300 bg-white bg-opacity-75 px-3 py-3.5 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter lg:table-cell"
                >
                  Description
                </th>
                <th
                  scope="col"
                  class="sticky top-0 z-10 hidden border-b border-gray-300 bg-white bg-opacity-75 px-3 py-3.5 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter lg:table-cell"
                >
                  Mime Type
                </th>
                <th
                  scope="col"
                  class="sticky top-0 z-10 hidden border-b border-gray-300 bg-white bg-opacity-75 px-3 py-3.5 text-left text-sm font-semibold text-gray-900 backdrop-blur backdrop-filter lg:table-cell"
                >
                  Size
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(file, i) in library" :key="i">
                <td
                  :class="[
                    i !== library.length - 1 ? 'border-b border-gray-200' : '',
                    'whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6 lg:pl-8',
                  ]"
                >
                  {{ file.name }}
                </td>
                <td
                  :class="[
                    i !== library.length - 1 ? 'border-b border-gray-200' : '',
                    'whitespace-nowrap hidden px-3 py-4 text-sm text-gray-500 sm:table-cell',
                  ]"
                >
                  {{ file.catalog }}
                </td>
                <td
                  :class="[
                    i !== library.length - 1 ? 'border-b border-gray-200' : '',
                    'whitespace-nowrap hidden px-3 py-4 text-sm text-gray-500 lg:table-cell',
                  ]"
                >
                  {{ file.description }}
                </td>
                <td
                  :class="[
                    i !== library.length - 1 ? 'border-b border-gray-200' : '',
                    'whitespace-nowrap hidden px-3 py-4 text-sm text-gray-500 lg:table-cell',
                  ]"
                >
                  {{ file.mimeType }}
                </td>
                <td
                  :class="[
                    i !== library.length - 1 ? 'border-b border-gray-200' : '',
                    'whitespace-nowrap hidden px-3 py-4 text-sm text-gray-500 lg:table-cell',
                  ]"
                >
                  {{ file.size }}
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
const library = [
  {
    name: "kubecon_plans.pdf",
    catalog: "Synadia Global Share",
    description: "World dominiation plans",
    mimeType: "file/pdf",
    size: "18273",
  },
  {
    name: "fireside_with_byron.wav",
    catalog: "Synadia Global Share",
    description: "Spend time with Byron and learn his ways",
    mimeType: "file/wav",
    size: "128478273",
  },
  {
    name: "Waterfalls - TLC",
    catalog: "KylesSlowJams",
    description: "Workin music",
    mimeType: "file/mp3",
    size: "128323",
  },
  {
    name: "LoveLetter.doc",
    catalog: "My Secrets",
    description: "The strongest feelings are ofter the hardest",
    mimeType: "file/doc",
    size: "287427",
  },
  {
    name: "BitcoinWallet",
    catalog: "Dope Filez",
    description: "Where I keep my lambo money",
    mimeType: "file/bin",
    size: "4",
  },
  {
    name: "Birthday Party Invite",
    catalog: "Dope Filez",
    description: "Its my party and ill cry if I want ot",
    mimeType: "file/pdg",
    size: "98723",
  },
];
</script>
