import { defineStore } from 'pinia'
import { jwtDecode } from 'jwt-decode'
import { useAuth0 } from '@auth0/auth0-vue'
import { createCurve } from 'nkeys.js'
import { natsStore } from './nats'
import type { Catalog } from '../types/types.ts'
import { JSONCodec, StringCodec } from 'nats.ws'
import { textFileStore } from './textfile'

export const userStore = defineStore('user', {
  state: () => ({
    jwt: '',
    nkey: '',
    user: null,
    oauth_id: '',
    authenticated: false,
    loading: false,
    decoded_jwt: null,
    catalogs: [] as Catalog[],
    files: [] as File[],
    xkey_seed: '',
    xkey_pub: '',
    catalog_online: false,
    pending_imports: 0
  }),
  actions: {
    setJWT(jwt) {
      this.jwt = jwt
      localStorage.setItem('natster_jwt', jwt)
    },
    setNkey(nkey) {
      this.nkey = nkey
      localStorage.setItem('natster_nkey', nkey)
    },
    setOauthId(id) {
      this.oauth_id = id
      localStorage.setItem('natster_oauth_id', id)
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
    getJWT(state) {
      return state.jwt !== '' ? state.jwt : localStorage.getItem('natster_jwt')
    },
    getNkey(state) {
      return state.nkey !== '' ? state.nkey : localStorage.getItem('natster_nkey')
    },
    getOauthId(state) {
      return state.oauth_id !== '' ? state.oauth_id : localStorage.getItem('natster_oauth_id')
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
    },
    getImportedCatalogs(state) {
      console.log('here')
      return state.catalogs.filter((c) => !c.pending_invite)
    }
  }
})
