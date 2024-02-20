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
    async ping() {
      try {
        await this.connection.ping()
      } catch (err) {
        console.error('nats ping err: ', err)
      }
    }
  }
})
