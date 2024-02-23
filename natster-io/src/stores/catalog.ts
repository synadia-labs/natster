import { defineStore } from 'pinia'
import { userStore } from './user'
import { natsStore } from './nats'
import { createCurve } from 'nkeys.js'
import { connect, jwtAuthenticator, JSONCodec } from 'nats.ws'
import type { Catalog, File } from '../types/types.ts'
import { textFileStore } from './textfile'

export const catalogStore = defineStore('catalog', {
  state: () => ({
    catalogs: [] as Catalog[],
    pending_catalogs: [] as Catalog[],
    shares_init: false,
    pending_init: false
  }),
  actions: {
    async getShares(init) {
      const nStore = natsStore()
      const uStore = userStore()
      await nStore.connection
        .request('natster.global.my.shares', '', { timeout: 5000 })
        .then((msg) => {
          let m = JSONCodec().decode(msg.data)
          if (m.code == 200) {
            m.data.forEach((c, i) => {
              if (c.to_account === uStore.getAccount) {
                const catalog: Catalog = {
                  selected: false,
                  to: c.to_account,
                  from: c.from_account,
                  name: c.catalog,
                  online: c.catalog_online,
                  pending_invite: false,
                  files: []
                }
                if (c.catalog == 'Synadia Hub' && init) {
                  catalog.selected = true
                }
                this.catalogs.push(catalog)
              }
            })
          }
        })
        .catch((err) => console.error('nats shares err: ', err))
        .finally(() => {
          this.shares_init = true
        })
    },
    async getLocalInbox() {
      const uStore = userStore()
      const nStore = natsStore()
      await nStore.connection
        .request('natster.local.inbox', '', { timeout: 5000 })
        .then((msg) => {
          let m = JSONCodec().decode(msg.data)
          if (m.code == 200) {
            uStore.catalog_online = true
            uStore.pending_imports = m.data.unimported_shares.length
            m.data.unimported_shares.forEach((c, i) => {
              if (c.to_account === uStore.getAccount) {
                const catalog: Catalog = {
                  to: c.to_account,
                  from: c.from_account,
                  name: c.catalog
                }
                this.pending_catalogs.push(catalog)
              }
            })
          }
        })
        .catch((err) => console.error('nats ping err: ', err))
        .finally(() => {
          this.pending_init = true
        })
    },
    setCatalogSelected(cat) {
      this.catalogs.forEach(async function (item, index) {
        if (cat.name == item.name) {
          if (item.selected) {
            item.files = [] as File[]
            item.selected = false
          } else {
            await natsStore()
              .connection.request('natster.catalog.' + cat.name + '.get', '', { timeout: 5000 })
              .then((m) => {
                item.files.push(...JSONCodec().decode(m.data).data.entries)
              })
              .catch((err) => {
                console.error('nats requestCatalogFiles err: ', err)
              })
              .finally(() => {
                item.selected = true
              })
          }
        }
      })
    },
    async viewFile(fileName, catalog, hash) {
      const tfStore = textFileStore()

      let xkey = createCurve()
      this.xkey_seed = new TextDecoder().decode(xkey.getSeed())
      this.xkey_pub = xkey.getPublicKey()

      var sender_xkey
      var fileArray
      const nStore = natsStore()
      const sub = nStore.connection.subscribe('natster.media.' + catalog + '.' + hash)
      ;(async () => {
        for await (const m of sub) {
          await new Promise((r) => setTimeout(r, 1000))
          let decrypted = xkey.open(m.data, sender_xkey)
          // let decrypted = Array.from(xkey.open(m.data, sender_xkey))
          // fileArray = [...decrypted]
          // console.log("decrypted: ", decrypted)
          // console.log("fileArray: ", fileArray)
          tfStore.showTextFile(fileName, new TextDecoder().decode(decrypted))
        }
        console.log('subscription closed')
      })()

      const dl_request = {
        hash: hash,
        target_xkey: this.xkey_pub
      }
      await nStore.connection
        .request('natster.catalog.' + catalog + '.download', JSON.stringify(dl_request), {
          timeout: 5000
        })
        .then((m) => {
          let data = JSONCodec().decode(m.data)
          console.log('data: ', data)
          sender_xkey = data.data.sender_xkey
          fileArray = new Array(data.data.total_bytes)
        })
        .catch((err) => {
          console.error('nats requestCatalogFiles err: ', err)
        })
    }
  },
  getters: {
    getCatalogs(state) {
      return state.catalogs
    },
    getPendingCatalogs(state) {
      return state.pending_catalogs
    },
    catalogsInitialized(state) {
      return state.shares_init && state.pending_init
    },
    getImportedCatalogs(state) {
      state.catalogs.forEach(function (tCatalog, index) {
        state.pending_catalogs.forEach(function (tPending, index) {
          if (tCatalog.name === tPending.name) {
            tCatalog.pending_invite = true
          }
        })
      })

      return state.catalogs.filter((c) => !c.pending_invite)
    }
  }
})
