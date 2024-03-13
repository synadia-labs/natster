import { defineStore } from 'pinia'
import { userStore } from './user'
import { natsStore } from './nats'
import { createCurve } from 'nkeys.js'
import { jwtAuthenticator, JSONCodec } from 'nats.ws'
import type { Catalog, File } from '../types/types.ts'
import { fileStore } from './file'
import { notificationStore } from './notification'
import { saveAs } from 'file-saver'

export const catalogStore = defineStore('catalog', {
  state: () => ({
    activeDownload: false,
    supportedMimeType: ['image/png', 'image/jpeg', 'video/mp4', 'text/plain', 'audio/mpeg'],
    numSelected: 0,
    catalogs: [] as Catalog[],
    pending_catalogs: [] as Catalog[],
    shares_init: false,
    pending_init: false
  }),
  actions: {
    subscribeToHeartbeats() {
      const nStore = natsStore()
      const uStore = userStore()

      const sub = nStore.connection.subscribe('natster.global-events.>')
        ; (async () => {
          for await (const msg of sub) {
            let m = JSONCodec().decode(msg.data)
            this.setOnlineAndCatalogRevision(m.catalog, m.revision)
            uStore.togglePing()

            this.scrubCatalogs()
          }
          console.log('subscription closed')
        })()
    },
    scrubCatalogs() {
      const disableCat = (catalog) => {
        this.setCatalogSelected(catalog)
      }

      this.catalogs.forEach(function(tCatalog) {

        // Toggles catalogs avaialability
        if (!tCatalog.online && Date.now() - new Date(tCatalog.lastSeen).getTime() < 1 * 60 * 1000) {
          tCatalog.online = true
          notificationStore().setNotification('Catalog Online', tCatalog.name + ' has come online')
        } else if (tCatalog.online && Date.now() - new Date(tCatalog.lastSeen).getTime() > 1 * 60 * 1000) {
          tCatalog.online = false
          notificationStore().setNotification(
            'Catalog Offline',
            tCatalog.name + ' has gone offline'
          )
          if (tCatalog.selected) {
            disableCat(tCatalog)
          }
        }
      })
    },
    subscribeToLocalHeartbeats() {
      const nStore = natsStore()
      const uStore = userStore()

      const sub = nStore.connection.subscribe('natster.local-events.heartbeat')
        ; (async () => {
          for await (const m of sub) {
            let msg = JSONCodec().decode(m.data)
            uStore.setLastSeenTS(new Date(Date.now()))
            this.setOnlineAndCatalogRevision(msg.catalog, msg.revision)
          }
          console.log('subscription closed')
        })()
    },
    setOnlineAndCatalogRevision(inCat, rev) {
      const nStore = notificationStore()
      var d = new Date(0)
      d.setUTCSeconds(rev)
      this.catalogs.forEach(async function(c, i) {
        if (c.name == inCat) {
          c.pending_invite = false
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
                  nStore.setNotification(
                    'Catalog failed to respond',
                    'The catalog ' + c.name + ' failed to respond. Moving offline until it heartbeats.'
                  )
                  if (c.selected) {
                    this.setCatalogSelected(c)
                  }
                  c.online = false
                })
            }
          }
        }
      })
    },
    setCatalogSelected(cat: Catalog) {
      const setDiff = (inDiff) => {
        this.numSelected += inDiff
      }

      const nStore = notificationStore()
      this.catalogs.forEach(async function(item) {
        if (cat.name == item.name) {
          if (item.selected) {
            item.files = [] as File[]
            item.selected = false
            setDiff(-1)
          } else {
            natsStore()
              .connection.request('natster.catalog.' + cat.name + '.get', '', { timeout: 5000 })
              .then((m) => {
                let msg = JSONCodec().decode(m.data)
                if (msg.code == 200) {
                  item.description = msg.data.description
                  item.image = msg.data.image
                  item.files.push(...msg.data.entries)
                  item.selected = true
                  setDiff(1)
                }
              })
              .catch((err) => {
                nStore.setNotification(
                  'Catalog failed to respond',
                  'The catalog ' + item.name + ' failed to respond. Moving offline until it heartbeats.'
                )
                item.selected = false
                item.files = [] as File[]
                item.online = false
                item.lastSeen = new Date(0)
              })
          }
        }
      })
    },
    async getShares() {
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
        .catch(() => {
          notificationStore().setNotification(
            'Unable to connect',
            'There is an issue with your connection'
          )
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
            const catalog: Catalog = {
              selected: false,
              to: m.data.account_key,
              from: m.data.account_key,
              name: m.data.catalog,
              online: true,
              lastSeen: Date.now(),
              pending_invite: false,
              status: m.data.revision,
              files: [] as File[]
            }
            this.catalogs.push(catalog)

            uStore.setLastSeenTS(Date.now())
            uStore.setPendingInvites(m.data.unimported_shares.length)
            uStore.togglePing()

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
        .catch((err) => {
          notificationStore().setNotification(
            'Not running Natster Catalog',
            'You do not appear to be running a natster catalog'
          )
          this.pending_init = true
        })
    },
    async downloadFile(fileName, catalog, hash, mimeType) {
      this.activeDownload = true
      let xkey = createCurve()
      this.xkey_seed = new TextDecoder().decode(xkey.getSeed())
      this.xkey_pub = xkey.getPublicKey()

      var fileArray
      const nStore = natsStore()
      const sub = nStore.connection.subscribe('natster.media.' + catalog.name + '.' + hash)
        ; (async () => {
          for await (const m of sub) {
            const chunkIdx = parseInt(m.headers.get('x-natster-chunk-idx'))
            const totalChunks = parseInt(m.headers.get('x-natster-total-chunks'))
            const senderXKey = m.headers.get('x-natster-sender-xkey')

            let decrypted = xkey.open(m.data, senderXKey)
            fileArray.push(decrypted)

            if (chunkIdx === totalChunks - 1) {
              this.activeDownload = false
              sub.unsubscribe()
            }
          }

          console.log('DOWNLOAD FILE', fileArray)
          var blob = new Blob(fileArray, { type: mimeType })
          console.log('DOWNLOAD FILE', blob)
          saveAs(blob, fileName)
        })()


      const dl_request = {
        hash: hash,
        transcode: false,
        target_xkey: this.xkey_pub
      }
      await nStore.connection
        .request('natster.catalog.' + catalog.name + '.download', JSON.stringify(dl_request), {
          timeout: 5000
        })
        .then((m) => {
          let data = JSONCodec().decode(m.data)
          fileArray = new Array()
        })
        .catch((err) => {
          console.error('nats download request err: ', err)
          notificationStore().setNotification('NATS request failed', err)
          sub.unsubscribe()
        })
    },
    async viewFile(fileName, fileDescription, catalog, hash, mimeType) {
      const fStore = fileStore()

      let xkey = createCurve()
      this.xkey_seed = new TextDecoder().decode(xkey.getSeed())
      this.xkey_pub = xkey.getPublicKey()

      var fileArray
      const nStore = natsStore()
      const sub = nStore.connection.subscribe('natster.media.' + catalog.name + '.' + hash)

        ; (async () => {
          let timeout
          let lastChunkReceivedTimestamp: number

          if (mimeType.toLowerCase().indexOf('video/') !== 0 || mimeType.toLowerCase().indexOf('audio/') !== 0) {
            fStore.load(fileName, fileDescription, mimeType, catalog, () => {
              if (timeout) {
                clearTimeout(timeout)
                timeout = null
              }

              fStore.endStream()
              sub.unsubscribe()
            })
          }

          for await (const m of sub) {
            const chunkIdx = parseInt(m.headers.get('x-natster-chunk-idx'))
            const senderXKey = m.headers.get('x-natster-sender-xkey')
            const totalChunks = parseInt(m.headers.get('x-natster-total-chunks'))
            const transcoding = m.headers.get('x-natster-transcoding') && m.headers.get('x-natster-transcoding') === 'true'
            const decrypted = xkey.open(m.data, senderXKey)
            lastChunkReceivedTimestamp = Date.now()

            if (mimeType.toLowerCase().indexOf('video/') === 0) {
              if (timeout) {
                clearTimeout(timeout)
                timeout = null
              }

              fStore.render(fileName, fileDescription, mimeType, decrypted, catalog)

              if (transcoding) {
                timeout = setTimeout(() => {
                  if (chunkIdx >= totalChunks - 1) {
                    // HACK-- x-natster-total-chunks is lower than the actual number of chunks when transcoding video/mp4 on-the-fly
                    // HACK-- this branch prevents slow streams from being canceled early while we are still transcoding
                    fStore.endStream()
                    timeout = null

                    sub.unsubscribe()
                  } else {
                    // TODO-- maintain a tolerance for max time we will wait for the next packet-- this can eventually replace the above HACK
                    // TODO-- navigator.connection.addEventListener('change', () => { // prevent/cancel streams })
                    console.log(`WARNING-- no packet received since ${lastChunkReceivedTimestamp}`)
                  }
                }, 5000)
              }
            } else if (mimeType.toLowerCase() === 'audio/mpeg') {
              fStore.render(fileName, fileDescription, mimeType, decrypted, catalog)
            } else {
              fileArray.push(decrypted)

              if (chunkIdx === totalChunks - 1) {
                console.log('VIEW FILE', fileArray)
                let blob = new Blob(fileArray, { type: mimeType })
                console.log('VIEW FILE', blob)
                fStore.render(fileName, fileDescription, mimeType, blob, catalog)
              }
            }

            if (!transcoding && chunkIdx === totalChunks - 1) {
              fStore.endStream()
              sub.unsubscribe()
            }
          }
          console.log('subscription closed')
        })()

      const dl_request = {
        hash: hash,
        transcode: true,
        target_xkey: this.xkey_pub
      }
      await nStore.connection
        .request('natster.catalog.' + catalog.name + '.download', JSON.stringify(dl_request), {
          timeout: 5000
        })
        .then((m) => {
          let data = JSONCodec().decode(m.data)
          fileArray = new Array()
        })
        .catch((err) => {
          console.error('nats download request err: ', err)
          notificationStore().setNotification('NATS request failed', err)
          sub.unsubscribe()
          fStore.reset()
        })
    },
    isMimeTypeSupported(inMime: string) {
      return this.supportedMimeType.includes(inMime)
    },
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
      state.catalogs.forEach(function(tCatalog) {
        if (tCatalog.from == 'AC5V4OC2POUAX4W4H7CKN5TQ5AKVJJ4AJ7XZKNER6P6DHKBYGVGJHSNC') {
          tCatalog.pending_invite = false // synadiahub is never pending
          return
        }
        state.pending_catalogs.forEach(function(tPending) {
          if (tCatalog.name === tPending.name) {
            tCatalog.pending_invite = true
          }
        })

      })

      return state.catalogs.filter((c) => !c.pending_invite)
    }
  }
})
