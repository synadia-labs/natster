import { defineStore } from 'pinia'
import { jwtDecode } from 'jwt-decode'
import { useAuth0 } from '@auth0/auth0-vue'
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
