import { defineStore } from 'pinia'
import { jwtDecode } from 'jwt-decode'
import { useAuth0 } from '@auth0/auth0-vue'
import { natsStore } from './nats'
import type { Catalog } from '../types/types.ts'
import { JSONCodec, StringCodec } from 'nats.ws'
import init, { get_xkeys, decrypt_chunk } from '../wasm/generate-xkeys/pkg/generate_xkeys.js'
import { textFileStore } from './textfile'

export const userStore = defineStore('user', {
  state: () => ({
    jwt: '',
    nkey: '',
    user: null,
    nats_code: '',
    authenticated: false,
    loading: false,
    decoded_jwt: null,
    catalogs: [] as Catalog[],
    files: [] as File[],
    xkey_seed: '',
    xkey_pub: ''
  }),
  actions: {
    setJWT(jwt) {
      this.jwt = jwt
    },
    setNkey(nkey) {
      this.nkey = nkey
    },
    setNatsCode(code) {
      this.nats_code = code
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

      await init()
      var buf = new Uint8Array(32)
      window.crypto.getRandomValues(buf)

      let xkey = JSON.parse(get_xkeys(buf))
      this.xkey_seed = xkey.seed
      this.xkey_pub = xkey.public

      var sender_xkey
      const nStore = natsStore()
      const sub = nStore.connection.subscribe('natster.media.' + catalog + '.' + hash)
      ;(async () => {
        for await (const m of sub) {
          await new Promise((r) => setTimeout(r, 1000))
          console.log('mine', this.xkey_seed)
          console.log('theirs', sender_xkey)
          console.log('data', m.data)
          let decrypted = decrypt_chunk(m.data, this.xkey_seed, sender_xkey)
          console.log('DECRYPTED: ', decrypted)

          tfStore.showTextFile(fileName, decrypted)
          //console.log(`[${sub.getProcessed()}]: ${StringCodec().decode(m.data)}`)
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
          console.log(m.data)
          sender_xkey = data.data.sender_xkey
        })
        .catch((err) => {
          console.error('nats requestCatalogFiles err: ', err)
        })
    }
  },
  getters: {
    getJWT(state) {
      return state.jwt
    },
    getNkey(state) {
      return state.nkey
    },
    getNatsCode(state) {
      return state.nats_code
    },
    getUser(state) {
      const { user } = useAuth0()
      return user
    },
    getCatalogs(state) {
      return state.catalogs
    },
    getAccount(state) {
      if (state.jwt == '') {
        return ''
      }
      const decoded_jwt = jwtDecode(state.jwt)
      return decoded_jwt.nats.issuer_account
    },
    getUserName(state) {
      const { user } = useAuth0()
      return user.name
    }
  }
})
