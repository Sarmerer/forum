import api from "@/router/api";

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
    async signUp({ dispatch, commit }, credentials) {
      await api
        .post(`auth/signup`, credentials)
        .then((response) => {
          response?.data?.data
            ? commit("setUser", response.data.data)
            : dispatch("attempt");
        })
        .catch((error) => console.log(error));
    },
    async signIn({ dispatch, commit }, credentials) {
      await api
        .post(`auth/signin`, credentials)
        .then((response) => {
          console.log(response.data.data);
          response?.data?.data
            ? commit("setUser", response.data.data)
            : dispatch("attempt");
        })
        .catch((error) => console.log(error));
    },
    async signOut({ commit }) {
      return api.post("auth/signout").finally(() => {
        commit("setUser", null);
      });
    },
    async attempt({ commit }) {
      return await api
        .get("auth/me")
        .then((response) => commit("setUser", response.data.data))
        .catch(() => commit("setUser", null));
    },
  },
};
