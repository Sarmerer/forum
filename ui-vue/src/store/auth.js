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
    async signIn({ dispatch }, credentials) {
      let response = await axios.post(`auth/signin`, credentials);
      dispatch("attempt", response.data.data);
    },
    async signOut({ commit }) {
      return axios.post("auth/signout").then(() => {
        commit("setUser", null);
      });
    },
    async attempt({ commit }) {
      try {
        let response = await axios.get("auth/me");
        commit("setUser", response.data.data);
      } catch (error) {
        commit("setUser", null);
      }
    },
  },
};
