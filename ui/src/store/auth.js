import api from "@/router/api";

export default {
  namespaced: true,
  state: {
    user: null,
    authError: null,
  },
  getters: {
    authenticated(state) {
      return state.user != null;
    },
    user(state) {
      return state.user;
    },
    authError(state) {
      return state.authError;
    },
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
    },
    setAuthError(state, error) {
      state.authError = error;
    },
  },
  actions: {
    async signUp({ dispatch, commit }, credentials) {
      await api
        .post(`auth/signup`, {
          login: credentials.login,
          email: credentials.email,
          password: credentials.password,
          admin: credentials.admin || false,
          admin_token: credentials.adminToken || "",
        })
        .then((response) => {
          let user = response?.data?.data;
          user ? commit("setUser", user) : dispatch("attempt");
        })
        .catch((error) => commit("setAuthError", error));
    },
    async signIn({ dispatch, commit }, credentials) {
      await api
        .post(`auth/signin`, credentials)
        .then((response) => {
          let user = response?.data?.data;
          user ? commit("setUser", user) : dispatch("attempt");
        })
        .catch((error) => commit("setAuthError", error));
    },
    async signOut({ commit }) {
      await api.post("auth/signout").finally(() => {
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
