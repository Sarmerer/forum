import axios from "axios";

export default {
  namespaced: true,
  state: {
    user: null,
  },
  getters: {
    authenticated(state) {
      return state.user != null;
    },
    user(state) {
      return state.user;
    },
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
    },
  },
  actions: {
    async signUp({ dispatch }, credentials) {
      await axios
        .post(`auth/signup`, credentials)
        .then(() => {
          dispatch("attempt");
        })
        .catch((error) => console.log(error));
    },
    async signIn({ dispatch }, credentials) {
      await axios
        .post(`auth/signin`, credentials)
        .then(() => {
          dispatch("attempt");
        })
        .catch((error) => console.log(error));
    },
    async signOut({ commit }) {
      //FIXME fix case when you log in from a new place, and then try to log out from the first session
      return axios.post("auth/signout").then(() => {
        commit("setUser", null);
      });
    },
    async attempt({ commit }) {
      return await axios
        .get("auth/me")
        .then((response) => commit("setUser", response.data.data))
        .catch(() => commit("setUser", null));
    },
  },
};
