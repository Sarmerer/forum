import api from "@/api/api";

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
    async OAuth({ dispatch, commit }, query) {
      await api
        .post(`oauth`, null, {
          params: query,
        })
        .then((response) => {
          let user = response?.data?.data;
          user ? commit("setUser", user) : dispatch("attempt");
        })
        .catch((error) => commit("setAuthError", error));
    },
    async verify({ dispatch, commit }, code) {
      await api
        .post(`auth/verify?code=${code}`)
        .then((response) => {
          let user = response?.data?.data;
          user ? commit("setUser", user) : dispatch("attempt");
        })
        .catch((error) => commit("setAuthError", error));
    },
    async signUp({ commit }, credentials) {
      await api
        .post(`auth/signup`, {
          login: credentials.login,
          email: credentials.email,
          password: credentials.password,
          admin_token: credentials.adminToken || "",
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
      await api.post("auth/signout").then(() => {
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
