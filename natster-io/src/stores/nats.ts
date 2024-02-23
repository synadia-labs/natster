import { defineStore } from 'pinia'
import { userStore } from './user'
import { connect, jwtAuthenticator, JSONCodec } from 'nats.ws'

export const natsStore = defineStore('nats', {
  state: () => ({
    name: 'natster_ui',
    servers: 'connect.ngs.global',
    connection: null,
    connected: false
  }),
  actions: {
    async connect() {
      const uStore = userStore()
      if (uStore.nkey != '' && uStore.jwt != '') {
        try {
          const conn = await connect({
            debug: true,
            ignoreClusterUpdates: true,
            servers: this.servers,
            reconnect: true,
            authenticator: jwtAuthenticator(uStore.jwt, new TextEncoder().encode(uStore.nkey))
          })
          this.connection = conn
        } catch (err) {
          console.error('nats connect err: ', err)
        } finally {
          this.connected = true
        }
      }
    },
    async disconnect() {
      try {
        if (this.connection) {
          await this.connection.close()
        }
      } catch (err) {
        console.error('nats disconnect err: ', err)
      } finally {
        this.connection = null
        this.connected = false
      }
    },
    async getShares(init) {
      const uStore = userStore()
      await this.connection
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
                uStore.catalogs.push(catalog)
              }
            })
          }
        })
        .catch((err) => console.error('nats shares err: ', err))
    },

    async ping() {
      const uStore = userStore()
      await this.connection
        .request('natster.local.inbox', '', { timeout: 5000 })
        .then((msg) => {
          let m = JSONCodec().decode(msg.data)
          if (m.code == 200) {
            uStore.catalog_online = true
            uStore.pending_imports = m.data.unimported_shares.length
          }
        })
        .catch((err) => console.error('nats ping err: ', err))
    }
  }
})
