import { defineStore } from 'pinia'
import { userStore } from './user'
import { natsStore } from './nats'
import { createCurve } from 'nkeys.js'
import { connect, jwtAuthenticator, JSONCodec } from 'nats.ws'
import type { Catalog, File } from '../types/types.ts'
import { textFileStore } from './textfile'
import { notificationStore } from './notification'

export const catalogStore = defineStore('catalog', {
  state: () => ({
    numSelected: 0,
    catalogs: [] as Catalog[],
    pending_catalogs: [] as Catalog[],
    shares_init: false,
    pending_init: false
  }),
  actions: {
    subscribeToHeartbeats() {
      const nStore = natsStore()
      const sub = nStore.connection.subscribe('natster.global-events.>')
      ;(async () => {
        for await (const msg of sub) {
          let m = JSONCodec().decode(msg.data)
          this.setOnlineAndCatalogRevision(m.catalog, m.revision)
        }
        console.log('subscription closed')
      })()
    },
    setOnlineAndCatalogRevision(inCat, rev) {
      const nStore = notificationStore()
      var d = new Date(0)
      d.setUTCSeconds(rev)
      this.catalogs.forEach(async function (c, i) {
        if (c.name == inCat) {
          c.lastSeen = Date.now()
          if (c.status != rev) {
            c.status = rev
            nStore.setNotification(
              'New Catalog Content!',
              'The catalog ' + c.name + ' has published new content'
            )
            if (c.selected) {
              natsStore()
                .connection.request('natster.catalog.' + c.name + '.get', '', { timeout: 5000 })
                .then((m) => {
                  c.files = [] as File[]
                  c.files.push(...JSONCodec().decode(m.data).data.entries)
                })
                .catch((err) => {
                  console.error('nats requestCatalogFiles err: ', err)
                })
            }
          }
        }
      })
    },
    setCatalogSelected(cat) {
      let selectedDiff = 0
      this.catalogs.forEach(async function (item, index) {
        if (cat.name == item.name) {
          if (item.selected) {
            selectedDiff = -1
            item.files = [] as File[]
            item.selected = false
          } else {
            natsStore()
              .connection.request('natster.catalog.' + cat.name + '.get', '', { timeout: 5000 })
              .then((m) => {
                item.files.push(...JSONCodec().decode(m.data).data.entries)
                item.selected = true
              })
              .catch((err) => {
                console.error('nats requestCatalogFiles err: ', err)
              })
            selectedDiff = 1
          }
        }
      })
      this.numSelected += selectedDiff
    },
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
                  description: 'PLACEHOLDER',
                  to: c.to_account,
                  from: c.from_account,
                  name: c.catalog,
                  online: c.catalog_online,
                  lastSeen: Date.now(),
                  pending_invite: false,
                  status: c.revision,
                  files: []
                }

                this.catalogs.push(catalog)
              }
            })
            this.shares_init = true
          }
        })
        .catch((err) => console.error('nats shares err: ', err))
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
            this.pending_init = true
          }
        })
        .catch((err) => console.error('nats ping err: ', err))
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
        })
        .catch((err) => {
          console.error('nats requestCatalogFiles err: ', err)
        })
    }
  },
  getters: {
    getNumCatalogsSelected(state) {
      return state.numSelected
    },
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

        if (!tCatalog.online && Date.now() - tCatalog.lastSeen < 1 * 60 * 1000) {
          tCatalog.online = true
          notificationStore().setNotification('Catalog Online', tCatalog.name + ' has come online')
        } else if (tCatalog.online && Date.now() - tCatalog.lastSeen > 1 * 60 * 1000) {
          tCatalog.online = false
          notificationStore().setNotification(
            'Catalog Offline',
            tCatalog.name + ' has gone offline'
          )
        }
      })

      return state.catalogs.filter((c) => !c.pending_invite)
    }
  }
})
