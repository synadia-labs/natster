import { defineStore } from "pinia";
import { jwtDecode } from "jwt-decode";

export const userStore = defineStore("user", {
  state: () => ({
    // TODO: YOU NEED TO MANUALLY DO THIS FOR NOW AND SET LOGGEDIN TO TRUE
    jwt: "",
    nkey: "",
    loggedIn: false,
    decoded_jwt: null,
    shares: [],
  }),
  actions: {
    logout() {
      this.jwt = "";
      this.nkey = "";
      this.decoded_jwt = null;
      this.share = [];
      this.loggedIn = false;
    },
  },
  getters: {
    myShares(state) {
      state.decoded_jwt = jwtDecode(state.jwt);
      var ret = [];
      state.shares.forEach(function (item, index) {
        if (item.from_account == state.decoded_jwt.nats.issuer_account) {
          ret.push(item);
        }
      });
      return ret;
    },
    getAccount(state) {
      const decoded_jwt = jwtDecode(state.jwt);
      return decoded_jwt.nats.issuer_account;
    },
    getUserName(state) {
      const decoded_jwt = jwtDecode(state.jwt);
      return decoded_jwt.name;
    },
    getUserPhotoUrl(state) {
      const decoded_jwt = jwtDecode(state.jwt);
      const tags = decoded_jwt.nats.tags;
      for (var i = 0; i < tags.length; i++) {
        let t = tags[i].trim();
        if (t.startsWith("photo_url:")) {
          return t.substring(t.indexOf(":") + 1);
        }
      }
      return (
        "https://ui-avatars.com/api/?name=" + decoded_jwt.name.replace(" ", "+")
      );
    },
  },
});
