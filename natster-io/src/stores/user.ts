import { defineStore } from 'pinia'
import { jwtDecode } from 'jwt-decode'
import { useAuth0 } from '@auth0/auth0-vue'

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
    last_seen_ts: new Date(0),
    ping: false,
    pending_imports: 0
  }),
  actions: {
    setJWT(jwt: string) {
      this.jwt = jwt
      localStorage.setItem('natster_jwt', jwt)
    },
    setNkey(nkey: string) {
      this.nkey = nkey
      localStorage.setItem('natster_nkey', nkey)
    },
    setOauthId(id: string) {
      this.oauth_id = id
      localStorage.setItem('natster_oauth_id', id)
    },
    setCatalogOnline(online: boolean) {
      this.catalog_online = online
    },
    setPendingInvites(pending) {
      this.pending_imports = pending
    },
    setLastSeenTS(ts: Date) {
      this.last_seen_ts = ts
    },
    togglePing() {
      this.ping = !this.ping
    }
  },
  getters: {
    hasJWT(state) {
      return state.jwt !== '' || localStorage.getItem('natster_jwt') !== null
    },
    hasNkey(state) {
      return state.nkey !== '' || localStorage.getItem('natster_nkey') !== null
    },
    getLastSeen(state) {
      return state.last_seen_ts
    },
    getCatalogOnline(state) {
      state.ping
      return (Date.now() - new Date(state.last_seen_ts).getTime() < (1 * 60 * 1000))
    },
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
