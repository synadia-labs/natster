import { defineStore } from 'pinia'
import { jwtDecode } from 'jwt-decode'
import { useAuth0 } from '@auth0/auth0-vue'

export const userStore = defineStore('user', {
  state: () => ({
    // TODO: YOU NEED TO MANUALLY DO THIS FOR NOW AND SET LOGGEDIN TO TRUE
    jwt: "",
    nkey: "",
    user: null,
    nats_code: '',
    authenticated: false,
    loading: false,
    decoded_jwt: null,
    shares: []
  }),
  getters: {
    getUser(state) {
      const { user } = useAuth0()
      console.log(state.user)
      return user
    },
    getNatsCode(state) {
      return state.nats_code
    },
    myShares(state) {
      state.decoded_jwt = jwtDecode(state.jwt)
      const ret = []
      state.shares.forEach(function (item, index) {
        if (item.from_account == state.decoded_jwt.nats.issuer_account) {
          ret.push(item)
        }
      })
      return ret
    },
    getAccount(state) {
      const decoded_jwt = jwtDecode(state.jwt)
      return decoded_jwt.nats.issuer_account
    },
    getUserName(state) {
      const decoded_jwt = jwtDecode(state.jwt)
      return decoded_jwt.name
    },
    getUserPhotoUrl(state) {
      const decoded_jwt = jwtDecode(state.jwt)
      const tags = decoded_jwt.nats.tags
      for (let i = 0; i < tags.length; i++) {
        const t = tags[i].trim()
        if (t.startsWith('photo_url:')) {
          return t.substring(t.indexOf(':') + 1)
        }
      }
      return 'https://ui-avatars.com/api/?name=' + decoded_jwt.name.replace(' ', '+')
    }
  }
})
